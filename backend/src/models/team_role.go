package models

import (
	"fmt"
)

// Role represents the role of a team member within a team.
type Role string

const (
	RoleOwner  Role = "owner"  // RoleOwner has full permissions to manage the team
	RoleAdmin  Role = "admin"  // RoleAdmin has full permissions to manage monitors and the team
	RoleMember Role = "member" // RoleMember can manage monitors and view team details
	RoleViewer Role = "viewer" // RoleViewer can only view monitor statuses and team details
)

var rolePermissions = map[Role][]Permission{
	RoleOwner: {
		PermissionTeamAdmin,
		PermissionMonitorAdmin,
	},
	RoleAdmin: {
		PermissionTeamEditor,
		PermissionMonitorAdmin,
	},
	RoleMember: {
		PermissionTeamReader,
		PermissionMonitorEditor,
	},
	RoleViewer: {
		PermissionTeamReader,
		PermissionMonitorReader,
	},
}

// Validate checks if the Role is one of the defined roles.
func (r *Role) Validate() error {
	switch *r {
	case RoleOwner, RoleAdmin, RoleMember, RoleViewer:
		return nil
	default:
		return fmt.Errorf("invalid team role: %s", *r)
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
