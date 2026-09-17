package monitors

type LatencyStats struct {
	Avg float64 `json:"avg"`
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

type StatusChangeStats struct {
	SecondsInCurrentStatus int64 `json:"secondsInCurrentStatus"`
}

type UptimeStats struct {
	UptimePercentage float64 `json:"uptimePercentage"`
}

type MonitorStats struct {
	Latency      LatencyStats      `json:"latency"`
	StatusChange StatusChangeStats `json:"statusChange"`
	Uptime       UptimeStats       `json:"uptime"`
}
