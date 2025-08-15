package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPermission(t *testing.T) {
	t.Run("Creates permission with correct fields", func(t *testing.T) {
		perm := NewPermission("test:permission", "Test Permission", "Test description")

		assert.Equal(t, "test:permission", perm.ID)
		assert.Equal(t, "Test Permission", perm.Name)
		assert.Equal(t, "Test description", perm.Description)
	})

	t.Run("Creates permission with empty fields", func(t *testing.T) {
		perm := NewPermission("", "", "")

		assert.Equal(t, "", perm.ID)
		assert.Equal(t, "", perm.Name)
		assert.Equal(t, "", perm.Description)
	})
}

func TestGetEffectivePermissions(t *testing.T) {
	t.Run("Permission with no implications", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionTeamReader)

		assert.Len(t, perms, 1)
		assert.Contains(t, perms, PermissionTeamReader)
	})

	t.Run("Permission with single level implication", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionTeamEditor)

		assert.Len(t, perms, 2)
		assert.Contains(t, perms, PermissionTeamEditor)
		assert.Contains(t, perms, PermissionTeamReader)
	})

	t.Run("Permission with multi-level implications", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionTeamAdmin)

		assert.Len(t, perms, 3)
		assert.Contains(t, perms, PermissionTeamAdmin)
		assert.Contains(t, perms, PermissionTeamEditor)
		assert.Contains(t, perms, PermissionTeamReader)
	})

	t.Run("Monitor permissions hierarchy", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionMonitorAdmin)

		assert.Len(t, perms, 3)
		assert.Contains(t, perms, PermissionMonitorAdmin)
		assert.Contains(t, perms, PermissionMonitorEditor)
		assert.Contains(t, perms, PermissionMonitorReader)
	})

	t.Run("Permission not in implications map", func(t *testing.T) {
		customPerm := NewPermission("custom:perm", "Custom", "Custom permission")
		perms := getEffectivePermissions(customPerm)

		assert.Len(t, perms, 1)
		assert.Contains(t, perms, customPerm)
	})
}

func TestTeamRoleHasPermissions(t *testing.T) {
	t.Run("Owner has all permissions", func(t *testing.T) {
		owner := TeamRoleOwner

		// Direct permissions
		assert.True(t, owner.HasPermissions(PermissionTeamAdmin))
		assert.True(t, owner.HasPermissions(PermissionMonitorAdmin))

		// Implied permissions
		assert.True(t, owner.HasPermissions(PermissionTeamEditor))
		assert.True(t, owner.HasPermissions(PermissionTeamReader))
		assert.True(t, owner.HasPermissions(PermissionMonitorEditor))
		assert.True(t, owner.HasPermissions(PermissionMonitorReader))

		// Multiple permissions
		assert.True(t, owner.HasPermissions(PermissionTeamAdmin, PermissionMonitorAdmin))
		assert.True(t, owner.HasPermissions(PermissionTeamReader, PermissionMonitorReader))
	})

	t.Run("Admin has correct permissions", func(t *testing.T) {
		admin := TeamRoleAdmin

		// Direct permissions
		assert.True(t, admin.HasPermissions(PermissionTeamEditor))
		assert.True(t, admin.HasPermissions(PermissionMonitorAdmin))

		// Implied permissions
		assert.True(t, admin.HasPermissions(PermissionTeamReader))
		assert.True(t, admin.HasPermissions(PermissionMonitorEditor))
		assert.True(t, admin.HasPermissions(PermissionMonitorReader))

		// Should not have
		assert.False(t, admin.HasPermissions(PermissionTeamAdmin))
	})

	t.Run("Member has limited permissions", func(t *testing.T) {
		member := TeamRoleMember

		// Direct permissions
		assert.True(t, member.HasPermissions(PermissionTeamReader))
		assert.True(t, member.HasPermissions(PermissionMonitorEditor))

		// Implied permissions
		assert.True(t, member.HasPermissions(PermissionMonitorReader))

		// Should not have
		assert.False(t, member.HasPermissions(PermissionTeamEditor))
		assert.False(t, member.HasPermissions(PermissionTeamAdmin))
		assert.False(t, member.HasPermissions(PermissionMonitorAdmin))
	})

	t.Run("Viewer has read-only permissions", func(t *testing.T) {
		viewer := TeamRoleViewer

		// Direct permissions
		assert.True(t, viewer.HasPermissions(PermissionTeamReader))
		assert.True(t, viewer.HasPermissions(PermissionMonitorReader))

		// Should not have any write permissions
		assert.False(t, viewer.HasPermissions(PermissionTeamEditor))
		assert.False(t, viewer.HasPermissions(PermissionTeamAdmin))
		assert.False(t, viewer.HasPermissions(PermissionMonitorEditor))
		assert.False(t, viewer.HasPermissions(PermissionMonitorAdmin))
	})

	t.Run("Empty permissions check", func(t *testing.T) {
		owner := TeamRoleOwner
		assert.False(t, owner.HasPermissions())
	})

	t.Run("Nil role", func(t *testing.T) {
		var nilRole *TeamRole
		assert.False(t, nilRole.HasPermissions(PermissionTeamReader))
	})

	t.Run("Invalid role", func(t *testing.T) {
		invalidRole := TeamRole("invalid")
		assert.False(t, invalidRole.HasPermissions(PermissionTeamReader))
	})

	t.Run("Multiple permissions check - all required", func(t *testing.T) {
		member := TeamRoleMember

		// Should have both
		assert.True(t, member.HasPermissions(PermissionTeamReader, PermissionMonitorReader))

		// Should not have one of them
		assert.False(t, member.HasPermissions(PermissionTeamReader, PermissionTeamAdmin))
		assert.False(t, member.HasPermissions(PermissionMonitorAdmin, PermissionMonitorReader))
	})

	t.Run("Permission not in system", func(t *testing.T) {
		owner := TeamRoleOwner
		unknownPerm := NewPermission("unknown:perm", "Unknown", "Unknown permission")

		assert.False(t, owner.HasPermissions(unknownPerm))
	})
}

func TestPermissionImplicationsConsistency(t *testing.T) {
	t.Run("All permissions in implications map exist", func(t *testing.T) {
		allPerms := []Permission{
			PermissionTeamAdmin,
			PermissionTeamEditor,
			PermissionTeamReader,
			PermissionMonitorAdmin,
			PermissionMonitorEditor,
			PermissionMonitorReader,
		}

		for perm, implications := range PermissionImplications {
			// Check that the key permission exists
			assert.Contains(t, allPerms, perm, "Permission %s in implications map doesn't exist", perm.ID)

			// Check that all implied permissions exist
			for _, implied := range implications {
				assert.Contains(t, allPerms, implied, "Implied permission %s doesn't exist", implied.ID)
			}
		}
	})
}
