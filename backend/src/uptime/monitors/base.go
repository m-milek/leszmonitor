package monitors

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/teris-io/shortid"
)

type IMonitor interface {
	Run() (IMonitorResponse, error)
	GetId() string
	GetName() string
	GetDescription() string
	GetInterval() int
	validate() error
	validateBase() error
	setBase(base baseMonitor)
	GetType() MonitorType
}

type baseMonitor struct {
	Id          string      `json:"id" bson:"id"`                   // Unique identifier for the monitor
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

func NewBaseMonitor(name, description string, interval int, ownerId string, monitorType MonitorType) *baseMonitor {
	return &baseMonitor{
		Id:          generateMonitorId(),
		Name:        name,
		Description: description,
		Interval:    interval,
		OwnerId:     ownerId,
		Type:        monitorType,
	}
}

func validateBaseMonitor(b IMonitor) error {
	if b.GetName() == "" {
		return fmt.Errorf("monitor name cannot be empty")
	}
	if b.GetInterval() <= 0 {
		return fmt.Errorf("monitor interval must be greater than zero")
	}
	if b.GetType() == "" {
		return fmt.Errorf("monitor type cannot be empty")
	}
	if b.GetId() == "" {
		return fmt.Errorf("monitor ID cannot be empty")
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

	if base.Id == "" {
		log.Trace().Msg("Monitor ID is empty, generating a new one")
		base.Id = generateMonitorId()
	} else {
		log.Trace().Msgf("Monitor ID is set: %s", base.Id)
	}

	monitorData.setBase(base)

	return nil
}

func generateMonitorId() string {
	id, err := shortid.Generate()
	if err != nil {
		panic(fmt.Sprintf("Failed to generate monitor ID: %v", err))
	}
	return id
}

func (b baseMonitor) Run() (IMonitorResponse, error) {
	return nil, fmt.Errorf("run method not implemented for baseMonitor")
}

func (b baseMonitor) GetId() string {
	return b.Id
}

func (b baseMonitor) GetName() string {
	return b.Name
}

func (b baseMonitor) GetDescription() string {
	return b.Description
}

func (b baseMonitor) GetInterval() int {
	return b.Interval
}

func (b baseMonitor) validate() error {
	return validateBaseMonitor(b)
}

func (b baseMonitor) validateBase() error {
	return validateBaseMonitor(b)
}

func (b baseMonitor) setBase(base baseMonitor) {
	return // baseMonitor does not have a setBase method, this is a no-op
}

func (b baseMonitor) GetType() MonitorType {
	return b.Type
}
