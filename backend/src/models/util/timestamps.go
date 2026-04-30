package util

import (
	"encoding/json"
	"time"
)

// Timestamps holds createdAt and updatedAt timestamps for database models.
type Timestamps struct {
	CreatedAt UTCTime  `json:"createdAt" db:"created_at"`
	UpdatedAt *UTCTime `json:"updatedAt" db:"updated_at"`
}

// UTCTime is a wrapper around time.Time that marshals time in UTC format.
type UTCTime struct {
	time.Time
}

// MarshalJSON implements the json.Marshaler interface for UTCTime.
// It converts the time to UTC before marshaling to JSON.
func (t UTCTime) MarshalJSON() ([]byte, error) {
	if t.Time.IsZero() {
		return []byte("null"), nil
	}
	utcTime := t.Time.UTC()
	return json.Marshal(utcTime.Format(time.RFC3339))
}
