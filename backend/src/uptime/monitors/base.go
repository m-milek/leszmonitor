package monitors

import (
	"encoding/json"
	"fmt"
)

type IMonitor interface {
	Run() (IMonitorResponse, error)
	GetName() string
	GetDescription() string
	GetInterval() int
	validate() error
	setBase(base baseMonitor)
}

type baseMonitor struct {
	Name        string      `json:"name" bson:"name"`               // Name of the monitor
	Description string      `json:"description" bson:"description"` // Description of the monitor
	Interval    int         `json:"interval" bson:"interval"`       // How often to run the monitor in seconds
	OwnerId     string      `json:"owner_id" bson:"owner_id"`       // ID of the owner of the monitor
	Type        MonitorType `json:"type" bson:"type"`               // Type of the monitor (http, ping, etc.)
}

type MonitorType string

const (
	Http MonitorType = "http"
	Ping MonitorType = "ping"
)

func (m *baseMonitor) validate() error {
	if m.Name == "" {
		return fmt.Errorf("monitor name cannot be empty")
	}
	if m.Interval <= 0 {
		return fmt.Errorf("monitor interval must be greater than zero")
	}
	if m.Type == "" {
		return fmt.Errorf("monitor type cannot be empty")
	}
	return nil
}

func UnmarshalMonitor(rawData []byte, monitorData IMonitor) error {
	var base baseMonitor
	if err := json.Unmarshal(rawData, &base); err != nil {
		return err
	}
	if err := json.Unmarshal(rawData, &monitorData); err != nil {
		return err
	}
	monitorData.setBase(base)

	return nil
}
