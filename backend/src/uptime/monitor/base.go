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
	GetID() pgtype.UUID
	GetDisplayID() string
	GenerateDisplayID()
	GetTeamID() pgtype.UUID
	SetTeamID(uuid pgtype.UUID)
	GetGroupID() pgtype.UUID
	SetGroupID(uuid pgtype.UUID)
	GetName() string
	GetDescription() string
	GetInterval() time.Duration
	GetType() MonitorConfigType
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
	case httpType:
		monitor := &httpMonitor{
			BaseMonitor: base,
			Config:      *config.(*httpConfig),
		}
		return monitor, nil
	case pingType:
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
	ID          pgtype.UUID       `json:"id"`
	DisplayID   string            `json:"displayId"`   // Unique identifier for the monitor
	TeamID      pgtype.UUID       `json:"teamId"`      // ID of the owner team of the monitor
	GroupID     pgtype.UUID       `json:"groupId"`     // ID of the owner group of the monitor
	Name        string            `json:"name"`        // Name of the monitor
	Description string            `json:"description"` // Description of the monitor
	Interval    int               `json:"interval"`    // How often to run the monitor in seconds
	Type        MonitorConfigType `json:"type"`        // Type of the monitor (httpType, pingType, etc.)
	util2.Timestamps
}

type MonitorConfigType string

const (
	httpType MonitorConfigType = "http"
	pingType MonitorConfigType = "ping"
)

type monitorTypeExtractor struct {
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
	if m.GetDisplayID() == "" {
		return fmt.Errorf("monitor DisplayID cannot be empty")
	}
	return nil
}

func (m *BaseMonitor) GenerateDisplayID() {
	m.DisplayID = util.IDFromString(m.GetName())
}

func (m *BaseMonitor) GetDisplayID() string {
	return m.DisplayID
}

func (m *BaseMonitor) GetID() pgtype.UUID {
	return m.ID
}

func (m *BaseMonitor) GetTeamID() pgtype.UUID {
	return m.TeamID
}

func (m *BaseMonitor) SetTeamID(teamID pgtype.UUID) {
	m.TeamID = teamID
}

func (m *BaseMonitor) GetGroupID() pgtype.UUID {
	return m.GroupID
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

func (m *BaseMonitor) SetGroupID(groupId pgtype.UUID) {
	m.GroupID = groupId
}
