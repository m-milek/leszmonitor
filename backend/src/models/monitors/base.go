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
	ID                     uuid.UUID        `json:"id" db:"id"`                                           // ID is the unique identifier for the monitor, generated as a UUID
	Slug                   string           `json:"slug" db:"slug"`                                       // Slug is unique in the project
	ProjectID              uuid.UUID        `json:"projectID" db:"project_id"`                            // ProjectID is used to associate the monitor with a project
	Name                   string           `json:"name" db:"name"`                                       // Name of the monitor
	Description            string           `json:"description" db:"description"`                         // Description of the monitor
	Interval               int              `json:"interval" db:"interval"`                               // Interval determines how often to run the monitor in seconds
	Type                   shared.ProbeType `json:"type" db:"kind"`                                       // Type of the monitor (http, tcp, etc.)
	ProbeConfig            string           `json:"probeConfig" db:"config"`                              // JSON string containing the specific configuration for the monitor type
	ResultRetentionSeconds int              `json:"resultRetentionSeconds" db:"result_retention_seconds"` // ResultRetentionSeconds determines how long to keep the monitor results in seconds
	util2.Timestamps
}

type Probe interface {
	Run(ctx context.Context, monitorID uuid.UUID) monitorresult.IMonitorResult
	Validate() error
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
	if m.ResultRetentionSeconds <= 0 {
		return fmt.Errorf("monitor result retention period must be greater than zero")
	}

	return nil
}

func (m *Monitor) GenerateSlug() {
	m.Slug = util.SlugFromString(m.Name)
}
