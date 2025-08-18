package models

import (
	"fmt"
	"github.com/m-milek/leszmonitor/util"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MonitorGroup struct {
	ObjectId    bson.ObjectID `json:"-" bson:"_id,omitempty"`         // MongoDB ObjectID for internal use
	Id          string        `json:"id" bson:"id"`                   // Unique identifier for the monitor group
	Name        string        `json:"name" bson:"name"`               // Name of the monitor group
	Description string        `json:"description" bson:"description"` // Description of the monitor group
	TeamId      bson.ObjectID `json:"-" bson:"teamId"`                // ID of the team that owns the monitor group
	MonitorIds  []string      `json:"monitorIds" bson:"monitorIds"`   // List of monitor IDs in the group
}

func NewMonitorGroup(name string, description string, team *Team) (*MonitorGroup, error) {
	group := &MonitorGroup{
		Id:          util.IdFromString(name),
		Name:        name,
		Description: description,
		TeamId:      team.ObjectId,
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
	if g.TeamId.IsZero() {
		return fmt.Errorf("team ID cannot be empty")
	}
	return nil
}

func (g *MonitorGroup) GenerateId() {
	g.Id = util.IdFromString(g.Name)
}
