package workers

import (
	"context"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

// monitorRunner owns the lifecycle of a single monitor. All state lives in the
// run() goroutine, so no locking is required. Edits arrive via updates (pushed
// by the worker); deletion/shutdown arrive via ctx cancellation.
type monitorRunner struct {
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
func (r *monitorRunner) run(ctx context.Context) {
	defer r.cancel()

	r.refreshLogger()
	ctx = log.WithContext(ctx, &r.logger)

	if err := r.monitor.Validate(); err != nil {
		r.logger.Error().Err(err).Msg("Initial validation failed; runner not started")
		r.onExit() // self-termination: ask the manager to drop us from the map
		return
	}

	r.logger.Debug().Int("interval", r.monitor.Interval).Msg("Starting probe runner")

	ticker := time.NewTicker(r.interval())
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.logger.Info().Msg("Stopping probe")
			return
		case mon := <-r.updates:
			r.applyUpdate(mon, ticker)
		case <-ticker.C:
			r.schedule()
		}
	}
}

// interval returns the monitor's interval as a time.Duration.
func (r *monitorRunner) interval() time.Duration {
	return time.Duration(r.monitor.Interval) * time.Second
}

// refreshLogger updates the runner's logger with the current monitor name. Called on startup and on edits.
func (r *monitorRunner) refreshLogger() {
	r.logger = r.baseLogger.With().Str("monitor_name", r.monitor.Name).Logger()
}

// push delivers an edit to the runner. Newest-wins and non-blocking, so a
// mid-probe runner never stalls the worker's dispatch loop.
func (r *monitorRunner) push(mon monitors.Monitor) {
	for {
		select {
		case r.updates <- mon:
			return
		default:
			select {
			case <-r.updates:
			default:
			}
		}
	}
}

// applyUpdate applies an incoming edit to the runner's state.
// If the interval has changed, the ticker is reset to the new duration.
func (r *monitorRunner) applyUpdate(update monitors.Monitor, ticker *time.Ticker) {
	oldInterval := r.monitor.Interval
	r.monitor = update
	r.refreshLogger()
	r.logger.Debug().Msg("Monitor updated")

	if r.monitor.Interval != oldInterval {
		r.logger.Info().Int("interval", r.monitor.Interval).Msg("Interval changed, resetting ticker")
		ticker.Reset(r.interval())
	}
}

// schedule sends the monitor to the executor if it is active.
func (r *monitorRunner) schedule() {
	if r.monitor.RunState != monitors.MonitorStateActive {
		r.logger.Trace().Str("state", string(r.monitor.RunState)).Msg("Skipping run - not active")
		return
	}

	monitors.MonitorExecuteChannel.Broadcast(monitors.MonitorExecuteMessage{Monitor: r.monitor})
}
