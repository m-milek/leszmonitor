package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTeam(t *testing.T) {
	t.Run("Creates team with valid data", func(t *testing.T) {
		name := "Test Team"
		description := "Test Description"
		ownerId := "owner123"

		team, _ := NewTeam(name, description, ownerId)

		assert.NotEmpty(t, team.Id)
		assert.Equal(t, name, team.Name)
		assert.Equal(t, description, team.Description)
		assert.NotNil(t, team.Members)
		assert.Len(t, team.Members, 1)
		assert.Equal(t, TeamRoleOwner, team.Members[ownerId])
		assert.NotEmpty(t, team.CreatedAt)
		assert.NotEmpty(t, team.UpdatedAt)

		// Verify timestamps are valid RFC3339
		_, err := time.Parse(time.RFC3339, team.CreatedAt)
		assert.NoError(t, err)
		_, err = time.Parse(time.RFC3339, team.UpdatedAt)
		assert.NoError(t, err)
	})

	t.Run("Fails to create team with empty name and description", func(t *testing.T) {
		team, err := NewTeam("", "", "owner123")

		assert.Error(t, err)
		assert.Nil(t, team)
	})

	t.Run("Each team gets unique ID", func(t *testing.T) {
		team1, _ := NewTeam("Team 1", "Desc 1", "owner1")
		team2, _ := NewTeam("Team 2", "Desc 2", "owner2")

		assert.NotEqual(t, team1.Id, team2.Id)
	})
}

func TestTeam_IsMember(t *testing.T) {
	t.Run("Owner is a member", func(t *testing.T) {
		ownerId := "owner123"
		team, _ := NewTeam("Test Team", "Description", ownerId)

		assert.True(t, team.IsMember(ownerId))
	})

	t.Run("Non-member returns false", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		assert.False(t, team.IsMember("nonmember"))
		assert.False(t, team.IsMember(""))
	})

	t.Run("Added member returns true", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("member1", TeamRoleMember)
		require.NoError(t, err)

		assert.True(t, team.IsMember("member1"))
	})

	t.Run("Works with nil Members map", func(t *testing.T) {
		team := &Team{
			Id:      "test",
			Name:    "Test",
			Members: nil,
		}

		assert.False(t, team.IsMember("anyone"))
	})
}

func TestTeam_IsAdmin(t *testing.T) {
	t.Run("Owner is admin", func(t *testing.T) {
		ownerId := "owner123"
		team, _ := NewTeam("Test Team", "Description", ownerId)

		assert.True(t, team.IsAdmin(ownerId))
	})

	t.Run("Admin role is admin", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("admin1", TeamRoleAdmin)
		require.NoError(t, err)

		assert.True(t, team.IsAdmin("admin1"))
	})

	t.Run("Member is not admin", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("member1", TeamRoleMember)
		require.NoError(t, err)

		assert.False(t, team.IsAdmin("member1"))
	})

	t.Run("Viewer is not admin", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("viewer1", TeamRoleViewer)
		require.NoError(t, err)

		assert.False(t, team.IsAdmin("viewer1"))
	})

	t.Run("Non-member is not admin", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		assert.False(t, team.IsAdmin("nonmember"))
		assert.False(t, team.IsAdmin(""))
	})
}

func TestTeam_AddMember(t *testing.T) {
	t.Run("Add member successfully", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		err := team.AddMember("member1", TeamRoleMember)

		assert.NoError(t, err)
		assert.Equal(t, TeamRoleMember, team.Members["member1"])
		assert.Len(t, team.Members, 2)
	})

	t.Run("Add multiple members", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		err := team.AddMember("admin1", TeamRoleAdmin)
		assert.NoError(t, err)

		err = team.AddMember("member1", TeamRoleMember)
		assert.NoError(t, err)

		err = team.AddMember("viewer1", TeamRoleViewer)
		assert.NoError(t, err)

		assert.Len(t, team.Members, 4)
		assert.Equal(t, TeamRoleAdmin, team.Members["admin1"])
		assert.Equal(t, TeamRoleMember, team.Members["member1"])
		assert.Equal(t, TeamRoleViewer, team.Members["viewer1"])
	})

	t.Run("Cannot add existing member", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		err := team.AddMember("owner123", TeamRoleMember)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already a member")
		assert.Equal(t, TeamRoleOwner, team.Members["owner123"]) // Role unchanged
	})

	t.Run("Add member to team with nil Members map", func(t *testing.T) {
		team := &Team{
			Id:      "test",
			Name:    "Test",
			Members: nil,
		}

		err := team.AddMember("member1", TeamRoleMember)

		assert.NoError(t, err)
		assert.NotNil(t, team.Members)
		assert.Equal(t, TeamRoleMember, team.Members["member1"])
	})

	t.Run("Add member with empty username", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		err := team.AddMember("", TeamRoleMember)

		assert.NoError(t, err) // Should succeed but might want to validate
		assert.Equal(t, TeamRoleMember, team.Members[""])
	})
}

