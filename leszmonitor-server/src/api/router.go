package api

import (
	"embed"
	"net/http"

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
	protectedRouter.HandleFunc(
		"POST /api/v1/monitors",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionWriter,
		)(
			h.Monitor.CreateMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionReader,
		)(
			h.Monitor.GetAllMonitorsHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionReader,
		)(
			h.Monitor.GetMonitorByIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/monitors/{monitorId}",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionWriter,
		)(
			h.Monitor.DeleteMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionWriter,
		)(
			h.Monitor.UpdateMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}/state",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionWriter,
		)(
			h.Monitor.UpdateMonitorStateByIDHandler,
		),
	)

	// Tags
	tags.RegisterRoutes(protectedRouter, h.Tag, func(perm auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
		return users.RequirePermission(h.AuthzMiddlewareService, perm)
	})

	// MonitorResults
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results/latest",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionReader,
		)(
			h.MonitorResults.GetLatestMonitorResultByMonitorIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionReader,
		)(
			h.MonitorResults.GetMonitorResultsByMonitorIDHandler,
		),
	)

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

	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/stats",
		users.RequirePermission(
			h.AuthzMiddlewareService,
			auth.PermissionReader,
		)(
			h.MonitorStats.GetLatencyStatsByMonitorIDHandler,
		),
	)

	// WebSocket
	publicRouter.HandleFunc("GET /api/ws", controllers.WebSocketConnectionHandler)

	// Health
	protectedRouter.HandleFunc("GET /api/v1/health", controllers.GetHealthCheckHandler)

	// SPA Handler for UI
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}
