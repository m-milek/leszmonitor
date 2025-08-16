package models

import (
	"fmt"
	"github.com/m-milek/leszmonitor/util"
)

type MonitorGroup struct {
	Id          string   `json:"id" bson:"_id"`                  // Unique identifier for the monitor group
	Name        string   `json:"name" bson:"name"`               // Name of the monitor group
	Description string   `json:"description" bson:"description"` // Description of the monitor group
	TeamId      string   `json:"teamId" bson:"teamId"`           // ID of the team that owns the monitor group
	MonitorIds  []string `json:"monitorIds" bson:"monitorIds"`   // List of monitor IDs in the group
}

func NewMonitorGroup(name string, description string, teamId string) (*MonitorGroup, error) {
	group := &MonitorGroup{
		Id:          util.IdFromString(name),
		Name:        name,
		Description: description,
		TeamId:      teamId,
		MonitorIds:  make([]string, 0),
	}
	err := group.Validate()

	if err != nil {
		return nil, err
	}
	return group, nil
}

func (g *MonitorGroup) AddMonitor(monitorId string) {
	if util.SliceContains(g.MonitorIds, monitorId) {
		return
	}
	g.MonitorIds = append(g.MonitorIds, monitorId)
}

func (g *MonitorGroup) RemoveMonitor(monitorId string) {
	for i, id := range g.MonitorIds {
		if id == monitorId {
			g.MonitorIds = append(g.MonitorIds[:i], g.MonitorIds[i+1:]...)
			break
		}
	}
}

func (g *MonitorGroup) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("monitor group name cannot be empty")
	}
	if g.TeamId == "" {
		return fmt.Errorf("team ID cannot be empty")
	}
	return nil
}
