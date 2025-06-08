package uptime

import (
	"context"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
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
		go runMonitor(ctx, monitor)
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

func runMonitor(ctx context.Context, monitor monitors.IMonitor) {
	logger.Uptime.Trace().Dur("every", monitor.GetInterval()).Msgf("Running monitor: %s", monitor.GetName())

	monitorMsgChannel := monitors.MessageBroadcaster.Subscribe()

	for {
		select {

		case <-time.After(monitor.GetInterval()):
			// Run the monitor's check method
			err := monitor.Validate()
			if err != nil {
				logger.Uptime.Error().Err(err).Msgf("Error validating monitor %s", monitor.GetName())
				return
			}
			logger.Uptime.Debug().Msgf("Running monitor %s", monitor.GetName())

			// Log the result of the check
			//logger.Uptime.Info().Msgf("Monitor %s result: %v", monitor.GetName(), result)

			//	// Optionally, you can save the result to the database or perform further actions
			//	err = db.SaveMonitorResult(monitor.GetId(), result)
			//	if err != nil {
			//		logger.Uptime.Error().Err(err).Msgf("Failed to save result for monitor %s", monitor.GetName())
			//	}
		case msg := <-monitorMsgChannel:
			if msg.Id != monitor.GetId() {
				break
			}
			if msg.Status == monitors.Edited {
				monitor = *msg.Monitor
				logger.Uptime.Debug().Msgf("Updating monitor: %s - %s", monitor.GetName(), monitor.GetId())
				break
			}
			if msg.Status == monitors.Deleted {
				logger.Uptime.Info().Msgf("Deleting monitor: %s - %s", monitor.GetName(), monitor.GetId())
				return
			}
			if msg.Status == monitors.Disabled {
				logger.Uptime.Info().Msgf("Disabling monitor: %s - %s", monitor.GetName(), monitor.GetId())
				return
			}
			if msg.Status == monitors.Enabled {
				break
			}
			if msg.Status == monitors.Created {
				break
			}
		case <-ctx.Done():
			logger.Uptime.Info().Msgf("Stopping monitor: %s - %s", monitor.GetName(), monitor.GetId())
			return
		}
	}
}
