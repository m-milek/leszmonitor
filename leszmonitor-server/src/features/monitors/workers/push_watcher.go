package workers

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

const durationBetweenPushChecks = 10 * time.Second

// PushWatcher emits results for push monitors
type PushWatcher struct {
	db        db.DB
	startedAt time.Time
	logger    zerolog.Logger
}

func NewPushWatcher(database db.DB) *PushWatcher {
	return &PushWatcher{
		db: database,
	}
}

func (w *PushWatcher) Run(ctx context.Context) {
	w.logger = log.FromContext(ctx).With().Str("component", "push_watcher").Logger()
	w.startedAt = time.Now().UTC()

	w.logger.Info().Msg("Starting push watcher...")

	ticker := time.NewTicker(durationBetweenPushChecks)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("Stopping push watcher")
			return
		case <-ticker.C:
			w.check(ctx, time.Now().UTC())
		}
	}
}

func (w *PushWatcher) check(ctx context.Context, now time.Time) {
	allMonitors, err := monitors.NewMonitorDAO(w.db.Querier()).GetAllMonitors(ctx)
	if err != nil {
		w.logger.Error().Err(err).Msg("Failed to fetch monitors for PushWatcher check")
		return
	}
	activePushMonitors := slices.DeleteFunc(allMonitors, func(monitor monitors.Monitor) bool {
		return monitor.Type != kind.PushConfigType || monitor.RunState != monitors.MonitorRunStateActive
	})

	for _, monitor := range activePushMonitors {
		w.checkMonitor(ctx, monitor, now)
	}
}

func (w *PushWatcher) checkMonitor(ctx context.Context, monitor monitors.Monitor, now time.Time) {
	latestResult, err := results.NewMonitorResultDAO(w.db.Querier()).GetLatestMonitorResultByMonitorID(ctx, monitor.ID)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		w.logger.Error().Err(err).Str("monitor_id", monitor.ID.String()).Msg("Failed to get latest monitor result for PushWatcher check")
		return
	}

	baseline := w.baseline(monitor, latestResult)

	probeConfig, err := probe.Parse[*probe.PushProbe](monitor.Type, monitor.ProbeConfig)
	if err != nil {
		w.logger.Error().Err(err).Str("monitor_id", monitor.ID.String()).Msg("Failed to parse push probe config")
		return
	}
	failThreshold := time.Duration(monitor.Interval+probeConfig.GracePeriodSeconds) * time.Second
	elapsed := now.Sub(baseline)
	if elapsed > failThreshold {
		result := results.NewMonitorResult(monitor.ID, kind.PushConfigType, kind.MonitorStatusDown, false, nil, &results.PushResultDetails{})
		result.AddFailure(results.FailureReasonPushMissedHeartbeat, nil, nil)
		monitors.MonitorRunChannel.Broadcast(monitors.MonitorRunMessage{
			Monitor: monitor,
			Result:  &result,
		})
	}
}

func (w *PushWatcher) baseline(monitor monitors.Monitor, latestResult results.IMonitorResult) time.Time {
	candidates := []time.Time{w.startedAt, monitor.CreatedAt}
	if monitor.UpdatedAt != nil {
		candidates = append(candidates, *monitor.UpdatedAt)
	}
	if latestResult != nil {
		candidates = append(candidates, latestResult.GetCreatedAt())
	}

	return slices.MaxFunc(candidates, time.Time.Compare)
}
