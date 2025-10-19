package monitors

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/common"
)

type monitorMessageStatus string

const (
	Created monitorMessageStatus = "created"
	Edited  monitorMessageStatus = "edited"
	Deleted monitorMessageStatus = "deleted"
	Stopped monitorMessageStatus = "stopped"
	Started monitorMessageStatus = "started"
)

type MonitorMessage struct {
	ID      pgtype.UUID
	Status  monitorMessageStatus
	Monitor *IConcreteMonitor
}

// MessageBroadcaster is a global channel for sending monitor messages.
var MessageBroadcaster = common.NewBroadcaster[MonitorMessage]()
