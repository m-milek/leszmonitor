package workers

import (
	"context"
	"sync"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

// MonitorExecutor runs monitor checks requested by the scheduler and publishes their results.
type MonitorExecutor struct {
	db     db.DB
	logger zerolog.Logger
}

func NewMonitorExecutor(database db.DB) *MonitorExecutor {
	return &MonitorExecutor{
		db: database,
	}
}

func (e *MonitorExecutor) Run(ctx context.Context) {
	e.logger = log.FromContext(ctx).With().Str("component", "monitor_executor").Logger()

	e.logger.Info().Msg("Starting monitor executor...")

	executeChannel := monitors.MonitorExecuteChannel.Subscribe()
	defer monitors.MonitorExecuteChannel.Unsubscribe(executeChannel)

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		select {
		case <-ctx.Done():
			e.logger.Info().Msg("Monitor executor shutting down...")
			return
		case msg := <-executeChannel:
			wg.Go(func() {
				e.execute(ctx, msg.Monitor)
			})
		}
	}
}

// execute runs the monitor's check and broadcasts the result.
func (e *MonitorExecutor) execute(ctx context.Context, monitor monitors.Monitor) {
	logger := e.logger.With().
		Str("monitor_id", monitor.ID.String()).
		Str("monitor_name", monitor.Name).
		Logger()
	ctx = log.WithContext(ctx, &logger)

	if err := monitor.Validate(); err != nil {
		logger.Error().Err(err).Msg("Monitor validation failed")
		return
	}

	probe, err := probe.UnmarshalProbeFromBytes(monitor.Type, []byte(monitor.ProbeConfig))
	if err != nil {
		logger.Error().Err(err).Msg("Failed to unmarshal probe config")
		return
	}
	if err := probe.Validate(); err != nil {
		logger.Error().Err(err).Msg("Probe config validation failed")
		return
	}

	logger.Trace().Msg("Running monitor")
	result, err := probe.Run(ctx, monitor.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Probe execution failed due to an error")
		return
	}
	logger.Info().Any("monitor_result", result).Msg("Monitor result")

	if result.GetStatus() != kind.MonitorStatusUp {
		if failures := result.GetFailures(); len(failures) > 0 {
			logger.Warn().
				Any("failures", failures).
				Msg("Monitor check failed (service down or misconfigured)")
		}
	}

	monitors.MonitorRunChannel.Broadcast(monitors.MonitorRunMessage{
		Result:  result,
		Monitor: monitor,
	})
}
