package statuschange

import (
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
)

type MonitorStatusChange struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	MonitorID      uuid.UUID  `json:"monitorId" db:"monitor_id"`
	CausedByID     *uuid.UUID `json:"causedById" db:"caused_by_id"`
	PreviousStatus string     `json:"previousStatus" db:"previous_status"`
	NextStatus     string     `json:"nextStatus" db:"next_status"`

	CreatedAt time.Time `json:"createdAt" db:"created_at"`
}

// MonitorStatusPeriod is a span of time during which a monitor had one status. EndedAt is nil while it is ongoing.
type MonitorStatusPeriod struct {
	Status          kind.MonitorStatus `json:"status"`
	StartedAt       time.Time          `json:"startedAt"`
	EndedAt         *time.Time         `json:"endedAt"`
	DurationSeconds int64              `json:"durationSeconds"`
}
