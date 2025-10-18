package models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// MonitorGroup represents a group of monitors assigned to a team.
// Monitor groups help organize monitors within a team.
// They are just for organizational purposes - access to monitors is controlled at the team level.
type MonitorGroup struct {
	ID pgtype.UUID `json:"id"` // ID is the UUID of the monitor group - database primary key
	util.DisplayIDFromName
	TeamID      pgtype.UUID `json:"-"`           // DisplayID of the team that owns the monitor group
	Description string      `json:"description"` // Description of the monitor group
	util.Timestamps
}

// NewMonitorGroup creates a new MonitorGroup instance and validates it.
func NewMonitorGroup(name string, description string, team *Team) (*MonitorGroup, error) {
	group := &MonitorGroup{
		Description: description,
		TeamID:      team.ID,
	}
	group.DisplayIDFromName.Init(name)
	err := group.Validate()

	if err != nil {
		return nil, err
	}
	return group, nil
}

// Validate checks if the MonitorGroup has valid Name and DisplayID.
func (g *MonitorGroup) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("monitor group name cannot be empty")
	}
	if g.DisplayID == "" {
		return fmt.Errorf("team DisplayID cannot be empty")
	}
	return nil
}
