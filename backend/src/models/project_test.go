package models

import (
	"testing"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
)

func TestNewProject_CreatesProjectAndValidates(t *testing.T) {
	org := &Org{ID: pgtype.UUID{}}

	project, err := NewProject("My Project!", "desc", org)

	assert.NoError(t, err)
	if assert.NotNil(t, project) {
		assert.Equal(t, "My Project!", project.Name)
		assert.Equal(t, "my-project", project.DisplayID)
		assert.Equal(t, "desc", project.Description)
		assert.Equal(t, org.ID, project.OrgID)
	}
}

func TestNewProject_EmptyName_ReturnsError(t *testing.T) {
	org := &Org{ID: pgtype.UUID{}}

	project, err := NewProject("", "anything", org)

	assert.Nil(t, project)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "project name cannot be empty")
	}
}

func TestProject_Validate_ErrorWhenDisplayIDEmpty(t *testing.T) {
	g := &Project{}
	g.Name = "Has Name"
	g.DisplayID = ""

	err := g.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "org DisplayID cannot be empty", err.Error())
	}
}

func TestProject_Validate_Success(t *testing.T) {
	g := &Project{}
	g.Name = "Any"
	g.DisplayID = "any"

	assert.NoError(t, g.Validate())
}
