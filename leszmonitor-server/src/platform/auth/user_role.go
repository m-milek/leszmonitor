package auth

import (
	"fmt"
)

// Role represents a user's role within the instance.
type Role string

const (
	RoleOwner  Role = "owner"  // RoleOwner has full (instance admin) permissions
	RoleAdmin  Role = "admin"  // RoleAdmin can create, edit, and delete monitors
	RoleWriter Role = "member" // RoleWriter can create and edit monitors
	RoleViewer Role = "viewer" // RoleViewer can only view monitors and their statuses
)

var rolePermissions = map[Role][]Permission{
	RoleOwner: {
		PermissionInstanceAdmin,
	},
	RoleAdmin: {
		PermissionAdmin,
	},
	RoleWriter: {
		PermissionWriter,
	},
	RoleViewer: {
		PermissionReader,
	},
}

// Validate checks if the Role is one of the defined roles.
func (r *Role) Validate() error {
	switch *r {
	case RoleOwner, RoleAdmin, RoleWriter, RoleViewer:
		return nil
	default:
		return fmt.Errorf("invalid user role: %s", *r)
	}
}

// HasPermissions checks if the Role includes all the specified permissions.
// It considers permission implications, so higher-level permissions
// automatically include lower-level permissions.
func (r *Role) HasPermissions(permissions ...Permission) bool {
	if r == nil || len(permissions) == 0 {
		return false
	}

	rolePerms, exists := rolePermissions[*r]
	if !exists {
		return false
	}

	// Build a set of all effective permissions for this role
	effectivePerms := make(map[string]bool)
	for _, perm := range rolePerms {
		for _, effectivePerm := range getEffectivePermissions(perm) {
			effectivePerms[effectivePerm.ID] = true
		}
	}

	// Check if all requested permissions are present
	for _, permission := range permissions {
		if !effectivePerms[permission.ID] {
			return false
		}
	}

	return true
}
