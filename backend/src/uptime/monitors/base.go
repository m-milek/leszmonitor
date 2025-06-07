package monitors

import (
	"fmt"
	"github.com/teris-io/shortid"
)

type IMonitor interface {
	Run() IMonitorResponse
	Validate() error
	GetId() string
	GetName() string
	GetDescription() string
	GetInterval() int
	GetType() MonitorConfigType
	GenerateId()
}

type IMonitorConfig interface {
	run() IMonitorResponse
	validate() error
}

type BaseMonitor struct {
	Id          string            `json:"id" bson:"id"`                   // Unique identifier for the monitor
	Name        string            `json:"name" bson:"name"`               // Name of the monitor
	Description string            `json:"description" bson:"description"` // Description of the monitor
	Interval    int               `json:"interval" bson:"interval"`       // How often to run the monitor in seconds
	OwnerId     string            `json:"owner_id" bson:"owner_id"`       // ID of the owner of the monitor
	Type        MonitorConfigType `json:"type" bson:"type"`               // Type of the monitor (http, ping, etc.)
}

type MonitorConfigType string

const (
	Http MonitorConfigType = "http"
	Ping MonitorConfigType = "ping"
)

func NewMonitor(name, description string, interval int, ownerId string, monitorType MonitorConfigType) *BaseMonitor {
	return &BaseMonitor{
		Id:          generateMonitorId(),
		Name:        name,
		Description: description,
		Interval:    interval,
		OwnerId:     ownerId,
		Type:        monitorType,
	}
}

func (m *BaseMonitor) Validate() error {
	if err := m.validateBase(); err != nil {
		return fmt.Errorf("monitor validation failed: %w", err)
	}
	return nil
}

//func (m *BaseMonitor) Run() (IMonitorResponse, error) {
//	if err := m.Validate(); err != nil {
//		return nil, fmt.Errorf("monitor validation failed: %w", err)
//	}
//
//	response := m.Config.run()
//
//	if len(response.GetErrors()) > 0 {
//		logger.Uptime.Debug().Any("errors", response.GetErrors()).Msg("BaseMonitor run encountered an error")
//	}
//
//	return response, nil
//}

//func (m *BaseMonitor) UnmarshalJSON(data []byte) error {
//	// Define an alias type to avoid the recursion
//	type alias BaseMonitor
//
//	// Create a map to hold the raw JSON
//	var rawMap map[string]json.RawMessage
//	if err := json.Unmarshal(data, &rawMap); err != nil {
//		return fmt.Errorf("failed to parse monitor JSON: %w", err)
//	}
//
//	// Extract the type field to determine the monitor type
//	var monitorType MonitorConfigType
//	if typeData, exists := rawMap["type"]; exists {
//		if err := json.Unmarshal(typeData, &monitorType); err != nil {
//			return fmt.Errorf("failed to parse monitor type: %w", err)
//		}
//	} else {
//		return fmt.Errorf("missing type field in monitor JSON")
//	}
//
//	// Create the appropriate config type
//	config := MapMonitorType(monitorType)
//	if config == nil {
//		return fmt.Errorf("unknown monitor type: %s", monitorType)
//	}
//
//	// Extract the config data
//	if configData, exists := rawMap["config"]; exists {
//		if err := json.Unmarshal(configData, &config); err != nil {
//			return fmt.Errorf("failed to unmarshal monitor configuration: %w", err)
//		}
//	} else {
//		return fmt.Errorf("missing config field in monitor JSON")
//	}
//
//	// Remove the config field from the raw map to avoid unmarshaling it twice
//	delete(rawMap, "config")
//
//	// Re-encode the remaining fields
//	modifiedData, err := json.Marshal(rawMap)
//	if err != nil {
//		return fmt.Errorf("failed to re-encode monitor data: %w", err)
//	}
//
//	// Unmarshal into the alias type (which will not trigger this custom UnmarshalJSON)
//	aux := (*alias)(m)
//	if err := json.Unmarshal(modifiedData, aux); err != nil {
//		return fmt.Errorf("failed to unmarshal monitor base fields: %w", err)
//	}
//
//	// Set the config
//	m.Config = config
//
//	return nil
//}

//func (m *BaseMonitor) UnmarshalBSON() ([]byte, error) {
//	return nil, fmt.Errorf("UnmarshalBSON is not implemented yet")
//}

func generateMonitorId() string {
	id, err := shortid.Generate()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate monitor ID: %v", err))
	}
	return id
}

func (m *BaseMonitor) validateBase() error {
	if m.GetName() == "" {
		return fmt.Errorf("monitor name cannot be empty")
	}
	if m.GetInterval() <= 0 {
		return fmt.Errorf("monitor interval must be greater than zero")
	}
	if m.GetType() == "" {
		return fmt.Errorf("monitor type cannot be empty")
	}
	if m.GetId() == "" {
		return fmt.Errorf("monitor ID cannot be empty")
	}
	return nil
}

func (m *BaseMonitor) GenerateId() {
	if m.Id == "" {
		m.Id = generateMonitorId()
	}
}

func (m *BaseMonitor) GetId() string {
	return m.Id
}

func (m *BaseMonitor) GetName() string {
	return m.Name
}

func (m *BaseMonitor) GetDescription() string {
	return m.Description
}

func (m *BaseMonitor) GetInterval() int {
	return m.Interval
}

func (m *BaseMonitor) GetType() MonitorConfigType {
	return m.Type
}

type MonitorTypeExtractor struct {
	Type MonitorConfigType `json:"type"`
}
