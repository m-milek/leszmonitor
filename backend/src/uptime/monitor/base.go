package monitors

import (
	"fmt"
	"github.com/m-milek/leszmonitor/util"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type IMonitor interface {
	Run() IMonitorResponse
	Validate() error
	GetId() string
	GetObjectId() bson.ObjectID
	GetName() string
	GetDescription() string
	GetInterval() time.Duration
	GetType() MonitorConfigType
	GenerateId()
	SetGroupId(string)
}

type IMonitorConfig interface {
	run() IMonitorResponse
	validate() error
}

type BaseMonitor struct {
	ObjectId    bson.ObjectID     `json:"-" bson:"_id"`
	Id          string            `json:"id" bson:"id"`                   // Unique identifier for the monitor
	Name        string            `json:"name" bson:"name"`               // Name of the monitor
	Description string            `json:"description" bson:"description"` // Description of the monitor
	Interval    int               `json:"interval" bson:"interval"`       // How often to run the monitor in seconds
	GroupId     string            `json:"groupId" bson:"groupId"`         // ID of the owner group of the monitor
	Type        MonitorConfigType `json:"type" bson:"type"`               // Type of the monitor (http, ping, etc.)
}

type MonitorConfigType string

const (
	Http MonitorConfigType = "http"
	Ping MonitorConfigType = "ping"
)

type MonitorTypeExtractor struct {
	Type MonitorConfigType `json:"type"`
}

func (m *BaseMonitor) Validate() error {
	if err := m.validateBase(); err != nil {
		return fmt.Errorf("monitor validation failed: %w", err)
	}
	return nil
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
	if m.GroupId == "" {
		return fmt.Errorf("monitor group ID cannot be empty")
	}
	return nil
}

func (m *BaseMonitor) GenerateId() {
	m.Id = util.IdFromString(m.GetName())
}

func (m *BaseMonitor) GetId() string {
	return m.Id
}

func (m *BaseMonitor) GetObjectId() bson.ObjectID {
	return m.ObjectId
}

func (m *BaseMonitor) GetName() string {
	return m.Name
}

func (m *BaseMonitor) GetDescription() string {
	return m.Description
}

func (m *BaseMonitor) GetInterval() time.Duration {
	return time.Duration(m.Interval) * time.Second
}

func (m *BaseMonitor) GetType() MonitorConfigType {
	return m.Type
}

func (m *BaseMonitor) SetGroupId(groupId string) {
	m.GroupId = groupId
}
