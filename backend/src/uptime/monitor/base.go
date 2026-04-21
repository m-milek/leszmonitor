package monitors

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	util2 "github.com/m-milek/leszmonitor/models/util"
	"github.com/m-milek/leszmonitor/util"
)

type IMonitor interface {
	Run() IMonitorResponse
	Validate() error
	GetID() pgtype.UUID
	GetSlug() string
	GenerateSlug()
	GetProjectID() pgtype.UUID
	SetProjectID(uuid pgtype.UUID)
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
	Slug        string            `json:"slug"`        // Unique identifier for the monitor
	ProjectID   pgtype.UUID       `json:"projectID"`   // ID of the project this monitor belongs to
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
	if m.GetSlug() == "" {
		return fmt.Errorf("monitor slug cannot be empty")
	}
	return nil
}

func (m *BaseMonitor) GenerateSlug() {
	m.Slug = util.IDFromString(m.GetName())
}

func (m *BaseMonitor) GetSlug() string {
	return m.Slug
}

func (m *BaseMonitor) GetID() pgtype.UUID {
	return m.ID
}

func (m *BaseMonitor) GetProjectID() pgtype.UUID {
	return m.ProjectID
}

func (m *BaseMonitor) SetProjectID(projectID pgtype.UUID) {
	m.ProjectID = projectID
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
