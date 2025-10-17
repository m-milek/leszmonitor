package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMonitorGroup(t *testing.T) {
	someObjectID := bson.NewObjectID()
	tests := []struct {
		name        string
		groupName   string
		description string
		teamId      bson.ObjectID
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
			group, err := NewMonitorGroup(tt.groupName, tt.description, &Team{
				Id: tt.teamId,
			})

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
				assert.Empty(t, group.MonitorIds)
			}
		})
	}
}

func TestMonitorGroup_AddMonitor(t *testing.T) {
	group := &MonitorGroup{
		DisplayId:  "test-id",
		Name:       "Test Group",
		TeamId:     bson.NewObjectID(),
		MonitorIds: []string{},
	}

	// Add first monitor
	group.AddMonitor("monitor-1")
	assert.Len(t, group.MonitorIds, 1)
	assert.Contains(t, group.MonitorIds, "monitor-1")

	// Add second monitor
	group.AddMonitor("monitor-2")
	assert.Len(t, group.MonitorIds, 2)
	assert.Contains(t, group.MonitorIds, "monitor-2")

	// Add duplicate monitor (shouldn't still add it)
	group.AddMonitor("monitor-1")
	assert.Len(t, group.MonitorIds, 2)
}

func TestMonitorGroup_RemoveMonitor(t *testing.T) {
	tests := []struct {
		name             string
		initialMonitors  []string
		removeId         string
		expectedMonitors []string
	}{
		{
			name:             "remove existing monitor",
			initialMonitors:  []string{"monitor-1", "monitor-2", "monitor-3"},
			removeId:         "monitor-2",
			expectedMonitors: []string{"monitor-1", "monitor-3"},
		},
		{
			name:             "remove non-existing monitor",
			initialMonitors:  []string{"monitor-1", "monitor-2"},
			removeId:         "monitor-99",
			expectedMonitors: []string{"monitor-1", "monitor-2"},
		},
		{
			name:             "remove from empty list",
			initialMonitors:  []string{},
			removeId:         "monitor-1",
			expectedMonitors: []string{},
		},
		{
			name:             "remove first monitor",
			initialMonitors:  []string{"monitor-1", "monitor-2", "monitor-3"},
			removeId:         "monitor-1",
			expectedMonitors: []string{"monitor-2", "monitor-3"},
		},
		{
			name:             "remove last monitor",
			initialMonitors:  []string{"monitor-1", "monitor-2", "monitor-3"},
			removeId:         "monitor-3",
			expectedMonitors: []string{"monitor-1", "monitor-2"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			group := &MonitorGroup{
				MonitorIds: make([]string, len(tt.initialMonitors)),
			}
			copy(group.MonitorIds, tt.initialMonitors)

			group.RemoveMonitor(tt.removeId)
			assert.Equal(t, tt.expectedMonitors, group.MonitorIds)
		})
	}
}

func TestMonitorGroup_Validate(t *testing.T) {
	tests := []struct {
		name    string
		group   *MonitorGroup
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid group",
			group: &MonitorGroup{
				Name:   "Valid Group",
				TeamId: bson.NewObjectID(),
			},
			wantErr: false,
		},
		{
			name: "empty name",
			group: &MonitorGroup{
				Name:   "",
				TeamId: bson.NewObjectID(),
			},
			wantErr: true,
			errMsg:  "monitor group name cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.group.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
