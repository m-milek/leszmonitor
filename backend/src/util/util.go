package util

import (
	"golang.org/x/exp/constraints"
	"time"
)

func SliceContains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func SliceMinMax[T constraints.Ordered](slice []T) (min T, max T) {
	if len(slice) == 0 {
		return
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
	return
}

func GetUnixTimestamp() int64 {
	return time.Now().Unix()
}

func GetUnixTimestampMillis() int64 {
	return time.Now().UnixMilli()
}
