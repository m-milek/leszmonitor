package stats

import (
	"github.com/m-milek/leszmonitor/features/monitors/kind"
)

type LatencyStats struct {
	Avg float64 `json:"avg"`
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type StatusChangeStats struct {
	SecondsInCurrentStatus int64 `json:"secondsInCurrentStatus"`
}

type UptimeStats struct {
	StatusToCount      map[kind.MonitorStatus]int     `json:"statusToCount"`
	StatusToPercentage map[kind.MonitorStatus]float64 `json:"statusToPercentage"`
}

type MonitorStats struct {
	Latency      LatencyStats      `json:"latency"`
	StatusChange StatusChangeStats `json:"statusChange"`
	Uptime       UptimeStats       `json:"uptime"`
}
