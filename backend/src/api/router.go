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
	protectedRouter.HandleFunc("GET /api/v1/users/{username}", h.User.GetUserHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/register", h.User.UserRegisterHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/login", h.User.UserLoginHandler)

	// Projects
	protectedRouter.HandleFunc("GET /api/v1/projects", h.Project.GetProjectsHandler)
	protectedRouter.HandleFunc("POST /api/v1/projects", h.Project.CreateProjectHandler)
	protectedRouter.HandleFunc("GET /api/v1/projects/{projectSlug}", middleware.RequireProjectPermission(h.AuthzMiddlewareService, models.PermissionProjectReader)(h.Project.GetProjectByIDHandler))
	protectedRouter.HandleFunc("PATCH /api/v1/projects/{projectSlug}", middleware.RequireProjectPermission(h.AuthzMiddlewareService, models.PermissionProjectEditor)(h.Project.UpdateProjectHandler))
	protectedRouter.HandleFunc("DELETE /api/v1/projects/{projectSlug}", middleware.RequireProjectPermission(h.AuthzMiddlewareService, models.PermissionProjectAdmin)(h.Project.DeleteProjectHandler))

	// Project Members
	protectedRouter.HandleFunc("POST /api/v1/projects/{projectSlug}/members", middleware.RequireProjectPermission(h.AuthzMiddlewareService, models.PermissionProjectEditor)(h.Project.AddProjectMemberHandler))
	protectedRouter.HandleFunc("DELETE /api/v1/projects/{projectSlug}/members", middleware.RequireProjectPermission(h.AuthzMiddlewareService, models.PermissionProjectEditor)(h.Project.RemoveProjectMemberHandler))
	protectedRouter.HandleFunc("PATCH /api/v1/projects/{projectSlug}/members/{userId}", middleware.RequireProjectPermission(h.AuthzMiddlewareService, models.PermissionProjectAdmin)(h.Project.ChangeProjectMemberRoleHandler))

	// Monitors
	protectedRouter.HandleFunc("GET /api/v1/monitors", h.Monitor.GetMonitorByProjectSlugHandler)
	protectedRouter.HandleFunc("POST /api/v1/monitors", h.Monitor.CreateMonitorHandler)
	protectedRouter.HandleFunc("GET /api/v1/monitors/{monitorId}", h.Monitor.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("DELETE /api/v1/monitors/{monitorId}", h.Monitor.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /api/v1/monitors/{monitorId}", h.Monitor.UpdateMonitorHandler)
	protectedRouter.HandleFunc("PATCH /api/v1/monitors/{monitorId}/state", h.Monitor.UpdateMonitorStateByIDHandler)
	protectedRouter.HandleFunc("GET /api/v1/projects/{projectSlug}/monitors/{monitorSlug}", h.Monitor.GetMonitorBySlugByProject)

	// Monitor Results
	protectedRouter.HandleFunc("GET /api/v1/monitors/{monitorId}/results/latest", h.MonitorResults.GetLatestMonitorResultByMonitorIDHandler)
	protectedRouter.HandleFunc("GET /api/v1/monitors/{monitorId}/results", h.MonitorResults.GetMonitorResultsByMonitorIDHandler)

	protectedRouter.HandleFunc("GET /api/v1/audit-log", h.AuditLog.GetAuditLogByQueryHandler)

	// WebSocket
	publicRouter.HandleFunc("GET /api/ws", controllers.WebSocketConnectionHandler)

	// Health
	protectedRouter.HandleFunc("GET /api/v1/health", controllers.GetHealthCheckHandler)

	// SPA Handler for frontend
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}
