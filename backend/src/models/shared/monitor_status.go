package shared

type MonitorStatus string

const (
	MonitorStatusUp          MonitorStatus = "up"
	MonitorStatusDown        MonitorStatus = "down"
	MonitorStatusPaused      MonitorStatus = "paused"
	MonitorStatusError       MonitorStatus = "error"
	MonitorStatusMaintenance MonitorStatus = "maintenance"
)
