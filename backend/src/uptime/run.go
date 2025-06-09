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
		go runMonitor(ctx, monitor)
	}

	monitorMsgChannel := monitors.MessageBroadcaster.Subscribe()
	defer monitors.MessageBroadcaster.Unsubscribe(monitorMsgChannel)

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
			//logger.Uptime.Trace().Msg("Worker Uptime - running")
		case msg := <-monitorMsgChannel:
			if msg.Status == monitors.Created {
				go runMonitor(ctx, *msg.Monitor)
			}
		}
	}
}

func runMonitor(ctx context.Context, monitor monitors.IMonitor) {
	logger.Uptime.Debug().Dur("every", monitor.GetInterval()).Msgf("Starting monitor runner: %s", monitor.GetName())

	monitorMsgChannel := monitors.MessageBroadcaster.Subscribe()
	defer monitors.MessageBroadcaster.Unsubscribe(monitorMsgChannel)

	for {
		select {

		case <-time.After(monitor.GetInterval()):
			// Run the monitor's check method
			err := monitor.Validate()
			if err != nil {
				logger.Uptime.Error().Err(err).Msgf("Error validating monitor %s - %s", monitor.GetName(), monitor.GetId())
				return
			}
			logger.Uptime.Trace().Msgf("Running monitor %s - %s", monitor.GetName(), monitor.GetId())

			response := monitor.Run()

			logger.Uptime.Debug().Any("monitor_response", response).Msgf("Monitor %s - %s response", monitor.GetName(), monitor.GetId())

			// Log the result of the check
			// logger.Uptime.Info().Msgf("Monitor %s result: %v", monitor.GetName(), result)

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
				logger.Uptime.Debug().Msgf("Monitor deleted, stopping runner: %s - %s", monitor.GetName(), monitor.GetId())
				return
			}
			if msg.Status == monitors.Disabled {
				logger.Uptime.Debug().Msgf("Disabling monitor: %s - %s", monitor.GetName(), monitor.GetId())
				return
			}
		case <-ctx.Done():
			logger.Uptime.Info().Msgf("Stopping monitor: %s - %s", monitor.GetName(), monitor.GetId())
			return
		}
	}
}
