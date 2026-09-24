package monitors

import "github.com/m-milek/leszmonitor/platform/events"

// MonitorLifecycleChannel distributes monitor lifecycle events to subscribers.
var MonitorLifecycleChannel = events.NewEventBus[MonitorLifecycleMessage]("monitor_lifecycle")

// MonitorRunChannel distributes monitor run events (e.g., tcp results) to subscribers.
var MonitorRunChannel = events.NewEventBus[MonitorRunMessage]("monitor_run")

// MonitorExecuteChannel carries scheduled monitor checks from the scheduler to the executor.
var MonitorExecuteChannel = events.NewEventBus[MonitorExecuteMessage]("monitor_execute")

type BroadcastMonitorPublisher struct{}

func (BroadcastMonitorPublisher) PublishLifecycle(msg MonitorLifecycleMessage) {
	MonitorLifecycleChannel.Broadcast(msg)
}
