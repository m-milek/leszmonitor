package models

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/util"
)

// ProjectMember represents a member of a project with a specific role.
type ProjectMember struct {
	ID       uuid.UUID `json:"id" db:"id"`             // ID is the UUID of the user who this instance represents
	Username string    `json:"username" db:"username"` // Username is the username of the user (for display purposes)
	Role     Role      `json:"role" db:"role"`         // Role is the role of the project member (e.g., owner, admin, member)
	util.Timestamps
}

// NewProjectMember creates a new ProjectMember instance and validates it.
func NewProjectMember(id uuid.UUID, role Role) (*ProjectMember, error) {
	member := &ProjectMember{
		ID:   id,
		Role: role,
	}
	err := member.Validate()
	return member, err
}

// Validate checks if the ProjectMember has a valid ID and Role.
func (pm *ProjectMember) Validate() error {
	if pm.ID == uuid.Nil {
		return fmt.Errorf("project member ID %s is not valid UUID", pm.ID.String())
	}
	return pm.Role.Validate()
}
