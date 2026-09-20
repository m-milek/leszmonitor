package downtime

import (
	"time"

	"github.com/google/uuid"
)

type HeartbeatRecord struct {
	BeatAt time.Time `db:"beat_at"`
}

type AppDowntime struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StartedAt time.Time `json:"startedAt" db:"started_at"`
	EndedAt   time.Time `json:"endedAt" db:"ended_at"`
}
