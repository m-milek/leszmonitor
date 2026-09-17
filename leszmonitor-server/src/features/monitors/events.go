package monitors

import "github.com/m-milek/leszmonitor/platform/events"

// MonitorLifecycleChannel distributes monitor lifecycle events to subscribers.
var MonitorLifecycleChannel = events.NewEventBus[MonitorLifecycleMessage]("monitor_lifecycle")

// MonitorRunChannel distributes monitor run events (e.g., tcp results) to subscribers.
var MonitorRunChannel = events.NewEventBus[MonitorRunMessage]("monitor_run")

type BroadcastMonitorPublisher struct{}

func (BroadcastMonitorPublisher) PublishLifecycle(msg MonitorLifecycleMessage) {
	MonitorLifecycleChannel.Broadcast(msg)
}
