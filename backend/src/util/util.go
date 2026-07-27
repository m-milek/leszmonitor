package util

import (
	"slices"

	"cmp"
	"time"
)

func SliceContains[T comparable](slice []T, item T) bool {
	return slices.Contains(slice, item)
}

func SliceMinMax[T cmp.Ordered](slice []T) (min T, max T) {
	if len(slice) == 0 {
		return min, max
	}
	min, max = slice[0], slice[0]
	for _, v := range slice {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return min, max
}

func GetUnixTimestamp() int64 {
	return time.Now().Unix()
}

func GetUnixTimestampMillis() int64 {
	return time.Now().UnixMilli()
}
