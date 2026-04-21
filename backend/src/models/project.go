package models

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// Project is the top-level organizational unit. Projects have multiple members with
// different roles and own all monitors directly.
type Project struct {
	ID pgtype.UUID `json:"id"` // ID is the UUID of the project - database primary key
	util.DisplayIDFromName
	Description string          `json:"description"` // Description of the project
	Members     []ProjectMember `json:"members"`     // Members of the project
	util.Timestamps
}

// NewProject creates a new Project instance with the given name, description, and owner UUID.
// The owner is added as the first member of the project with the "owner" role.
func NewProject(name string, description string, ownerID pgtype.UUID) (*Project, error) {
	initialMembers := []ProjectMember{
		{
			ID:   ownerID,
			Role: RoleOwner,
		},
	}

	project := &Project{
		Description: description,
		Members:     initialMembers,
	}
	project.DisplayIDFromName.Init(name)

	err := project.Validate()
	if err != nil {
		return nil, err
	}
	return project, nil
}

// IsMember checks if a user with the given UUID is a member of the project.
func (p *Project) IsMember(userID pgtype.UUID) bool {
	for _, member := range p.Members {
		if member.ID == userID {
			return true
		}
	}
	return false
}

// GetMember retrieves the ProjectMember with the given UUID.
func (p *Project) GetMember(userID pgtype.UUID) *ProjectMember {
	for i := range p.Members {
		if p.Members[i].ID == userID {
			return &p.Members[i]
		}
	}
	return nil
}

// ChangeMemberRole changes the role of a project member with the given UUID to the specified role.
func (p *Project) ChangeMemberRole(userID pgtype.UUID, role Role) error {
	if !p.IsMember(userID) {
		return fmt.Errorf("user %s is not a member of the project", userID)
	}

	for i, member := range p.Members {
		if member.ID == userID {
			p.Members[i].Role = role
			break
		}
	}

	return nil
}

// Validate checks if the Project has valid Name, Description, and Members.
// It also validates each member's role.
func (p *Project) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if p.DisplayID == "" {
		return fmt.Errorf("project DisplayID cannot be empty")
	}
	if len(p.Members) == 0 {
		return fmt.Errorf("project must have at least one member")
	}
	for i, member := range p.Members {
		if err := member.Validate(); err != nil {
			return fmt.Errorf("invalid member %d: %w", i, err)
		}
		if err := member.Role.Validate(); err != nil {
			return fmt.Errorf("invalid role for member %d: %w", i, err)
		}
	}
	return nil
}
