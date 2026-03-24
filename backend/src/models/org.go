package models

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	util2 "github.com/m-milek/leszmonitor/models/util"
)

// Org represents a set of users working together. Orgs can have multiple members with different roles.
// They are used to manage access to monitors. Each org can own multiple projects and monitors in them.
type Org struct {
	ID pgtype.UUID `json:"id"` // ID is the UUID of the org - database primary key
	util2.DisplayIDFromName
	Description string      `json:"description"` // Description of the org
	Members     []OrgMember `json:"members"`     // Members of the org
	util2.Timestamps
}

// NewOrg creates a new Org instance with the given name, description, and owner UUID.
// The owner is added as the first member of the org with the "owner" role.
func NewOrg(name string, description string, ownerID pgtype.UUID) (*Org, error) {
	initialMembers := []OrgMember{
		{
			ID:   ownerID,
			Role: RoleOwner,
		},
	}

	org := &Org{
		Description: description,
		Members:     initialMembers,
	}
	org.DisplayIDFromName.Init(name)

	err := org.Validate()

	return org, err
}

// IsMember checks if a user with the given UUID is a member of the org.
func (t *Org) IsMember(userID pgtype.UUID) bool {
	exists := false
	for _, member := range t.Members {
		if member.ID == userID {
			exists = true
			break
		}
	}
	return exists
}

// GetMember retrieves the OrgMember with the given UUID.
func (t *Org) GetMember(userID pgtype.UUID) *OrgMember {
	for _, member := range t.Members {
		if member.ID == userID {
			return &member
		}
	}
	return nil
}

// ChangeMemberRole changes the role of a org member with the given UUID to the specified role.
func (t *Org) ChangeMemberRole(userID pgtype.UUID, role Role) error {
	if !t.IsMember(userID) {
		return fmt.Errorf("user %s is not a member of the org", userID)
	}

	for i, member := range t.Members {
		if member.ID == userID {
			t.Members[i].Role = role
			break
		}
	}

	return nil
}

// Validate checks if the Org has valid Name, Description, and Members.
// It also validates each member's role.
func (t *Org) Validate() error {
	if t.Name == "" {
		return fmt.Errorf("org name cannot be empty")
	}
	if len(t.Members) == 0 {
		return fmt.Errorf("org must have at least one member")
	}
	for username, role := range t.Members {
		if err := role.Validate(); err != nil {
			return fmt.Errorf("invalid role for user %d: %w", username, err)
		}
	}
	return nil
}
