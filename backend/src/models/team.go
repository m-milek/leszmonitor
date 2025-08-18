package models

import (
	"fmt"
	"github.com/m-milek/leszmonitor/util"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type Team struct {
	ObjectId    bson.ObjectID       `json:"-" bson:"_id,omitempty"`         // MongoDB ObjectID for internal use
	Id          string              `json:"id" bson:"id"`                   // Unique identifier for the team
	Name        string              `json:"name" bson:"name"`               // Name of the team
	Description string              `json:"description" bson:"description"` // Description of the team
	Members     map[string]TeamRole `json:"members" bson:"members"`         // Map of team members with their roles
	CreatedAt   string              `json:"createdAt" bson:"createdAt"`     // Creation timestamp of the team
	UpdatedAt   string              `json:"updatedAt" bson:"updatedAt"`     // Last update timestamp of the team
}

func NewTeam(name string, description string, ownerId string) (*Team, error) {
	initialMembers := map[string]TeamRole{
		ownerId: TeamRoleOwner, // The owner is automatically added
	}

	team := &Team{
		Id:          util.IdFromString(name),
		Name:        name,
		Description: description,
		Members:     initialMembers,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}

	err := team.Validate()
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (t *Team) IsMember(username string) bool {
	_, exists := t.Members[username]
	return exists
}

func (t *Team) IsAdmin(username string) bool {
	role, exists := t.Members[username]
	return exists && (role == TeamRoleOwner || role == TeamRoleAdmin)
}

func (t *Team) AddMember(username string, role TeamRole) error {
	if t.IsMember(username) {
		return fmt.Errorf("user %s is already a member of the team", username)
	}

	if t.Members == nil {
		t.Members = make(map[string]TeamRole)
	}
	t.Members[username] = role
	t.updateTimestamps()

	return nil
}

func (t *Team) RemoveMember(username string) {
	if t.Members == nil {
		return
	}
	delete(t.Members, username)
	t.updateTimestamps()
}

func (t *Team) ChangeMemberRole(username string, role TeamRole) error {
	if !t.IsMember(username) {
		return fmt.Errorf("user %s is not a member of the team", username)
	}

	if role.Validate() != nil {
		return fmt.Errorf("invalid role: %s", role)
	}

	t.Members[username] = role
	t.updateTimestamps()
	return nil
}

func (t *Team) updateTimestamps() {
	now := time.Now().Format(time.RFC3339)
	t.UpdatedAt = now
	if t.CreatedAt == "" {
		t.CreatedAt = now
	}
}

func (t *Team) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("team name cannot be empty")
	}
	if t.Description == "" {
		t.Description = "No description provided"
	}
	if len(t.Members) == 0 {
		return fmt.Errorf("team must have at least one member")
	}
	for username, role := range t.Members {
		if err := role.Validate(); err != nil {
			return fmt.Errorf("invalid role for user %s: %w", username, err)
		}
	}
	return nil
}

func (t *Team) GenerateId() {
	t.Id = util.IdFromString(t.Name)
}
