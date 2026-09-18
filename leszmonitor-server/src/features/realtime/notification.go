package realtime

import (
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/results"
)

var (
	notificationTypeMonitorRun = "monitor_run"
)

type baseNotification struct {
	Type string `json:"type"`
}

type monitorRunNotification struct {
	baseNotification

	MonitorID string                 `json:"monitorId"`
	Response  results.IMonitorResult `json:"response"`
}

func newMonitorRunNotification(result results.IMonitorResult, monitor monitors.Monitor) *monitorRunNotification {
	return &monitorRunNotification{
		baseNotification: baseNotification{
			Type: notificationTypeMonitorRun,
		},
		Response:  result,
		MonitorID: monitor.ID.String(),
	}
}
