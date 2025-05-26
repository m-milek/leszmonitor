package monitors

import (
	"fmt"
)

type IMonitor interface {
	Run(client httpClient) (IMonitorResponse, error)
	GetName() string
	GetDescription() string
	GetInterval() int
	GetTimeout() int
}

type baseMonitor struct {
	Name        string      `json:"name" bson:"name"`
	Description string      `json:"description" bson:"description"`
	Interval    int         `json:"interval" bson:"interval"` // in seconds
	Timeout     int         `json:"timeout" bson:"timeout"`   // in seconds
	OwnerId     string      `json:"owner_id" bson:"owner_id"`
	Type        MonitorType `json:"type" bson:"type"`
}

type MonitorType string

const (
	Http MonitorType = "http"
)

func (m *baseMonitor) validate() error {
	if m.Name == "" {
		return fmt.Errorf("monitor name cannot be empty")
	}
	if m.Interval <= 0 {
		return fmt.Errorf("monitor interval must be greater than zero")
	}
	if m.Timeout <= 0 {
		return fmt.Errorf("monitor timeout must be greater than zero")
	}
	if m.Type == "" {
		return fmt.Errorf("monitor type cannot be empty")
	}
	return nil
}
