package models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/util"
)

type Team struct {
	Id          pgtype.UUID  `json:"id"`          // MongoDB ObjectID for internal use
	DisplayId   string       `json:"displayId"`   // Unique identifier for the team
	Name        string       `json:"name"`        // Name of the team
	Description string       `json:"description"` // Description of the team
	Members     []TeamMember `json:"members"`     // Map of team members with their roles
	Timestamps
}

func NewTeam(name string, description string, ownerId pgtype.UUID) (*Team, error) {
	initialMembers := []TeamMember{
		{
			Id:   ownerId,
			Role: TeamRoleOwner,
		},
	}

	team := &Team{
		DisplayId:   util.IdFromString(name),
		Name:        name,
		Description: description,
		Members:     initialMembers,
	}

	err := team.Validate()
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (t *Team) IsMember(userId pgtype.UUID) bool {
	exists := false
	for _, member := range t.Members {
		if member.Id == userId {
			exists = true
			break
		}
	}
	return exists
}

func (t *Team) GetMember(userId pgtype.UUID) *TeamMember {
	for _, member := range t.Members {
		if member.Id == userId {
			return &member
		}
	}
	return nil
}

func (t *Team) ChangeMemberRole(userId pgtype.UUID, role TeamRole) error {
	if !t.IsMember(userId) {
		return fmt.Errorf("user %s is not a member of the team", userId)
	}

	if role.Validate() != nil {
		return fmt.Errorf("invalid role: %s", role)
	}

	for i, member := range t.Members {
		if member.Id == userId {
			t.Members[i].Role = role
			break
		}
	}

	return nil
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
	t.DisplayId = util.IdFromString(t.Name)
}
