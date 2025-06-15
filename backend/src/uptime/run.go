package uptime

import (
	"context"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
	"sync"
	"time"
)

// StartUptimeWorker initializes the uptime worker that periodically checks all monitors.
// It retrieves all monitors from the database and starts a goroutine for each monitor to run its checks.
// It also listens for monitor creation messages to start new monitors dynamically.
// The worker runs until the context is done, allowing for graceful shutdown.
func StartUptimeWorker(ctx context.Context) {
	logger.Uptime.Info().Msg("Starting uptime worker...")

	allMonitors, err := db.GetAllMonitors()

	if err != nil {
		logger.Uptime.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return
	}
	logger.Uptime.Info().Any("monitors", allMonitors).Msgf("Found %d monitors to check", len(allMonitors))

	for _, monitor := range allMonitors {
		childContext, cancel := context.WithCancel(ctx)
		go runMonitor(childContext, cancel, monitor)
	}

	monitorMsgChannel := monitors.MessageBroadcaster.Subscribe()
	defer monitors.MessageBroadcaster.Unsubscribe(monitorMsgChannel)

	for {
		select {
		case <-ctx.Done():
			logger.Uptime.Info().Msg("Uptime worker shutting down...")
			return
		case msg := <-monitorMsgChannel:
			if msg.Status == monitors.Created {
				childContext, cancel := context.WithCancel(ctx)
				go runMonitor(childContext, cancel, *msg.Monitor)
			}
		}
	}
}

func runMonitor(ctx context.Context, cancelSelf context.CancelFunc, monitor monitors.IMonitor) {
	defer cancelSelf()
	logger.Uptime.Debug().Dur("interval", monitor.GetInterval()).Msgf("Starting monitor runner: %s", monitor.GetName())

	monitorMsgChannel := monitors.MessageBroadcaster.Subscribe()
	defer monitors.MessageBroadcaster.Unsubscribe(monitorMsgChannel)

	tickerChangedChannel := make(chan struct{}, 1)
	monitorMutex := sync.Mutex{}
	tickerMutex := sync.Mutex{}

	ticker := time.NewTicker(monitor.GetInterval())
	defer ticker.Stop()

	// Validate the monitor before starting
	err := monitor.Validate()
	if err != nil {
		logger.Uptime.Error().Err(err).Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Error validating monitor")
		return
	}

	// Start a goroutine to listen for monitor updates and handle them asynchronously from the main monitor loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Uptime.Info().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Stopping monitor due to context cancellation")
				return
			case msg := <-monitorMsgChannel:
				monitorMutex.Lock()
				shouldExit := false
				func() {
					defer monitorMutex.Unlock()
					if msg.Id != monitor.GetId() {
						logger.Uptime.Trace().Str("id", msg.Id).Msgf("Ignoring message for monitor %s", monitor.GetName())
						return
					}

					logger.Uptime.Trace().Str("id", msg.Id).Msgf("Received message for monitor %s", monitor.GetName())

					switch msg.Status {
					case monitors.Edited:
						logger.Uptime.Debug().Str("name", monitor.GetName()).Str("id", monitor.GetId()).Msg("Updating monitor")
						oldInterval := monitor.GetInterval()
						monitor = *msg.Monitor
						if monitor.GetInterval() != oldInterval {
							logger.Uptime.Info().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Changing monitor interval to %s", monitor.GetInterval())
							tickerChangedChannel <- struct{}{}
						}

					case monitors.Deleted, monitors.Disabled:
						logger.Uptime.Info().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Stopping monitor due to deletion or disablement")
						cancelSelf()
						shouldExit = true
					}
				}()
				if shouldExit {
					logger.Uptime.Info().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Exiting monitor handler loop")
					return
				}
			}
		}
	}()

	for {
		select {
		case <-tickerChangedChannel:
			logger.Uptime.Debug().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Monitor interval changed, restarting ticker")
			tickerMutex.Lock()
			ticker.Reset(monitor.GetInterval())
			tickerMutex.Unlock()
			continue
		case <-ticker.C:
			func() {
				monitorMutex.Lock()
				defer monitorMutex.Unlock()

				err := monitor.Validate()
				if err != nil {
					logger.Uptime.Error().Err(err).Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Error validating monitor")
					return
				}

				logger.Uptime.Trace().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Running monitor")
				response := monitor.Run()
				logger.Uptime.Debug().Str("id", monitor.GetId()).Str("name", monitor.GetName()).Any("monitor_response", response).Msg("Monitor response")

				if response.IsError() {
					logger.Uptime.Error().Errs("errors", response.GetErrors()).Str("id", monitor.GetId()).Str("name", monitor.GetName()).Msgf("Monitor check failed")
				}

				//	// Optionally, you can save the result to the database or perform further actions
				//	err = db.SaveMonitorResult(monitor.GetId(), result)
				//	if err != nil {
				//		logger.Uptime.Error().Err(err).Msgf("Failed to save result for monitor %s", monitor.GetName())
				//	}
			}()
		case <-ctx.Done():
			return
		}
	}
}
