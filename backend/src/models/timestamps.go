package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

type Timestamps struct {
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RawTimestamps struct {
	CreatedAt pgtype.Timestamptz
	UpdatedAt pgtype.Timestamptz
}

func (t *Timestamps) SetTimestamps(timestamps RawTimestamps) {
	t.CreatedAt = timestamps.CreatedAt.Time
	t.UpdatedAt = timestamps.UpdatedAt.Time
}
