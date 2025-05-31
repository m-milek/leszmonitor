package uptime

import (
	"context"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"time"
)

func StartUptimeWorker(ctx context.Context) {
	logger.Uptime.Info().Msg("Starting uptime worker...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	allMonitors, err := db.GetAllMonitors()

	if err != nil {
		logger.Uptime.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return
	}
	logger.Uptime.Info().Any("monitors", allMonitors).Msgf("Found %d monitors to check", len(allMonitors))

	for _, monitor := range allMonitors {
		logger.Uptime.Debug().Type("type", monitor).Msgf("Starting uptime monitor: %s", monitor.GetName())
	}

	for {
		select {
		case <-ctx.Done():
			logger.Uptime.Info().Msg("Uptime worker shutting down...")
			return
		case <-ticker.C:
			//res, err := db.Ping()
			//if err != nil {
			//	println(res)
			//}
			// logger.Uptime.Trace().Msg("Worker Uptime - running")
		}
	}
}
