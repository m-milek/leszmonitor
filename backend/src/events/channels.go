package events

import "github.com/m-milek/leszmonitor/models/monitors"

// MonitorLifecycleChannel distributes monitor lifecycle events to subscribers.
var MonitorLifecycleChannel = newBroadcaster[monitors.MonitorLifecycleMessage]()

// MonitorRunChannel distributes monitor run events (e.g., tcp results) to subscribers.
var MonitorRunChannel = newBroadcaster[monitors.MonitorRunMessage]()
