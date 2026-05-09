package workers

import (
	"context"
	"time"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
)

const durationBetweenCleanups = time.Duration(60) * time.Second

func StartDataCleanupWorker(ctx context.Context) {
	logger := log.FromContext(ctx).With().Str("component", "data_cleanup_worker").Logger()
	ctx = log.WithContext(ctx, &logger)

	logger.Info().Msg("Starting data cleanup worker...")

	ticker := time.NewTicker(durationBetweenCleanups)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Data cleanup worker shutting down...")
			return
		case <-ticker.C:
			allMonitors, err := db.Get().Monitors().GetAllMonitors(ctx)

			if err != nil {
				logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
				return
			}
			logger.Debug().Msgf("Starting data cleanup for %d monitors", len(allMonitors))
			for _, monitor := range allMonitors {
				_, err := db.Get().MonitorResults().DeleteMonitorResultsOlderThanDuration(ctx, monitor.ID, time.Duration(monitor.ResultRetentionSeconds)*time.Second)
				if err != nil {
					logger.Error().Err(err).Str("monitor_id", monitor.ID.String()).Msg("Failed to delete old monitor results")
				} else {
					logger.Debug().Str("monitor_id", monitor.ID.String()).Msg("Deleted old monitor results successfully")
				}
			}
		}
	}
}
