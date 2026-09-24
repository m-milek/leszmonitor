package workers

import (
	"context"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

// monitorTicker schedules a single monitor at its interval. All state lives in the
// run() goroutine, so no locking is required. Edits arrive via updates (pushed
// by the scheduler); deletion/shutdown arrive via ctx cancellation.
type monitorTicker struct {
	monitor monitors.Monitor
	db      db.DB
	cancel  context.CancelFunc
	updates chan monitors.Monitor
	onExit  func() // called for self-termination

	baseLogger zerolog.Logger
	logger     zerolog.Logger
}

// run handles the lifecycle of a single monitor.
// It runs checks at intervals, applying edits, and self-terminates on invalid configuration.
func (t *monitorTicker) run(ctx context.Context) {
	defer t.cancel()

	t.refreshLogger()
	ctx = log.WithContext(ctx, &t.logger)

	if err := t.monitor.Validate(); err != nil {
		t.logger.Error().Err(err).Msg("Initial validation failed; ticker not started")
		t.onExit() // self-termination: ask the manager to drop us from the map
		return
	}

	t.logger.Debug().Int("interval", t.monitor.Interval).Msg("Starting monitor ticker")

	ticker := time.NewTicker(t.interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			t.logger.Info().Msg("Stopping monitor ticker")
			return
		case mon := <-t.updates:
			t.applyUpdate(mon, ticker)
		case <-ticker.C:
			t.schedule()
		}
	}
}

// interval returns the monitor's interval as a time.Duration.
func (t *monitorTicker) interval() time.Duration {
	return time.Duration(t.monitor.Interval) * time.Second
}

// refreshLogger updates the ticker's logger with the current monitor name. Called on startup and on edits.
func (t *monitorTicker) refreshLogger() {
	t.logger = t.baseLogger.With().Str("monitor_name", t.monitor.Name).Logger()
}

// push delivers an edit to the ticker. Newest-wins and non-blocking, so a
// busy ticker never stalls the scheduler's dispatch loop.
func (t *monitorTicker) push(mon monitors.Monitor) {
	for {
		select {
		case t.updates <- mon:
			return
		default:
			select {
			case <-t.updates:
			default:
			}
		}
	}
}

// applyUpdate applies an incoming edit to the ticker's state.
// If the interval has changed, the ticker is reset to the new duration.
func (t *monitorTicker) applyUpdate(update monitors.Monitor, ticker *time.Ticker) {
	oldInterval := t.monitor.Interval
	t.monitor = update
	t.refreshLogger()
	t.logger.Debug().Msg("Monitor updated")

	if t.monitor.Interval != oldInterval {
		t.logger.Info().Int("interval", t.monitor.Interval).Msg("Interval changed, resetting ticker")
		ticker.Reset(t.interval())
	}
}

// schedule sends the monitor to the executor if it is active.
func (t *monitorTicker) schedule() {
	if t.monitor.RunState != monitors.MonitorStateActive {
		t.logger.Trace().Str("state", string(t.monitor.RunState)).Msg("Skipping run - not active")
		return
	}

	monitors.MonitorExecuteChannel.Broadcast(monitors.MonitorExecuteMessage{Monitor: t.monitor})
}
