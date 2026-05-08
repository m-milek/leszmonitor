package events

import "github.com/m-milek/leszmonitor/models/monitors"

// MonitorLifecycleChannel distributes monitor lifecycle events to subscribers.
var MonitorLifecycleChannel = newEventBus[monitors.MonitorLifecycleMessage]("monitor_lifecycle")

// MonitorRunChannel distributes monitor run events (e.g., tcp results) to subscribers.
var MonitorRunChannel = newEventBus[monitors.MonitorRunMessage]("monitor_run")
