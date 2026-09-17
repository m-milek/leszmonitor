package api

import (
	"embed"
	"net/http"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/users"

	"github.com/m-milek/leszmonitor/api/controllers"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/auth"
)

func SetupRouters(
	publicRouter *http.ServeMux,
	protectedRouter *http.ServeMux,
	staticFiles embed.FS,
	h Handlers,
) {
	// Users
	users.RegisterRoutes(publicRouter, protectedRouter, h.User)

	// Monitors
	monitors.RegisterRoutes(
		protectedRouter,
		h.Monitor,
		h.MonitorResults,
		h.MonitorStats,
		requirePermission(h),
	)

	// Tags
	tags.RegisterRoutes(protectedRouter, h.Tag, requirePermission(h))

	protectedRouter.HandleFunc(
		"GET /api/v1/audit-log",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionInstanceAdmin,
		)(
			h.AuditLog.GetAuditLogByQueryHandler,
		),
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/instance-metadata",
		h.InstanceMetadata.GetInstanceMetadataHandler,
	)

	// WebSocket
	publicRouter.HandleFunc("GET /api/ws", controllers.WebSocketConnectionHandler)

	// Health
	protectedRouter.HandleFunc("GET /api/v1/health", controllers.GetHealthCheckHandler)

	// SPA Handler for UI
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}

// requirePermission builds the permission middleware factory handed to feature routers.
func requirePermission(h Handlers) func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(perm auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
		return users.RequirePermission(h.AuthzMiddlewareService, perm)
	}
}
