package models

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// ProjectMember represents a member of a project with a specific role.
type ProjectMember struct {
	ID       pgtype.UUID `json:"id"`       // ID is the UUID of the user who this instance represents
	Username string      `json:"username"` // Username is the username of the user (for display purposes)
	Role     Role        `json:"role"`     // Role is the role of the project member (e.g., owner, admin, member)
	util.Timestamps
}

// NewProjectMember creates a new ProjectMember instance and validates it.
func NewProjectMember(id pgtype.UUID, role Role) (*ProjectMember, error) {
	member := &ProjectMember{
		ID:   id,
		Role: role,
	}
	err := member.Validate()
	return member, err
}

// Validate checks if the ProjectMember has a valid ID and Role.
func (pm *ProjectMember) Validate() error {
	if !pm.ID.Valid {
		return fmt.Errorf("project member ID %s is not valid UUID", pm.ID.String())
	}
	return pm.Role.Validate()
}
