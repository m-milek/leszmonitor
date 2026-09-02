package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPermission(t *testing.T) {
	t.Parallel()
	t.Run("Creates permission with correct fields", func(t *testing.T) {
		t.Parallel()
		perm := newPermission("test:permission", "Test Permission", "Test description")

		assert.Equal(t, "test:permission", perm.ID)
		assert.Equal(t, "Test Permission", perm.Name)
		assert.Equal(t, "Test description", perm.Description)
	})

	t.Run("Creates permission with empty fields", func(t *testing.T) {
		t.Parallel()
		perm := newPermission("", "", "")

		assert.Empty(t, perm.ID)
		assert.Empty(t, perm.Name)
		assert.Empty(t, perm.Description)
	})
}

func TestGetEffectivePermissions(t *testing.T) {
	t.Parallel()
	t.Run("Permission with no implications", func(t *testing.T) {
		t.Parallel()
		perms := getEffectivePermissions(PermissionReader)

		assert.Len(t, perms, 1)
		assert.Contains(t, perms, PermissionReader)
	})

	t.Run("Permission with single level implication", func(t *testing.T) {
		t.Parallel()
		perms := getEffectivePermissions(PermissionWriter)

		assert.Len(t, perms, 2)
		assert.Contains(t, perms, PermissionWriter)
		assert.Contains(t, perms, PermissionReader)
	})

	t.Run("Permission with multi-level implications", func(t *testing.T) {
		t.Parallel()
		perms := getEffectivePermissions(PermissionAdmin)

		assert.Len(t, perms, 3)
		assert.Contains(t, perms, PermissionAdmin)
		assert.Contains(t, perms, PermissionWriter)
		assert.Contains(t, perms, PermissionReader)
	})

	t.Run("Instance admin permissions hierarchy", func(t *testing.T) {
		t.Parallel()
		perms := getEffectivePermissions(PermissionInstanceAdmin)

		assert.Contains(t, perms, PermissionInstanceAdmin)
		assert.Contains(t, perms, PermissionAdmin)
		assert.Contains(t, perms, PermissionWriter)
		assert.Contains(t, perms, PermissionReader)
	})

	t.Run("Permission not in implications map", func(t *testing.T) {
		t.Parallel()
		customPerm := newPermission("custom:perm", "Custom", "Custom permission")
		perms := getEffectivePermissions(customPerm)

		assert.Len(t, perms, 1)
		assert.Contains(t, perms, customPerm)
	})
}

func TestRoleHasPermissions(t *testing.T) {
	t.Parallel()
	t.Run("Owner has all permissions", func(t *testing.T) {
		t.Parallel()
		owner := RoleOwner

		assert.True(t, owner.HasPermissions(PermissionInstanceAdmin))
		assert.True(t, owner.HasPermissions(PermissionAdmin))
		assert.True(t, owner.HasPermissions(PermissionWriter))
		assert.True(t, owner.HasPermissions(PermissionReader))

		assert.True(t, owner.HasPermissions(PermissionInstanceAdmin, PermissionAdmin))
		assert.True(t, owner.HasPermissions(PermissionWriter, PermissionReader))
	})

	t.Run("Admin has correct permissions", func(t *testing.T) {
		t.Parallel()
		admin := RoleAdmin

		assert.True(t, admin.HasPermissions(PermissionAdmin))
		assert.True(t, admin.HasPermissions(PermissionWriter))
		assert.True(t, admin.HasPermissions(PermissionReader))

		assert.False(t, admin.HasPermissions(PermissionInstanceAdmin))
	})

	t.Run("Member has limited permissions", func(t *testing.T) {
		t.Parallel()
		member := RoleWriter

		assert.True(t, member.HasPermissions(PermissionWriter))
		assert.True(t, member.HasPermissions(PermissionReader))

		assert.False(t, member.HasPermissions(PermissionAdmin))
		assert.False(t, member.HasPermissions(PermissionInstanceAdmin))
	})

	t.Run("Viewer has read-only permissions", func(t *testing.T) {
		t.Parallel()
		viewer := RoleViewer

		assert.True(t, viewer.HasPermissions(PermissionReader))

		assert.False(t, viewer.HasPermissions(PermissionWriter))
		assert.False(t, viewer.HasPermissions(PermissionAdmin))
		assert.False(t, viewer.HasPermissions(PermissionInstanceAdmin))
	})

	t.Run("Empty permissions check", func(t *testing.T) {
		t.Parallel()
		owner := RoleOwner
		assert.False(t, owner.HasPermissions())
	})

	t.Run("Nil role", func(t *testing.T) {
		t.Parallel()
		var nilRole *Role
		assert.False(t, nilRole.HasPermissions(PermissionReader))
	})

	t.Run("Invalid role", func(t *testing.T) {
		t.Parallel()
		invalidRole := Role("invalid")
		assert.False(t, invalidRole.HasPermissions(PermissionReader))
	})

	t.Run("Multiple permissions check - all required", func(t *testing.T) {
		t.Parallel()
		member := RoleWriter

		assert.True(t, member.HasPermissions(PermissionWriter, PermissionReader))
		assert.False(t, member.HasPermissions(PermissionReader, PermissionAdmin))
		assert.False(t, member.HasPermissions(PermissionAdmin, PermissionReader))
	})

	t.Run("Permission not in system", func(t *testing.T) {
		t.Parallel()
		owner := RoleOwner
		unknownPerm := newPermission("unknown:perm", "Unknown", "Unknown permission")

		assert.False(t, owner.HasPermissions(unknownPerm))
	})
}

func TestPermissionImplicationsConsistency(t *testing.T) {
	t.Parallel()
	t.Run("All permissions in implications map exist", func(t *testing.T) {
		t.Parallel()
		allPerms := []Permission{
			PermissionInstanceAdmin,
			PermissionAdmin,
			PermissionWriter,
			PermissionReader,
		}

		for perm, implications := range permissionImplications {
			assert.Contains(t, allPerms, perm, "Permission %s in implications map doesn't exist", perm.ID)

			for _, implied := range implications {
				assert.Contains(t, allPerms, implied, "Implied permission %s doesn't exist", implied.ID)
			}
		}
	})
}
