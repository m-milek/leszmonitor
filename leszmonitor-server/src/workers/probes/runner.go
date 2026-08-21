package probes

import (
	"context"
	"time"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/models/shared"
	"github.com/rs/zerolog"
)

// probeRunner owns the lifecycle of a single monitor. All state lives in the
// run() goroutine, so no locking is required. Edits arrive via updates (pushed
// by the worker); deletion/shutdown arrive via ctx cancellation.
type probeRunner struct {
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
func (r *probeRunner) run(ctx context.Context) {
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
			r.runCheck(ctx)
		}
	}
}

// interval returns the monitor's interval as a time.Duration.
func (r *probeRunner) interval() time.Duration {
	return time.Duration(r.monitor.Interval) * time.Second
}

// refreshLogger updates the runner's logger with the current monitor name. Called on startup and on edits.
func (r *probeRunner) refreshLogger() {
	r.logger = r.baseLogger.With().Str("monitor_name", r.monitor.Name).Logger()
}

// push delivers an edit to the runner. Newest-wins and non-blocking, so a
// mid-probe runner never stalls the worker's dispatch loop.
func (r *probeRunner) push(mon monitors.Monitor) {
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
func (r *probeRunner) applyUpdate(update monitors.Monitor, ticker *time.Ticker) {
	oldInterval := r.monitor.Interval
	r.monitor = update
	r.refreshLogger()
	r.logger.Debug().Msg("Monitor updated")

	if r.monitor.Interval != oldInterval {
		r.logger.Info().Int("interval", r.monitor.Interval).Msg("Interval changed, resetting ticker")
		ticker.Reset(r.interval())
	}
}

// runCheck executes the monitor's check and handles the result.
func (r *probeRunner) runCheck(ctx context.Context) {
	if r.monitor.RunState != monitors.MonitorStateActive {
		r.logger.Trace().Str("state", string(r.monitor.RunState)).Msg("Skipping run - not active")
		return
	}

	if err := r.monitor.Validate(); err != nil {
		r.logger.Error().Err(err).Msg("Monitor validation failed")
		return
	}

	probe, err := monitors.UnmarshalProbeFromBytes(r.monitor.Type, []byte(r.monitor.ProbeConfig))
	if err != nil {
		r.logger.Error().Err(err).Msg("Failed to unmarshal probe config")
		return
	}
	if err := probe.Validate(); err != nil {
		r.logger.Error().Err(err).Msg("Probe config validation failed")
		return
	}

	r.logger.Trace().Msg("Running monitor")
	result, err := probe.Run(ctx, r.monitor.ID)
	if err != nil {
		r.logger.Error().Err(err).Msg("Probe execution failed due to an error")
	}
	r.logger.Info().Any("monitor_result", result).Msg("Monitor result")

	if result.GetStatus() != shared.MonitorStatusUp {
		d := result.GetErrorDetails()
		if d.ErrorMessage != "" || len(d.Errors) > 0 {
			r.logger.Error().
				Str("error_message", d.ErrorMessage).
				Strs("errors", d.Errors).
				Msg("Monitor encountered internal error")
		}
		if len(d.Failures) > 0 {
			r.logger.Warn().
				Strs("failures", d.Failures).
				Msg("Monitor check failed (service down or misconfigured)")
		}
	}

	if _, err := r.db.MonitorResults().InsertMonitorResult(ctx, result); err != nil {
		r.logger.Error().Err(err).Msg("Failed to insert monitor result")
	}

	events.MonitorRunChannel.Broadcast(monitors.MonitorRunMessage{
		Result:  result,
		Monitor: r.monitor,
	})
}
