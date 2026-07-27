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
	publicRouter.HandleFunc("POST /api/v1/auth/register", h.User.UserRegisterHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/login", h.User.UserLoginHandler)

	// Projects
	protectedRouter.HandleFunc("GET /api/v1/projects", h.Project.GetProjectsHandler)
	protectedRouter.HandleFunc("POST /api/v1/projects", h.Project.CreateProjectHandler)
	protectedRouter.HandleFunc(
		"GET /api/v1/projects/{projectSlug}",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectReader,
			middleware.SlugSourcePath,
		)(
			h.Project.GetProjectByIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/projects/{projectSlug}",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
			middleware.SlugSourcePath,
		)(
			h.Project.UpdateProjectHandler,
		),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/projects/{projectSlug}",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectAdmin,
			middleware.SlugSourcePath,
		)(
			h.Project.DeleteProjectHandler,
		),
	)

	// Project Members
	protectedRouter.HandleFunc(
		"POST /api/v1/projects/{projectSlug}/members",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
			middleware.SlugSourcePath,
		)(
			h.Project.AddProjectMemberHandler,
		),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/projects/{projectSlug}/members",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
			middleware.SlugSourcePath,
		)(
			h.Project.RemoveProjectMemberHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/projects/{projectSlug}/members/{userId}",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectAdmin,
			middleware.SlugSourcePath,
		)(
			h.Project.ChangeProjectMemberRoleHandler,
		),
	)

	// Monitors
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectReader,
			middleware.SlugSourceQuery,
		)(
			h.Monitor.GetMonitorByProjectSlugHandler,
		),
	)
	protectedRouter.HandleFunc(
		"POST /api/v1/monitors",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
			middleware.SlugSourceQuery,
		)(
			h.Monitor.CreateMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}",
		middleware.RequireMonitorPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectReader,
		)(
			h.Monitor.GetMonitorByIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/monitors/{monitorId}",
		middleware.RequireMonitorPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
		)(
			h.Monitor.DeleteMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}",
		middleware.RequireMonitorPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
		)(
			h.Monitor.UpdateMonitorHandler,
		),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}/state",
		middleware.RequireMonitorPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectEditor,
		)(
			h.Monitor.UpdateMonitorStateByIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/projects/{projectSlug}/monitors/{monitorSlug}",
		middleware.RequireProjectPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectReader,
			middleware.SlugSourcePath,
		)(
			h.Monitor.GetMonitorBySlugByProject,
		),
	)

	// MonitorResults
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results/latest",
		middleware.RequireMonitorPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectReader,
		)(
			h.MonitorResults.GetLatestMonitorResultByMonitorIDHandler,
		),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results",
		middleware.RequireMonitorPermission(
			h.AuthzMiddlewareService,
			models.PermissionProjectReader,
		)(
			h.MonitorResults.GetMonitorResultsByMonitorIDHandler,
		),
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/audit-log",
		middleware.RequireProjectPermissionByIDQuery(
			h.AuthzMiddlewareService,
			models.PermissionProjectAdmin,
		)(
			h.AuditLog.GetAuditLogByQueryHandler,
		),
	)

	// WebSocket
	publicRouter.HandleFunc("GET /api/ws", controllers.WebSocketConnectionHandler)

	// Health
	protectedRouter.HandleFunc("GET /api/v1/health", controllers.GetHealthCheckHandler)

	// SPA Handler for frontend
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}
