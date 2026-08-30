package models

import "github.com/google/uuid"

type MonitorStatusChange struct {
	ID             uuid.UUID `json:"id" db:"id"`
	MonitorID      uuid.UUID `json:"monitorId" db:"monitor_id"`
	CausedByID     uuid.UUID `json:"causedById" db:"caused_by_id"`
	PreviousStatus string    `json:"previousStatus" db:"previous_status"`
	NextStatus     string    `json:"nextStatus" db:"next_status"`

	CreatedAt string `json:"createdAt" db:"created_at"`
}
