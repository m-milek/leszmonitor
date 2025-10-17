package models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/util"
)

type MonitorGroup struct {
	Id          pgtype.UUID `json:"id"`          // MongoDB ObjectID for internal use
	DisplayId   string      `json:"displayId"`   // Unique identifier for the monitor group
	Name        string      `json:"name"`        // Name of the monitor group
	Description string      `json:"description"` // Description of the monitor group
	TeamId      pgtype.UUID `json:"-"`           // DisplayId of the team that owns the monitor group
}

func NewMonitorGroup(name string, description string, team *Team) (*MonitorGroup, error) {
	group := &MonitorGroup{
		DisplayId:   util.IdFromString(name),
		Name:        name,
		Description: description,
		TeamId:      team.Id,
	}
	err := group.Validate()

	if err != nil {
		return nil, err
	}
	return group, nil
}

func (g *MonitorGroup) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("monitor group name cannot be empty")
	}
	if g.DisplayId == "" {
		return fmt.Errorf("team DisplayId cannot be empty")
	}
	return nil
}

func (g *MonitorGroup) GenerateId() {
	g.DisplayId = util.IdFromString(g.Name)
}
