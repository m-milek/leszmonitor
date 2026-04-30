package util

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimestamps_MarshalJSON_FieldTagsAndNullableUpdatedAt(t *testing.T) {
	created := time.Date(2025, 1, 2, 3, 4, 5, 0, time.UTC)
	expectedCreated := created.Format(time.RFC3339)

	ts := Timestamps{
		CreatedAt: created,
		UpdatedAt: nil,
	}

	b, err := json.Marshal(ts)
	assert.NoError(t, err)

	expected := fmt.Sprintf("{\"createdAt\":\"%s\",\"updatedAt\":null}", expectedCreated)
	assert.Equal(t, expected, string(b))
}
