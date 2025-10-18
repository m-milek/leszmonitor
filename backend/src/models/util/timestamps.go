package util

import (
	"encoding/json"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

// Timestamps holds createdAt and updatedAt timestamps for database models.
type Timestamps struct {
	CreatedAt utcTimestamptz `json:"createdAt"`
	UpdatedAt utcTimestamptz `json:"updatedAt"`
}

// utcTimestamptz is a wrapper around pgtype.Timestamptz that marshals time in UTC format.
type utcTimestamptz struct {
	pgtype.Timestamptz
}

// MarshalJSON implements the json.Marshaler interface for utcTimestamptz.
// It converts the time to UTC before marshaling to JSON.
func (t utcTimestamptz) MarshalJSON() ([]byte, error) {
	if !t.Valid {
		return []byte("null"), nil
	}
	utcTime := t.Time.UTC()
	return json.Marshal(utcTime.Format(time.RFC3339))
}
