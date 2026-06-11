package services

import "github.com/m-milek/leszmonitor/models/monitors"

type MonitorEventPublisher interface {
	PublishLifecycle(msg monitors.MonitorLifecycleMessage)
}