func TestTeam_RemoveMember(t *testing.T) {
	t.Run("Remove existing member", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("member1", TeamRoleMember)
		require.NoError(t, err)

		team.RemoveMember("member1")

		assert.False(t, team.IsMember("member1"))
		assert.Len(t, team.Members, 1)
	})

	t.Run("Remove non-existent member", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		originalLen := len(team.Members)

		team.RemoveMember("nonexistent")

		assert.Len(t, team.Members, originalLen)
	})

	t.Run("Remove owner", func(t *testing.T) {
		ownerId := "owner123"
		team, _ := NewTeam("Test Team", "Description", ownerId)

		team.RemoveMember(ownerId)

		assert.False(t, team.IsMember(ownerId))
		assert.Len(t, team.Members, 0)
	})

	t.Run("Remove from nil Members map", func(t *testing.T) {
		team := &Team{
			Id:      "test",
			Name:    "Test",
			Members: nil,
		}

		// Should not panic
		team.RemoveMember("anyone")
		assert.Nil(t, team.Members)
	})

	t.Run("Remove all members", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		team.AddMember("member1", TeamRoleMember)
		team.AddMember("member2", TeamRoleMember)

		team.RemoveMember("owner123")
		team.RemoveMember("member1")
		team.RemoveMember("member2")

		assert.Len(t, team.Members, 0)
		assert.NotNil(t, team.Members) // Map still exists, just empty
	})
}

func TestTeam_ChangeMemberRole(t *testing.T) {
	t.Run("Change member role successfully", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("member1", TeamRoleMember)
		require.NoError(t, err)

		err = team.ChangeMemberRole("member1", TeamRoleAdmin)

		assert.NoError(t, err)
		assert.Equal(t, TeamRoleAdmin, team.Members["member1"])
	})

	t.Run("Change owner role", func(t *testing.T) {
		ownerId := "owner123"
		team, _ := NewTeam("Test Team", "Description", ownerId)

		err := team.ChangeMemberRole(ownerId, TeamRoleMember)

		assert.NoError(t, err)
		assert.Equal(t, TeamRoleMember, team.Members[ownerId])
	})

	t.Run("Cannot change non-member role", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		err := team.ChangeMemberRole("nonmember", TeamRoleAdmin)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a member")
		assert.False(t, team.IsMember("nonmember"))
	})

	t.Run("Cannot change to invalid role", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("member1", TeamRoleMember)
		require.NoError(t, err)

		err = team.ChangeMemberRole("member1", TeamRole("invalid"))

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid role")
		assert.Equal(t, TeamRoleMember, team.Members["member1"]) // Role unchanged
	})

	t.Run("Role transitions", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")
		err := team.AddMember("user1", TeamRoleViewer)
		require.NoError(t, err)

		// Viewer -> Member
		err = team.ChangeMemberRole("user1", TeamRoleMember)
		assert.NoError(t, err)
		assert.Equal(t, TeamRoleMember, team.Members["user1"])

		// Member -> Admin
		err = team.ChangeMemberRole("user1", TeamRoleAdmin)
		assert.NoError(t, err)
		assert.Equal(t, TeamRoleAdmin, team.Members["user1"])

		// Admin -> Owner
		err = team.ChangeMemberRole("user1", TeamRoleOwner)
		assert.NoError(t, err)
		assert.Equal(t, TeamRoleOwner, team.Members["user1"])

		// Owner -> Viewer
		err = team.ChangeMemberRole("user1", TeamRoleViewer)
		assert.NoError(t, err)
		assert.Equal(t, TeamRoleViewer, team.Members["user1"])
	})
}

func TestTeam_Timestamps(t *testing.T) {
	t.Run("Timestamps are RFC3339 format", func(t *testing.T) {
		team, _ := NewTeam("Test Team", "Description", "owner123")

		createdTime, err := time.Parse(time.RFC3339, team.CreatedAt)
		assert.NoError(t, err)

		updatedTime, err := time.Parse(time.RFC3339, team.UpdatedAt)
		assert.NoError(t, err)

		assert.WithinDuration(t, time.Now(), createdTime, 5*time.Second)
		assert.WithinDuration(t, time.Now(), updatedTime, 5*time.Second)
	})
}
