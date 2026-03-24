package models

import (
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// Project represents a group of monitors assigned to a org.
// Projects help organize monitors within an org.
// They are just for organizational purposes - access to monitors is controlled at the org level.
type Project struct {
	ID pgtype.UUID `json:"id"` // ID is the UUID of the org - database primary key
	util.DisplayIDFromName
	OrgID       pgtype.UUID `json:"-"`           // DisplayID of the org that owns the project
	Description string      `json:"description"` // Description of the org
	util.Timestamps
}

// NewProject creates a new Project instance and validates it.
func NewProject(name string, description string, org *Org) (*Project, error) {
	project := &Project{
		Description: description,
		OrgID:       org.ID,
	}
	project.DisplayIDFromName.Init(name)
	err := project.Validate()

	if err != nil {
		return nil, err
	}
	return project, nil
}

// Validate checks if the Project has valid Name and DisplayID.
func (g *Project) Validate() error {
	if g.Name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if g.DisplayID == "" {
		return fmt.Errorf("org DisplayID cannot be empty")
	}
	return nil
}
