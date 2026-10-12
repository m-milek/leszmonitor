package auditlog

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

// RegisterRoutes registers the audit log API routes.
func RegisterRoutes(
	router *http.ServeMux,
	c AuditLogAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"GET /api/v1/audit-log",
		requirePermission(auth.PermissionInstanceAdmin)(c.GetAuditLogByQueryHandler),
	)
}
