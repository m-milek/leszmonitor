package models

import (
	"github.com/m-milek/leszmonitor/util"
)

type MonitorGroup struct {
	Id          string   `json:"id" bson:"_id"`                  // Unique identifier for the monitor group
	Name        string   `json:"name" bson:"name"`               // Name of the monitor group
	Description string   `json:"description" bson:"description"` // Description of the monitor group
	TeamId      string   `json:"teamId" bson:"teamId"`           // ID of the team that owns the monitor group
	MonitorIds  []string `json:"monitorIds" bson:"monitorIds"`   // List of monitor IDs in the group
}

func NewMonitorGroup(name string, description string, teamId string) *MonitorGroup {
	return &MonitorGroup{
		Id:          util.IdFromString(name),
		Name:        name,
		Description: description,
		TeamId:      teamId,
		MonitorIds:  make([]string, 0),
	}
}

func (g *MonitorGroup) AddMonitor(monitorId string) {
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
