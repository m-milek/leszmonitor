package monitors

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/common"
)

type MonitorMessageStatus string

const (
	Created MonitorMessageStatus = "created"
	Edited  MonitorMessageStatus = "edited"
	Deleted MonitorMessageStatus = "deleted"
	Stopped MonitorMessageStatus = "stopped"
	Started MonitorMessageStatus = "started"
)

type MonitorMessage struct {
	Id      pgtype.UUID
	Status  MonitorMessageStatus
	Monitor *IConcreteMonitor
}

// MessageBroadcaster is a global channel for sending monitor messages.
var MessageBroadcaster = common.NewBroadcaster[MonitorMessage]()
