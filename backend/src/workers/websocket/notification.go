package websocket

import (
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/models/monitors"
)

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
	MonitorID string                       `json:"monitorId"`
	Response  monitorresult.IMonitorResult `json:"response"`
}

func newMonitorRunNotification(result monitorresult.IMonitorResult, monitor monitors.Monitor) *monitorRunNotification {
	return &monitorRunNotification{
		baseNotification: baseNotification{
			Type: notificationTypeMonitorRun,
		},
		Response:  result,
		MonitorID: monitor.ID.String(),
	}
}

func (n *monitorRunNotification) GetType() string {
	return n.Type
}

func (n *monitorRunNotification) GetPayload() any {
	return n.Response
}
