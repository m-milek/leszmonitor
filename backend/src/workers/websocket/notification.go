package websocket

import "github.com/m-milek/leszmonitor/models/monitors"

var (
	notificationTypeMonitorRun = "monitor_run"
)

type iNotification interface {
	GetType() string
	GetPayload() any
}

type baseNotification struct {
	Type string `json:"type"`
}

type monitorRunNotification struct {
	baseNotification
	MonitorID string                  `json:"monitorId"`
	Response  monitors.IMonitorResult `json:"response"`
}

func newMonitorRunNotification(result monitors.IMonitorResult, monitor monitors.IMonitor) *monitorRunNotification {
	return &monitorRunNotification{
		baseNotification: baseNotification{
			Type: notificationTypeMonitorRun,
		},
		Response:  result,
		MonitorID: monitor.GetID().String(),
	}
}

func (n *monitorRunNotification) GetType() string {
	return n.Type
}

func (n *monitorRunNotification) GetPayload() any {
	return n.Response
}
