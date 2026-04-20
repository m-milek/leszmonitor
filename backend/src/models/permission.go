package models

// Permission represents a specific action that can be performed within the system.
// Examples include "read:monitor", "edit:project", etc.
type Permission struct {
	ID          string `json:"id"`          // ID of the permission - a unique identifier, e.g., "read:monitor"
	Name        string `json:"name"`        // Name of the permission. Used for display purposes.
	Description string `json:"description"` // Description of the permission. Used for display purposes.
}

// newPermission creates a new Permission instance.
func newPermission(id, name, description string) Permission {
	return Permission{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

var PermissionMonitorAdmin = newPermission("admin:monitor", "Delete Monitors", "Allows deleting monitors.")
var PermissionMonitorEditor = newPermission("edit:monitor", "Manage Monitors", "Allows editing and creating monitors.")
var PermissionMonitorReader = newPermission("read:monitor", "Read Monitors", "Allows reading monitor details and statuses.")

var PermissionProjectAdmin = newPermission("admin:project", "Administrate Project", "Allows deleting projects.")
var PermissionProjectEditor = newPermission("edit:project", "Manage Project", "Allows adding and removing members from the project.")
var PermissionProjectReader = newPermission("read:project", "Read Projects", "Allows reading project details and members.")

// permissionImplications defines which permissions imply other permissions.
// For example, having ProjectAdmin permission implies having lower ProjectEditor and ProjectReader permissions.
var permissionImplications = map[Permission][]Permission{
	PermissionProjectAdmin:  {PermissionProjectEditor},
	PermissionProjectEditor: {PermissionProjectReader},
	PermissionProjectReader: {}, // No implications
	PermissionMonitorAdmin:  {PermissionMonitorEditor},
	PermissionMonitorEditor: {PermissionMonitorReader},
	PermissionMonitorReader: {}, // No implications
}

// getEffectivePermissions expands a single permission to include all implied permissions.
// For example, if a user has ProjectAdmin permission, this function will return ProjectAdmin, ProjectEditor, and ProjectReader permissions.
func getEffectivePermissions(perm Permission) []Permission {
	result := []Permission{perm}

	if implied, exists := permissionImplications[perm]; exists {
		for _, impliedPerm := range implied {
			result = append(result, getEffectivePermissions(impliedPerm)...)
		}
	}

	return result
}
