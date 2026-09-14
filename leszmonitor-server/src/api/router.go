package api

import (
	"embed"
	"net/http"

	"github.com/m-milek/leszmonitor/api/controllers"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/models"
)

func SetupRouters(
	publicRouter *http.ServeMux,
	protectedRouter *http.ServeMux,
	staticFiles embed.FS,
	h Handlers,
) {
	// Users
	protectedRouter.HandleFunc("GET /api/v1/users", h.User.GetAllUsersHandler)
	protectedRouter.HandleFunc(
		"GET /api/v1/users/{username}",
		middleware.RequireSelf("username")(h.User.GetUserHandler),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/users/{username}/role",
		middleware.RequireInstanceAdmin()(h.User.SetUserRoleHandler),
	)
	publicRouter.HandleFunc("POST /api/v1/auth/register", h.User.UserRegisterHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/login", h.User.UserLoginHandler)

	// Monitors
	protectedRouter.HandleFunc(
		"POST /api/v1/monitors",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionWriter,
		)(
			h.Monitor.CreateMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionReader,
		)(
			h.Monitor.GetMonitorByIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/monitors/{monitorId}",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionWriter,
		)(
			h.Monitor.DeleteMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionWriter,
		)(
			h.Monitor.UpdateMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}/state",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionWriter,
		)(
			h.Monitor.UpdateMonitorStateByIDHandler,
		),
	)

	// MonitorResults
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results/latest",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionReader,
		)(
			h.MonitorResults.GetLatestMonitorResultByMonitorIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionReader,
		)(
			h.MonitorResults.GetMonitorResultsByMonitorIDHandler,
		),
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/audit-log",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionInstanceAdmin,
		)(
			h.AuditLog.GetAuditLogByQueryHandler,
		),
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/instance-metadata",
		h.InstanceMetadata.GetInstanceMetadataHandler,
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/stats/latency",
		middleware.RequirePermission(
			h.AuthzMiddlewareService,
			models.PermissionReader,
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
