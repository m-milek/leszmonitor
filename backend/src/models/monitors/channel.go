package monitors

import (
	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/monitorresult"
)

type monitorLifecycleState string

const (
	Created monitorLifecycleState = "created"
	Edited  monitorLifecycleState = "edited"
	Deleted monitorLifecycleState = "deleted"
	Stopped monitorLifecycleState = "stopped"
	Started monitorLifecycleState = "started"
)

type MonitorLifecycleMessage struct {
	ID      uuid.UUID
	Status  monitorLifecycleState
	Monitor *IConcreteMonitor
}

type MonitorRunMessage struct {
	Result  monitorresult.IMonitorResult
	Monitor IMonitor
}
