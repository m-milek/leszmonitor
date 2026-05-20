package monitors

import (
	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/monitorresult"
)

type MonitorLifecycleState string

const (
	Created MonitorLifecycleState = "created"
	Edited  MonitorLifecycleState = "edited"
	Deleted MonitorLifecycleState = "deleted"
)

type MonitorLifecycleMessage struct {
	ID      uuid.UUID
	Status  MonitorLifecycleState
	Monitor *Monitor
}

type MonitorRunMessage struct {
	Monitor Monitor
	Result  monitorresult.IMonitorResult
}
