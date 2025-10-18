package monitors

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	util2 "github.com/m-milek/leszmonitor/models/util"
	"github.com/m-milek/leszmonitor/util"
	"time"
)

type IMonitor interface {
	Run() IMonitorResponse
	Validate() error
	GetId() pgtype.UUID
	GetDisplayId() string
	GenerateDisplayId()
	GetTeamId() pgtype.UUID
	SetTeamId(uuid pgtype.UUID)
	GetGroupId() pgtype.UUID
	GetName() string
	GetDescription() string
	GetInterval() time.Duration
	GetType() MonitorConfigType
	SetGroupId(uuid pgtype.UUID)
}

type IConcreteMonitor interface {
	IMonitor
	GetConfig() IMonitorConfig
	SetConfig(IMonitorConfig)
}

type IMonitorConfig interface {
	run() IMonitorResponse
	validate() error
}

func NewConcreteMonitor(base BaseMonitor, config IMonitorConfig) (IConcreteMonitor, error) {
	switch base.Type {
	case Http:
		monitor := &HttpMonitor{
			BaseMonitor: base,
			Config:      *config.(*HttpConfig),
		}
		return monitor, nil
	case Ping:
		monitor := &PingMonitor{
			BaseMonitor: base,
			Config:      *config.(*PingConfig),
		}
		return monitor, nil
	default:
		return nil, fmt.Errorf("unknown monitor type: %s", base.Type)
	}
}

type BaseMonitor struct {
	Id          pgtype.UUID       `json:"id"`
	DisplayId   string            `json:"displayId"`   // Unique identifier for the monitor
	TeamId      pgtype.UUID       `json:"teamId"`      // Id of the owner team of the monitor
	GroupId     pgtype.UUID       `json:"groupId"`     // Id of the owner group of the monitor
	Name        string            `json:"name"`        // Name of the monitor
	Description string            `json:"description"` // Description of the monitor
	Interval    int               `json:"interval"`    // How often to run the monitor in seconds
	Type        MonitorConfigType `json:"type"`        // Type of the monitor (http, ping, etc.)
	util2.Timestamps
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
	if m.GetDisplayId() == "" {
		return fmt.Errorf("monitor DisplayID cannot be empty")
	}
	return nil
}

func (m *BaseMonitor) GenerateDisplayId() {
	m.DisplayId = util.IDFromString(m.GetName())
}

func (m *BaseMonitor) GetDisplayId() string {
	return m.DisplayId
}

func (m *BaseMonitor) GetId() pgtype.UUID {
	return m.Id
}

func (m *BaseMonitor) GetTeamId() pgtype.UUID {
	return m.TeamId
}

func (m *BaseMonitor) SetTeamId(teamId pgtype.UUID) {
	m.TeamId = teamId
}

func (m *BaseMonitor) GetGroupId() pgtype.UUID {
	return m.GroupId
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

func (m *BaseMonitor) SetGroupId(groupId pgtype.UUID) {
	m.GroupId = groupId
}
