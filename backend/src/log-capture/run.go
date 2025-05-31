package log_capture

import (
	"context"
	"github.com/m-milek/leszmonitor/logger"
	"time"
)

func StartLogWorker(ctx context.Context) {
	logger.Logs.Info().Msg("Starting log worker...")
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			//logger.Logs.Info().Msg("Log worker shutting down...")
			return
		case <-ticker.C:
			// logger.Logs.Trace().Msg("Worker LogCapture - running")
		}
	}
}
