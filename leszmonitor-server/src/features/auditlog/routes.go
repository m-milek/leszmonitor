package auditlog

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

// RegisterRoutes registers the audit log API routes on the protected router.
func RegisterRoutes(
	protectedRouter *http.ServeMux,
	c AuditLogAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	protectedRouter.HandleFunc(
		"GET /api/v1/audit-log",
		requirePermission(auth.PermissionInstanceAdmin)(c.GetAuditLogByQueryHandler),
	)
}
