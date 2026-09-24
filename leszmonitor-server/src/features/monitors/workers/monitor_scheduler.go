package workers

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

// MonitorScheduler supervises the set of live monitor tickers.
type MonitorScheduler struct {
	mu      sync.RWMutex
	db      db.DB
	tickers map[uuid.UUID]*monitorTicker
	exited  chan uuid.UUID
	logger  zerolog.Logger
}

// NewMonitorScheduler returns a pointer: MonitorScheduler holds a sync.RWMutex and must
// never be copied (go vet copylocks).
func NewMonitorScheduler(database db.DB) *MonitorScheduler {
	return &MonitorScheduler{
		db:      database,
		tickers: make(map[uuid.UUID]*monitorTicker),
		exited:  make(chan uuid.UUID, 100),
	}
}

func (w *MonitorScheduler) Run(ctx context.Context) {
	w.logger = log.FromContext(ctx).With().Str("component", "monitor_scheduler").Logger()

	w.logger.Info().Msg("Starting monitor scheduler...")

	monitorMsgChannel := monitors.MonitorLifecycleChannel.Subscribe()
	defer monitors.MonitorLifecycleChannel.Unsubscribe(monitorMsgChannel)

	allMonitors, err := monitors.NewMonitorDAO(w.db.Querier()).GetAllMonitors(ctx)
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
			w.logger.Info().Msg("Monitor scheduler shutting down...")
			return
		case msg := <-monitorMsgChannel:
			w.dispatchLifecycleMessage(ctx, msg)
		case exitedID := <-w.exited:
			w.remove(exitedID)
		}
	}
}

func (w *MonitorScheduler) dispatchLifecycleMessage(ctx context.Context, msg monitors.MonitorLifecycleMessage) {
	switch msg.Status {
	case monitors.Created:
		if msg.Monitor != nil {
			w.start(ctx, *msg.Monitor)
		}
	case monitors.Edited:
		updated, err := monitors.NewMonitorDAO(w.db.Querier()).GetMonitorByID(ctx, msg.ID)
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

func (w *MonitorScheduler) start(ctx context.Context, monitor monitors.Monitor) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, exists := w.tickers[monitor.ID]; exists {
		w.logger.Warn().
			Str("monitor_id", monitor.ID.String()).
			Str("monitor_name", monitor.Name).
			Msg("Ticker already running; ignoring start request")
		return
	}

	childContext, cancel := context.WithCancel(ctx)

	ticker := &monitorTicker{
		monitor:    monitor,
		db:         w.db,
		cancel:     cancel,
		updates:    make(chan monitors.Monitor, 1),
		onExit:     func() { w.notifyExit(childContext, monitor.ID) },
		baseLogger: w.logger.With().Str("monitor_id", monitor.ID.String()).Logger(),
	}

	w.tickers[monitor.ID] = ticker

	go ticker.run(childContext)
}

// stop cancels a running ticker and removes it.
func (w *MonitorScheduler) stop(id uuid.UUID) {
	w.mu.Lock()
	defer w.mu.Unlock()

	r, ok := w.tickers[id]
	if !ok {
		w.logger.Warn().Str("monitor_id", id.String()).Msg("Attempted to stop a ticker, but it was not found")
		return
	}

	delete(w.tickers, id)
	r.cancel()
}

// remove drops a ticker that has already exited on its own.
func (w *MonitorScheduler) remove(id uuid.UUID) {
	w.mu.Lock()
	delete(w.tickers, id)
	w.mu.Unlock()
}

// get returns a monitorTicker for a given monitor uuid.UUID.
func (w *MonitorScheduler) get(id uuid.UUID) *monitorTicker {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.tickers[id]
}

// notifyExit sends a message to the worker's main loop that a ticker has exited on its own and can be removed from the map.
// This is used for self-termination when a ticker detects an invalid configuration and fails.
func (w *MonitorScheduler) notifyExit(ctx context.Context, id uuid.UUID) {
	select {
	case w.exited <- id:
	case <-ctx.Done():
	}
}

// ActiveCount returns the number of currently active tickers.
func (w *MonitorScheduler) ActiveCount() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.tickers)
}
