package util

import (
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestUTCTimestamptz_MarshalJSON_ReturnsNullForInvalid(t *testing.T) {
	val := utcTimestamptz{Timestamptz: pgtype.Timestamptz{Valid: false}}

	b, err := val.MarshalJSON()

	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), b)
}

func TestUTCTimestamptz_MarshalJSON_UTCAndRFC3339(t *testing.T) {
	// Time with non-UTC zone should be converted to UTC and formatted as RFC3339
	loc := time.FixedZone("UTC+1", 1*60*60)
	in := time.Date(2024, 12, 31, 23, 0, 0, 0, loc)
	expectedUTC := in.UTC().Format(time.RFC3339)

	val := utcTimestamptz{Timestamptz: pgtype.Timestamptz{Time: in, Valid: true}}

	b, err := val.MarshalJSON()

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("\"%s\"", expectedUTC), string(b))
}

func TestTimestamps_MarshalJSON_FieldTagsAndNestedMarshaler(t *testing.T) {
	// createdAt should be RFC3339 UTC string, updatedAt should be null
	loc := time.FixedZone("UTC+2", 2*60*60)
	created := time.Date(2025, 1, 2, 3, 4, 5, 0, loc)
	expectedCreated := created.UTC().Format(time.RFC3339)

	ts := Timestamps{
		CreatedAt: utcTimestamptz{Timestamptz: pgtype.Timestamptz{Time: created, Valid: true}},
		UpdatedAt: utcTimestamptz{Timestamptz: pgtype.Timestamptz{Valid: false}},
	}

	b, err := json.Marshal(ts)
	assert.NoError(t, err)

	expected := fmt.Sprintf("{\"createdAt\":\"%s\",\"updatedAt\":null}", expectedCreated)
	assert.Equal(t, expected, string(b))
}
