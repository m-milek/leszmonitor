package models

type Permission struct {
	ID          string `json:"id"`          // eg. "read:monitor"
	Name        string `json:"name"`        // Name of the permission
	Description string `json:"description"` // Description of the permission
}

func NewPermission(id, name, description string) Permission {
	return Permission{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

var PermissionMonitorAdmin = NewPermission("admin:monitor", "Delete Monitors", "Allows deleting monitors.")
var PermissionMonitorEditor = NewPermission("edit:monitor", "Manage Monitors", "Allows editing and creating monitors.")
var PermissionMontiorReader = NewPermission("read:monitor", "Read Monitors", "Allows reading monitor details and statuses.")

var PermissionTeamAdmin = NewPermission("admin:team", "Administrate Team", "Allows deleting teams.")
var PermissionTeamEditor = NewPermission("edit:team", "Manage Team", "Allows adding and removing members from the team.")
var PermissionTeamReader = NewPermission("read:team", "Read Teams", "Allows reading team details and members.")

var PermissionImplications = map[Permission][]Permission{
	PermissionTeamAdmin:     {PermissionTeamEditor},
	PermissionTeamEditor:    {PermissionTeamReader},
	PermissionTeamReader:    {}, // No implications
	PermissionMonitorAdmin:  {PermissionMonitorEditor},
	PermissionMonitorEditor: {PermissionMontiorReader},
	PermissionMontiorReader: {}, // No implications
}

var RolePermissions = map[TeamRole][]Permission{
	TeamRoleOwner: {
		PermissionTeamAdmin,
		PermissionMonitorAdmin,
	},
	TeamRoleAdmin: {
		PermissionTeamEditor,
		PermissionMonitorAdmin,
	},
	TeamRoleMember: {
		PermissionTeamReader,
		PermissionMonitorEditor,
	},
	TeamRoleViewer: {
		PermissionTeamReader,
		PermissionMontiorReader,
	},
}

// getEffectivePermissions expands a single permission to include all implied permissions
func getEffectivePermissions(perm Permission) []Permission {
	result := []Permission{perm}

	if implied, exists := PermissionImplications[perm]; exists {
		for _, impliedPerm := range implied {
			result = append(result, getEffectivePermissions(impliedPerm)...)
		}
	}

	return result
}
