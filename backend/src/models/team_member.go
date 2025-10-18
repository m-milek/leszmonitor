package models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// TeamMember represents a member of a team with a specific role.
type TeamMember struct {
	ID   pgtype.UUID `json:"id"`   // ID is the UUID of the user who this instance represents
	Role Role        `json:"role"` // Role is the role of the team member (e.g., owner, admin, member)
	util.Timestamps
}

// NewTeamMember creates a new TeamMember instance and validates it.
func NewTeamMember(id pgtype.UUID, role Role) (*TeamMember, error) {
	member := &TeamMember{
		ID:   id,
		Role: role,
	}
	err := member.Validate()
	return member, err
}

// Validate checks if the TeamMember has a valid ID and Role.
func (tm *TeamMember) Validate() error {
	if !tm.ID.Valid {
		return fmt.Errorf("team member ID %s is not valid UUID", tm.ID.String())
	}
	return tm.Role.Validate()
}
