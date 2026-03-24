package models

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// OrgMember represents a member of a org with a specific role.
type OrgMember struct {
	ID       pgtype.UUID `json:"id"`       // ID is the UUID of the user who this instance represents
	Username string      `json:"username"` // Username is the username of the user (for display purposes)
	Role     Role        `json:"role"`     // Role is the role of the org member (e.g., owner, admin, member)
	util.Timestamps
}

// NewOrgMember creates a new OrgMember instance and validates it.
func NewOrgMember(id pgtype.UUID, role Role) (*OrgMember, error) {
	member := &OrgMember{
		ID:   id,
		Role: role,
	}
	err := member.Validate()
	return member, err
}

// Validate checks if the OrgMember has a valid ID and Role.
func (tm *OrgMember) Validate() error {
	if !tm.ID.Valid {
		return fmt.Errorf("org member ID %s is not valid UUID", tm.ID.String())
	}
	return tm.Role.Validate()
}
