package uptime

import (
	"context"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
	"sync"
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

	editMonitorMutex := sync.Mutex{}

	shutdownChannel := make(chan struct{})

	// Start a goroutine to listen for monitor updates and handle them asynchronously from the main monitor loop
	go func() {
		for msg := range monitorMsgChannel {
			if msg.Id != monitor.GetId() {
				continue
			}

			editMonitorMutex.Lock()

			switch msg.Status {
			case monitors.Edited:
				monitor = *msg.Monitor
				logger.Uptime.Debug().Str("name", monitor.GetName()).Str("id", monitor.GetId()).Msg("Updating monitor")
			case monitors.Deleted, monitors.Disabled:
				shutdownChannel <- struct{}{}
				editMonitorMutex.Unlock()
				return
			}

			editMonitorMutex.Unlock()
		}
	}()

	for {
		select {
		case <-time.After(monitor.GetInterval()):
			// Run the monitor's check method
			err := monitor.Validate()
			if err != nil {
				logger.Uptime.Error().Err(err).Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Error validating monitor")
				return
			}
			logger.Uptime.Trace().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Running monitor")

			response := monitor.Run()

			logger.Uptime.Debug().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Any("monitor_response", response).Msg("Monitor response")

			// Log the result of the check
			// logger.Uptime.Info().Msgf("Monitor %s result: %v", monitor.GetName(), result)

			//	// Optionally, you can save the result to the database or perform further actions
			//	err = db.SaveMonitorResult(monitor.GetId(), result)
			//	if err != nil {
			//		logger.Uptime.Error().Err(err).Msgf("Failed to save result for monitor %s", monitor.GetName())
			//	}
		case shutdownMsg := <-shutdownChannel:
			if shutdownMsg == struct{}{} {
				logger.Uptime.Info().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Stopping monitor")
				editMonitorMutex.Unlock()
				return
			}
		case <-ctx.Done():
			logger.Uptime.Info().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Stopping monitor")
			return
		}
	}
}
