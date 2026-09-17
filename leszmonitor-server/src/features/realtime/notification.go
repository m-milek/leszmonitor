package realtime

import (
	"github.com/m-milek/leszmonitor/features/monitors"
)

var (
	notificationTypeMonitorRun = "monitor_run"
)

type baseNotification struct {
	Type string `json:"type"`
}

type monitorRunNotification struct {
	baseNotification

	MonitorID string                  `json:"monitorId"`
	Response  monitors.IMonitorResult `json:"response"`
}

func newMonitorRunNotification(result monitors.IMonitorResult, monitor monitors.Monitor) *monitorRunNotification {
	return &monitorRunNotification{
		baseNotification: baseNotification{
			Type: notificationTypeMonitorRun,
		},
		Response:  result,
		MonitorID: monitor.ID.String(),
	}
}
