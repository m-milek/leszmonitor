package auth

// Permission represents a specific action that can be performed within the system.
// Examples include "read:monitor", "edit:monitor", etc.
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

var PermissionInstanceAdmin = newPermission("admin:instance", "Instance Admin", "Full access to all resources and settings in the instance.")
var PermissionAdmin = newPermission("admin:monitor", "Delete Monitors", "Allows deleting monitors.")
var PermissionWriter = newPermission("edit:monitor", "Manage Monitors", "Allows editing and creating monitors.")
var PermissionReader = newPermission(
	"read:monitor",
	"Read Monitors",
	"Allows reading monitor details and statuses.",
)

// permissionImplications defines which permissions imply other permissions.
// For example, having PermissionAdmin implies having the lower PermissionWriter and PermissionReader permissions.
var permissionImplications = map[Permission][]Permission{
	PermissionInstanceAdmin: {PermissionAdmin, PermissionWriter, PermissionReader},
	PermissionAdmin:         {PermissionWriter},
	PermissionWriter:        {PermissionReader},
	PermissionReader:        {}, // No implications
}

// getEffectivePermissions expands a single permission to include all implied permissions.
// For example, if a user has PermissionAdmin, this function returns PermissionAdmin, PermissionWriter, and PermissionReader.
func getEffectivePermissions(perm Permission) []Permission {
	result := []Permission{perm}

	if implied, exists := permissionImplications[perm]; exists {
		for _, impliedPerm := range implied {
			result = append(result, getEffectivePermissions(impliedPerm)...)
		}
	}

	return result
}
