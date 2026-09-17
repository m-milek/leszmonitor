package monitors

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/util"
)

type Monitor struct {
	util.Timestamps
	ID                     uuid.UUID       `json:"id"                     db:"id"`                       // ID is the unique identifier for the monitor, generated as a UUID
	Slug                   string          `json:"slug"                   db:"slug"`                     // Slug is unique
	Name                   string          `json:"name"                   db:"name"`                     // Name of the monitor
	Description            string          `json:"description"            db:"description"`              // Description of the monitor
	Interval               int             `json:"interval"               db:"interval"`                 // Interval determines how often to run the monitor in seconds
	Type                   ProbeType       `json:"type"                   db:"kind"`                     // Type of the monitor (http, tcp, etc.)
	ProbeConfig            string          `json:"probeConfig"            db:"config"`                   // JSON string containing the specific configuration for the monitor type
	ResultRetentionSeconds int             `json:"resultRetentionSeconds" db:"result_retention_seconds"` // ResultRetentionSeconds determines how long to keep the monitor results in seconds
	RunState               MonitorRunState `json:"runState" db:"run_state"`                              // RunState indicates whether the monitor is currently running or stopped
	OwnerID                uuid.UUID       `json:"ownerId"                db:"owner_id"`                 // OwnerID is the user who created the monitor. Informational metadata only; not used for visibility filtering.
	TagIDs                 []uuid.UUID     `json:"tagIds"                 db:"-"`
}

type Probe interface {
	Run(ctx context.Context, monitorID uuid.UUID) (IMonitorResult, error)
	Validate() error
}

type MonitorRunState string

const (
	MonitorStateActive  MonitorRunState = "active"
	MonitorStateStopped MonitorRunState = "paused"
)

func IsValidMonitorState(state string) bool {
	return state == string(MonitorStateActive) || state == string(MonitorStateStopped)
}

func InitializeFromPayload(payload Monitor, ownerID uuid.UUID) *Monitor {
	return &Monitor{
		ID:                     uuid.New(),
		OwnerID:                ownerID,
		Slug:                   payload.Slug,
		Name:                   payload.Name,
		Description:            payload.Description,
		Interval:               payload.Interval,
		Type:                   payload.Type,
		ProbeConfig:            payload.ProbeConfig,
		ResultRetentionSeconds: int((48 * time.Hour).Seconds()), // TODO: Make this configurable later
		RunState:               MonitorStateActive,
		TagIDs:                 payload.TagIDs,
	}
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
	if !IsValidMonitorState(string(m.RunState)) {
		return fmt.Errorf("monitor state must be either 'running' or 'stopped'")
	}

	return nil
}

func (m *Monitor) GenerateSlug() {
	m.Slug = util.SlugFromString(m.Name)
}
