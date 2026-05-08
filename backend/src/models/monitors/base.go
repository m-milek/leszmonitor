package monitors

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	util2 "github.com/m-milek/leszmonitor/models/util"
	"github.com/m-milek/leszmonitor/util"
)

type Monitor struct {
	ID          uuid.UUID        `json:"id" db:"id"`     // ID is the unique identifier for the monitor, generated as a UUID
	Slug        string           `json:"slug" db:"slug"` // Slug is unique in the project
	ProjectSlug string           `json:"projectSlug"`    // ProjectSlug is used to associate the monitor with a project
	Name        string           `json:"name"`           // Name of the monitor
	Description string           `json:"description"`    // Description of the monitor
	Interval    int              `json:"interval"`       // Interval determines how often to run the monitor in seconds
	Type        shared.ProbeType `json:"type"`           // Type of the monitor (http, tcp, etc.)
	ProbeConfig string           `json:"probeConfig"`    // JSON string containing the specific configuration for the monitor type
	util2.Timestamps
}

type Probe interface {
	Run(ctx context.Context, monitorID uuid.UUID) monitorresult.IMonitorResult
	Validate() error
}

type monitorTypeExtractor struct {
	Type shared.ProbeType `json:"type"`
}

func (m *Monitor) Validate() error {
	if uuid.Validate(m.ID.String()) != nil {
		return fmt.Errorf("monitor ID cannot be null")
	}
	if m.Name == "" {
		return fmt.Errorf("monitor name cannot be empty")
	}
	if m.Interval <= 0 {
		return fmt.Errorf("monitor interval must be greater than zero")
	}
	if m.Type == "" {
		return fmt.Errorf("monitor type cannot be empty")
	}
	if m.Slug == "" {
		return fmt.Errorf("monitor slug cannot be empty")
	}
	return nil
}

func (m *Monitor) GenerateSlug() {
	m.Slug = util.SlugFromString(m.Name)
}
