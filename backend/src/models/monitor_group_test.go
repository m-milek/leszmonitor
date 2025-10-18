package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMonitorGroup_CreatesGroupAndValidates(t *testing.T) {
	team := &Team{ID: pgtype.UUID{}}

	group, err := NewMonitorGroup("My Group!", "desc", team)

	assert.NoError(t, err)
	if assert.NotNil(t, group) {
		assert.Equal(t, "My Group!", group.Name)
		assert.Equal(t, "my-group", group.DisplayID)
		assert.Equal(t, "desc", group.Description)
		assert.Equal(t, team.ID, group.TeamID)
	}
}

func TestNewMonitorGroup_EmptyName_ReturnsError(t *testing.T) {
	team := &Team{ID: pgtype.UUID{}}

	group, err := NewMonitorGroup("", "anything", team)

	assert.Nil(t, group)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "monitor group name cannot be empty")
	}
}

func TestMonitorGroup_Validate_ErrorWhenDisplayIDEmpty(t *testing.T) {
	g := &MonitorGroup{}
	g.Name = "Has Name"
	g.DisplayID = ""

	err := g.Validate()

	if assert.Error(t, err) {
		assert.Equal(t, "team DisplayID cannot be empty", err.Error())
	}
}

func TestMonitorGroup_Validate_Success(t *testing.T) {
	g := &MonitorGroup{}
	g.Name = "Any"
	g.DisplayID = "any"

	assert.NoError(t, g.Validate())
}
