package models

import (
	"fmt"
	"github.com/teris-io/shortid"
	"time"
)

type Team struct {
	Id          string              `json:"id" bson:"id"`                   // Unique identifier for the team
	Name        string              `json:"name" bson:"name"`               // Name of the team
	Description string              `json:"description" bson:"description"` // Description of the team
	Members     map[string]TeamRole `json:"members" bson:"members"`         // Map of team members with their roles
	CreatedAt   string              `json:"createdAt" bson:"createdAt"`     // Creation timestamp of the team
	UpdatedAt   string              `json:"updatedAt" bson:"updatedAt"`     // Last update timestamp of the team
}

func NewTeam(name string, description string, ownerId string) *Team {
	initialMembers := map[string]TeamRole{
		ownerId: TeamRoleOwner, // The owner is automatically added
	}

	return &Team{
		Id:          shortid.MustGenerate(),
		Name:        name,
		Description: description,
		Members:     initialMembers,
		CreatedAt:   time.Now().Format(time.RFC3339),
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
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
	t.UpdatedAt = time.Now().Format(time.RFC3339)

	return nil
}

func (t *Team) RemoveMember(username string) {
	if t.Members == nil {
		return
	}
	delete(t.Members, username)
	t.UpdatedAt = time.Now().Format(time.RFC3339)
}

func (t *Team) ChangeMemberRole(username string, role TeamRole) error {
	if !t.IsMember(username) {
		return fmt.Errorf("user %s is not a member of the team", username)
	}

	if role.Validate() != nil {
		return fmt.Errorf("invalid role: %s", role)
	}

	t.Members[username] = role
	t.UpdatedAt = time.Now().Format(time.RFC3339)
	return nil
}
