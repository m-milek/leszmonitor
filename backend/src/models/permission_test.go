package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPermission(t *testing.T) {
	t.Run("Creates permission with correct fields", func(t *testing.T) {
		perm := newPermission("test:permission", "Test Permission", "Test description")

		assert.Equal(t, "test:permission", perm.ID)
		assert.Equal(t, "Test Permission", perm.Name)
		assert.Equal(t, "Test description", perm.Description)
	})

	t.Run("Creates permission with empty fields", func(t *testing.T) {
		perm := newPermission("", "", "")

		assert.Equal(t, "", perm.ID)
		assert.Equal(t, "", perm.Name)
		assert.Equal(t, "", perm.Description)
	})
}

func TestGetEffectivePermissions(t *testing.T) {
	t.Run("Permission with no implications", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionProjectReader)

		assert.Len(t, perms, 1)
		assert.Contains(t, perms, PermissionProjectReader)
	})

	t.Run("Permission with single level implication", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionProjectEditor)

		assert.Len(t, perms, 2)
		assert.Contains(t, perms, PermissionProjectEditor)
		assert.Contains(t, perms, PermissionProjectReader)
	})

	t.Run("Permission with multi-level implications", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionProjectAdmin)

		assert.Len(t, perms, 3)
		assert.Contains(t, perms, PermissionProjectAdmin)
		assert.Contains(t, perms, PermissionProjectEditor)
		assert.Contains(t, perms, PermissionProjectReader)
	})

	t.Run("Monitor permissions hierarchy", func(t *testing.T) {
		perms := getEffectivePermissions(PermissionMonitorAdmin)

		assert.Len(t, perms, 3)
		assert.Contains(t, perms, PermissionMonitorAdmin)
		assert.Contains(t, perms, PermissionMonitorEditor)
		assert.Contains(t, perms, PermissionMonitorReader)
	})

	t.Run("Permission not in implications map", func(t *testing.T) {
		customPerm := newPermission("custom:perm", "Custom", "Custom permission")
		perms := getEffectivePermissions(customPerm)

		assert.Len(t, perms, 1)
		assert.Contains(t, perms, customPerm)
	})
}

func TestProjectRoleHasPermissions(t *testing.T) {
	t.Run("Owner has all permissions", func(t *testing.T) {
		owner := RoleOwner

		assert.True(t, owner.HasPermissions(PermissionProjectAdmin))
		assert.True(t, owner.HasPermissions(PermissionMonitorAdmin))

		assert.True(t, owner.HasPermissions(PermissionProjectEditor))
		assert.True(t, owner.HasPermissions(PermissionProjectReader))
		assert.True(t, owner.HasPermissions(PermissionMonitorEditor))
		assert.True(t, owner.HasPermissions(PermissionMonitorReader))

		assert.True(t, owner.HasPermissions(PermissionProjectAdmin, PermissionMonitorAdmin))
		assert.True(t, owner.HasPermissions(PermissionProjectReader, PermissionMonitorReader))
	})

	t.Run("Admin has correct permissions", func(t *testing.T) {
		admin := RoleAdmin

		assert.True(t, admin.HasPermissions(PermissionProjectEditor))
		assert.True(t, admin.HasPermissions(PermissionMonitorAdmin))

		assert.True(t, admin.HasPermissions(PermissionProjectReader))
		assert.True(t, admin.HasPermissions(PermissionMonitorEditor))
		assert.True(t, admin.HasPermissions(PermissionMonitorReader))

		assert.False(t, admin.HasPermissions(PermissionProjectAdmin))
	})

	t.Run("Member has limited permissions", func(t *testing.T) {
		member := RoleMember

		assert.True(t, member.HasPermissions(PermissionProjectReader))
		assert.True(t, member.HasPermissions(PermissionMonitorEditor))
		assert.True(t, member.HasPermissions(PermissionMonitorReader))

		assert.False(t, member.HasPermissions(PermissionProjectEditor))
		assert.False(t, member.HasPermissions(PermissionProjectAdmin))
		assert.False(t, member.HasPermissions(PermissionMonitorAdmin))
	})

	t.Run("Viewer has read-only permissions", func(t *testing.T) {
		viewer := RoleViewer

		assert.True(t, viewer.HasPermissions(PermissionProjectReader))
		assert.True(t, viewer.HasPermissions(PermissionMonitorReader))

		assert.False(t, viewer.HasPermissions(PermissionProjectEditor))
		assert.False(t, viewer.HasPermissions(PermissionProjectAdmin))
		assert.False(t, viewer.HasPermissions(PermissionMonitorEditor))
		assert.False(t, viewer.HasPermissions(PermissionMonitorAdmin))
	})

	t.Run("Empty permissions check", func(t *testing.T) {
		owner := RoleOwner
		assert.False(t, owner.HasPermissions())
	})

	t.Run("Nil role", func(t *testing.T) {
		var nilRole *Role
		assert.False(t, nilRole.HasPermissions(PermissionProjectReader))
	})

	t.Run("Invalid role", func(t *testing.T) {
		invalidRole := Role("invalid")
		assert.False(t, invalidRole.HasPermissions(PermissionProjectReader))
	})

	t.Run("Multiple permissions check - all required", func(t *testing.T) {
		member := RoleMember

		assert.True(t, member.HasPermissions(PermissionProjectReader, PermissionMonitorReader))
		assert.False(t, member.HasPermissions(PermissionProjectReader, PermissionProjectAdmin))
		assert.False(t, member.HasPermissions(PermissionMonitorAdmin, PermissionMonitorReader))
	})

	t.Run("Permission not in system", func(t *testing.T) {
		owner := RoleOwner
		unknownPerm := newPermission("unknown:perm", "Unknown", "Unknown permission")

		assert.False(t, owner.HasPermissions(unknownPerm))
	})
}

func TestPermissionImplicationsConsistency(t *testing.T) {
	t.Run("All permissions in implications map exist", func(t *testing.T) {
		allPerms := []Permission{
			PermissionProjectAdmin,
			PermissionProjectEditor,
			PermissionProjectReader,
			PermissionMonitorAdmin,
			PermissionMonitorEditor,
			PermissionMonitorReader,
		}

		for perm, implications := range permissionImplications {
			assert.Contains(t, allPerms, perm, "Permission %s in implications map doesn't exist", perm.ID)

			for _, implied := range implications {
				assert.Contains(t, allPerms, implied, "Implied permission %s doesn't exist", implied.ID)
			}
		}
	})
}
