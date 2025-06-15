package monitors

import "github.com/m-milek/leszmonitor/common"

type MonitorMessageStatus string

const (
	Created  MonitorMessageStatus = "created"
	Edited   MonitorMessageStatus = "edited"
	Deleted  MonitorMessageStatus = "deleted"
	Disabled MonitorMessageStatus = "disabled"
	Enabled  MonitorMessageStatus = "enabled"
)

type MonitorMessage struct {
	Id      string
	Status  MonitorMessageStatus
	Monitor *IMonitor
}

// MessageBroadcaster is a global channel for sending monitor messages.
var MessageBroadcaster = common.NewBroadcaster[MonitorMessage]()
