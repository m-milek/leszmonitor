package models

import (
	"github.com/jackc/pgx/v5/pgtype"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMonitorGroup(t *testing.T) {
	someObjectID := pgtype.UUID{}
	tests := []struct {
		name        string
		groupName   string
		description string
		teamId      pgtype.UUID
		wantErr     bool
		errMsg      string
	}{
		{
			name:        "valid monitor group",
			groupName:   "Production Monitors",
			description: "All production environment monitors",
			teamId:      someObjectID,
			wantErr:     false,
		},
		{
			name:        "empty name",
			groupName:   "",
			description: "Description",
			teamId:      someObjectID,
			wantErr:     true,
			errMsg:      "monitor group name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group, err := NewMonitorGroup(tt.groupName, tt.description, &Team{Id: tt.teamId})

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Nil(t, group)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, group)
				assert.NotEmpty(t, group.DisplayId)
				assert.Equal(t, tt.groupName, group.Name)
				assert.Equal(t, tt.description, group.Description)
				assert.Equal(t, tt.teamId, group.TeamId)
			}
		})
	}
}
