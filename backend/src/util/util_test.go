package util

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestSliceContains(t *testing.T) {
	setupTest := func() []string {
		array := []string{"apple", "banana", "cherry"}
		return array
	}

	array := setupTest()

	assert.True(t, SliceContains(array, "banana"))
	assert.False(t, SliceContains(array, "orange"))
}

func TestSliceMinMax(t *testing.T) {
	t.Run("Test SliceMinMax with integers", func(t *testing.T) {
		slice := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
		min, max := SliceMinMax(slice)

		assert.Equal(t, 1, min)
		assert.Equal(t, 9, max)
	})

	t.Run("Test SliceMinMax with empty slice", func(t *testing.T) {
		emptySlice := []int{}
		minEmpty, maxEmpty := SliceMinMax(emptySlice)
		assert.Equal(t, 0, minEmpty)
		assert.Equal(t, 0, maxEmpty)
	})

	t.Run("Test SliceMinMax with strings", func(t *testing.T) {
		slice := []string{"apple", "banana", "cherry"}
		min, max := SliceMinMax(slice)

		assert.Equal(t, "apple", min)
		assert.Equal(t, "cherry", max)
	})
}

func TestGetUnixTimestamp(t *testing.T) {
	timestamp := GetUnixTimestamp()
	assert.Greater(t, timestamp, int64(0), "Unix timestamp should be greater than 0")
	assert.GreaterOrEqual(t, len(strconv.FormatInt(timestamp, 10)), 10, "Unix timestamp should be at least 10 digits long")
}

func TestGetUnixTimestampMillis(t *testing.T) {
	timestamp := GetUnixTimestampMillis()
	assert.Greater(t, timestamp, int64(0), "Unix timestamp in milliseconds should be greater than 0")
	assert.GreaterOrEqual(t, len(strconv.FormatInt(timestamp, 10)), 13, "Unix timestamp in milliseconds should be at least 13 digits long")
}
