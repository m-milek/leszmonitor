package monitors

import "github.com/jackc/pgx/v5/pgtype"

type monitorLifecycleState string

const (
	Created monitorLifecycleState = "created"
	Edited  monitorLifecycleState = "edited"
	Deleted monitorLifecycleState = "deleted"
	Stopped monitorLifecycleState = "stopped"
	Started monitorLifecycleState = "started"
)

type MonitorLifecycleMessage struct {
	ID      pgtype.UUID
	Status  monitorLifecycleState
	Monitor *IConcreteMonitor
}

type MonitorRunMessage struct {
	Result  *IMonitorResponse
	Monitor *IMonitor
}
