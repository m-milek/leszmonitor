package probes

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/rs/zerolog"
)

// ProbesManager supervises the set of live monitor runners.
type ProbesManager struct {
	mu     sync.RWMutex
	db     db.DB
	probes map[uuid.UUID]*probeRunner
	exited chan uuid.UUID
	logger zerolog.Logger
}

// NewProbesManager returns a pointer: ProbesManager holds a sync.RWMutex and must
// never be copied (go vet copylocks).
func NewProbesManager(database db.DB) *ProbesManager {
	return &ProbesManager{
		db:     database,
		probes: make(map[uuid.UUID]*probeRunner),
		exited: make(chan uuid.UUID, 100),
	}
}

func (w *ProbesManager) Run(ctx context.Context) {
	w.logger = log.FromContext(ctx).With().Str("component", "probes_worker").Logger()

	w.logger.Info().Msg("Starting probes worker...")

	monitorMsgChannel := events.MonitorLifecycleChannel.Subscribe()
	defer events.MonitorLifecycleChannel.Unsubscribe(monitorMsgChannel)

	allMonitors, err := w.db.Monitors().GetAllMonitors(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return
	}
	w.logger.Debug().Msgf("Found %d monitors to check", len(allMonitors))

	for _, monitor := range allMonitors {
		w.start(ctx, monitor)
	}

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("Probes worker shutting down...")
			return
		case msg := <-monitorMsgChannel:
			w.dispatch(ctx, msg)
		case exitedID := <-w.exited:
			w.remove(exitedID)
		}
	}
}

func (w *ProbesManager) dispatch(ctx context.Context, msg monitors.MonitorLifecycleMessage) {
	switch msg.Status {
	case monitors.Created:
		if msg.Monitor != nil {
			w.start(ctx, *msg.Monitor)
		}
	case monitors.Edited:
		updated, err := w.db.Monitors().GetMonitorByID(ctx, msg.ID)
		if err != nil {
			w.logger.Error().Err(err).Str("monitor_id", msg.ID.String()).Msg("Failed to refetch monitor after edit")
			return
		}
		if r := w.get(msg.ID); r != nil {
			r.push(*updated)
		}
	case monitors.Deleted:
		w.stop(msg.ID)
	}
}

func (w *ProbesManager) start(ctx context.Context, monitor monitors.Monitor) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.probes[monitor.ID]; exists {
		w.logger.Warn().
			Str("monitor_id", monitor.ID.String()).
			Str("monitor_name", monitor.Name).
			Msg("Probe already running; ignoring start request")
		return
	}

	childContext, cancel := context.WithCancel(ctx)

	probe := &probeRunner{
		monitor:    monitor,
		db:         w.db,
		cancel:     cancel,
		updates:    make(chan monitors.Monitor, 1),
		onExit:     func() { w.notifyExit(childContext, monitor.ID) },
		baseLogger: w.logger.With().Str("monitor_id", monitor.ID.String()).Logger(),
	}

	w.probes[monitor.ID] = probe

	go probe.run(childContext)
}

// stop cancels a running probe and removes it.
func (w *ProbesManager) stop(id uuid.UUID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	r, ok := w.probes[id]
	if !ok {
		w.logger.Warn().Str("monitor_id", id.String()).Msg("Attempted to stop a probe, but it was not found")
		return
	}

	delete(w.probes, id)
	r.cancel()
}

// remove drops a probe that has already exited on its own.
func (w *ProbesManager) remove(id uuid.UUID) {
	w.mu.Lock()
	delete(w.probes, id)
	w.mu.Unlock()
}

// get returns a probeRunner for a given monitor uuid.UUID.
func (w *ProbesManager) get(id uuid.UUID) *probeRunner {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.probes[id]
}

// notifyExit sends a message to the worker's main loop that a probe has exited on its own and can be removed from the map.
// This is used for self-termination when a probe detects an invalid configuration and fails.
func (w *ProbesManager) notifyExit(ctx context.Context, id uuid.UUID) {
	select {
	case w.exited <- id:
	case <-ctx.Done():
	}
}

// ActiveCount returns the number of currently active probes.
func (w *ProbesManager) ActiveCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.probes)
}
