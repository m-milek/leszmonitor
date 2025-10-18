package models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	util2 "github.com/m-milek/leszmonitor/models/util"
)

// Team represents a group of users working together. Teams can have multiple members with different roles.
// They are used to manage access to monitors. Each team can own multiple monitor groups and monitors.
type Team struct {
	ID pgtype.UUID `json:"id"` // ID is the UUID of the team - database primary key
	util2.DisplayIDFromName
	Description string       `json:"description"` // Description of the team
	Members     []TeamMember `json:"members"`     // Members of the team
	util2.Timestamps
}

// NewTeam creates a new Team instance with the given name, description, and owner UUID.
// The owner is added as the first member of the team with the "owner" role.
func NewTeam(name string, description string, ownerID pgtype.UUID) (*Team, error) {
	initialMembers := []TeamMember{
		{
			ID:   ownerID,
			Role: RoleOwner,
		},
	}

	team := &Team{
		Description: description,
		Members:     initialMembers,
	}
	team.DisplayIDFromName.Init(name)

	err := team.Validate()

	return team, err
}

// IsMember checks if a user with the given UUID is a member of the team.
func (t *Team) IsMember(userID pgtype.UUID) bool {
	exists := false
	for _, member := range t.Members {
		if member.ID == userID {
			exists = true
			break
		}
	}
	return exists
}

// GetMember retrieves the TeamMember with the given UUID.
func (t *Team) GetMember(userID pgtype.UUID) *TeamMember {
	for _, member := range t.Members {
		if member.ID == userID {
			return &member
		}
	}
	return nil
}

// ChangeMemberRole changes the role of a team member with the given UUID to the specified role.
func (t *Team) ChangeMemberRole(userID pgtype.UUID, role Role) error {
	if !t.IsMember(userID) {
		return fmt.Errorf("user %s is not a member of the team", userID)
	}

	for i, member := range t.Members {
		if member.ID == userID {
			t.Members[i].Role = role
			break
		}
	}

	return nil
}

// Validate checks if the Team has valid Name, Description, and Members.
// It also validates each member's role.
func (t *Team) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("team name cannot be empty")
	}
	if len(t.Members) == 0 {
		return fmt.Errorf("team must have at least one member")
	}
	for username, role := range t.Members {
		if err := role.Validate(); err != nil {
			return fmt.Errorf("invalid role for user %d: %w", username, err)
		}
	}
	return nil
}
