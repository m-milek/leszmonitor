package monitors

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	util2 "github.com/m-milek/leszmonitor/models/util"
	"github.com/m-milek/leszmonitor/util"
)

type IMonitor interface {
	Run() monitorresult.IMonitorResult
	Validate() error
	GetID() uuid.UUID
	GetSlug() string
	GenerateSlug()
	GetProjectSlug() string
	SetProjectSlug(slug string)
	GetName() string
	GetDescription() string
	GetInterval() time.Duration
	GetType() shared.MonitorConfigType
}

type IConcreteMonitor interface {
	IMonitor
	GetConfig() IMonitorConfig
	SetConfig(IMonitorConfig)
}

type IMonitorConfig interface {
	run(id uuid.UUID, monitorType shared.MonitorConfigType) monitorresult.IMonitorResult
	validate() error
}

func NewConcreteMonitor(base BaseMonitor, config IMonitorConfig) (IConcreteMonitor, error) {
	switch base.Type {
	case shared.HttpConfigType:
		monitor := &HttpMonitor{
			BaseMonitor: base,
			Config:      *config.(*HttpConfig),
		}
		return monitor, nil
	case shared.PingConfigType:
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
	ID          uuid.UUID                `json:"id" db:"id"`
	Slug        string                   `json:"slug" db:"slug"` // Unique identifier for the monitor
	ProjectSlug string                   `json:"projectSlug"`    // Slug of the project this monitor belongs to
	Name        string                   `json:"name"`           // Name of the monitor
	Description string                   `json:"description"`    // Description of the monitor
	Interval    int                      `json:"interval"`       // How often to run the monitor in seconds
	Type        shared.MonitorConfigType `json:"type"`           // Type of the monitor (httpType, pingType, etc.)
	util2.Timestamps
}

type monitorTypeExtractor struct {
	Type shared.MonitorConfigType `json:"type"`
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
	m.Slug = util.SlugFromString(m.GetName())
}

func (m *BaseMonitor) GetSlug() string {
	return m.Slug
}

func (m *BaseMonitor) GetID() uuid.UUID {
	return m.ID
}

func (m *BaseMonitor) GetProjectSlug() string {
	return m.ProjectSlug
}

func (m *BaseMonitor) SetProjectSlug(slug string) {
	m.ProjectSlug = slug
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

func (m *BaseMonitor) GetType() shared.MonitorConfigType {
	return m.Type
}
