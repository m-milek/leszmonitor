package workers

import (
	"context"
	"sync"
	"time"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models/monitors"
)

// StartUptimeWorker initializes the uptime worker that periodically checks all monitors.
// It retrieves all monitors from the database and starts a goroutine for each monitor to run its checks.
// It also listens for monitor creation messages to start new monitors dynamically.
// The worker runs until the context is done, allowing for graceful shutdown.
func StartUptimeWorker(ctx context.Context) {
	log.Uptime.Info().Msg("Starting uptime worker...")

	allMonitors, err := db.Get().Monitors().GetAllMonitors(ctx)

	if err != nil {
		log.Uptime.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return
	}
	log.Uptime.Info().Any("monitors", allMonitors).Msgf("Found %d monitors to check", len(allMonitors))

	for _, monitor := range allMonitors {
		childContext, cancel := context.WithCancel(ctx)
		go runMonitor(childContext, cancel, monitor)
	}

	monitorMsgChannel := events.MonitorLifecycleChannel.Subscribe()
	defer events.MonitorLifecycleChannel.Unsubscribe(monitorMsgChannel)

	for {
		select {
		case <-ctx.Done():
			log.Uptime.Info().Msg("Uptime worker shutting down...")
			return
		case msg := <-monitorMsgChannel:
			if msg.Status == monitors.Created {
				childContext, cancel := context.WithCancel(ctx)
				go runMonitor(childContext, cancel, *msg.Monitor)
			}
		}
	}
}

func runMonitor(ctx context.Context, cancelSelf context.CancelFunc, monitor monitors.Monitor) {
	defer cancelSelf()
	log.Uptime.Debug().Int("interval", monitor.Interval).Msgf("Starting monitor runner: %s", monitor.Name)

	monitorMsgChannel := events.MonitorLifecycleChannel.Subscribe()
	defer events.MonitorLifecycleChannel.Unsubscribe(monitorMsgChannel)

	tickerChangedChannel := make(chan struct{}, 1)
	monitorMutex := sync.Mutex{}
	tickerMutex := sync.Mutex{}

	duration := time.Duration(monitor.Interval) * time.Second
	ticker := time.NewTicker(duration)
	defer ticker.Stop()

	// Validate the monitor before starting
	err := monitor.Validate()
	if err != nil {
		log.Uptime.Error().Err(err).Str("id", monitor.Slug).Str("name", monitor.Name).Msgf("Error validating monitor")
		return
	}

	// Start a goroutine to listen for monitor updates and handle them asynchronously from the main monitor loop
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Uptime.Info().Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Stopping monitor due to context cancellation")
				return
			case msg := <-monitorMsgChannel:
				monitorMutex.Lock()
				shouldExit := false
				func() {
					defer monitorMutex.Unlock()
					if msg.ID != monitor.ID {
						log.Uptime.Trace().Str("id", msg.ID.String()).Msgf("Ignoring message for monitor %s", monitor.Name)
						return
					}

					log.Uptime.Trace().Str("id", msg.ID.String()).Msgf("Received message for monitor %s", monitor.Name)

					switch msg.Status {
					case monitors.Edited:
						// refetch the monitor
						newMonitor, err := db.Get().Monitors().GetMonitorByID(ctx, msg.ID)
						if err != nil {
							log.Uptime.Error().Err(err).Str("id", msg.ID.String()).Msg("Failed to refetch monitor after edit")
							return
						}
						log.Uptime.Debug().Str("name", monitor.Name).Str("id", monitor.ID.String()).Msg("Updating monitor")
						oldInterval := monitor.Interval
						monitor = *newMonitor
						if monitor.Interval != oldInterval {
							log.Uptime.Info().Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Changing monitor interval to %d", monitor.Interval)
							tickerChangedChannel <- struct{}{}
						}

					case monitors.Deleted, monitors.Stopped:
						log.Uptime.Info().Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Stopping monitor due to deletion or disablement")
						cancelSelf()
						shouldExit = true
					}
				}()
				if shouldExit {
					log.Uptime.Info().Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Exiting monitor handler loop")
					return
				}
			}
		}
	}()

	for {
		select {
		case <-tickerChangedChannel:
			log.Uptime.Debug().Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Monitor interval changed, restarting ticker")
			tickerMutex.Lock()
			newDuration := time.Duration(monitor.Interval) * time.Second
			ticker.Reset(newDuration)
			tickerMutex.Unlock()
			continue
		case <-ticker.C:
			func() {
				monitorMutex.Lock()
				defer monitorMutex.Unlock()

				err := monitor.Validate()
				if err != nil {
					log.Uptime.Error().Err(err).Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Error validating monitor")
					return
				}

				log.Uptime.Trace().Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Running monitor")
				probe, err := monitors.UnmarshalProbeFromBytes(monitor.Type, []byte(monitor.ProbeConfig))
				if err != nil {
					log.Uptime.Error().Err(err).Str("id", monitor.ID.String()).Str("name", monitor.Name).Msgf("Error unmarshalling probe config")
					return
				}
				result := probe.Run(monitor.ID)
				log.Uptime.Debug().Str("id", monitor.ID.String()).Str("name", monitor.Name).Any("monitor_result", result).Msg("Monitor result")

				events.MonitorRunChannel.Broadcast(monitors.MonitorRunMessage{
					Result:  result,
					Monitor: monitor,
				})

				if !result.GetIsSuccess() {
					errDetails := result.GetErrorDetails()
					if errDetails.ErrorMessage != "" || len(errDetails.Errors) > 0 {
						log.Uptime.Error().
							Str("error_message", errDetails.ErrorMessage).
							Strs("errors", errDetails.Errors).
							Str("id", monitor.ID.String()).
							Str("name", monitor.Name).
							Msg("Monitor encountered internal error")
					}
					if len(errDetails.Failures) > 0 {
						log.Uptime.Warn().
							Strs("failures", errDetails.Failures).
							Str("id", monitor.ID.String()).
							Str("name", monitor.Name).
							Msg("Monitor check failed (service down or misconfigured)")
					}
				}

				_, err = db.Get().MonitorResults().InsertMonitorResult(ctx, result)
				if err != nil {
					return
				}
			}()
		case <-ctx.Done():
			return
		}
	}
}
