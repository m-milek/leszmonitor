package models

import (
	"fmt"
)

type TeamRole string

const (
	TeamRoleOwner  TeamRole = "owner"  // Owner, has full permissions to manage the team
	TeamRoleAdmin  TeamRole = "admin"  // Admin, has full permissions to manage monitors and the team
	TeamRoleMember TeamRole = "member" // Member, can manage monitors and view team details
	TeamRoleViewer TeamRole = "viewer" // Viewer, can only view monitor statuses and team details
)

func (r *TeamRole) Validate() error {
	switch *r {
	case TeamRoleOwner, TeamRoleAdmin, TeamRoleMember, TeamRoleViewer:
		return nil
	default:
		return fmt.Errorf("invalid team role: %s", *r)
	}
}

func (r *TeamRole) HasPermissions(permissions ...Permission) bool {
	if r == nil || len(permissions) == 0 {
		return false
	}

	rolePerms, exists := RolePermissions[*r]
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
