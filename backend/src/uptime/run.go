package uptime

import (
	"context"
	"github.com/m-milek/leszmonitor/logger"
	"time"
)

func StartUptimeWorker(ctx context.Context) {
	logger.Uptime.Info().Msg("Starting uptime worker...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logger.Uptime.Info().Msg("Uptime worker shutting down...")
			return
		case <-ticker.C:
			logger.Uptime.Trace().Msg("Worker Uptime - running")
		}
	}
}
