package api

import (
	"embed"
	"net/http"

	"github.com/m-milek/leszmonitor/api/handlers"
)

func SetupRouters(
	publicRouter *http.ServeMux,
	protectedRouter *http.ServeMux,
	staticFiles embed.FS,
) {
	protectedRouter.HandleFunc("GET /api/v1/health", handlers.GetHealthCheckHandler)

	// Users
	protectedRouter.HandleFunc("GET /api/v1/users", handlers.GetAllUsersHandler)
	protectedRouter.HandleFunc("GET /api/v1/users/{username}", handlers.GetUserHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/register", handlers.UserRegisterHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/login", handlers.UserLoginHandler)

	// Projects
	protectedRouter.HandleFunc("GET /api/v1/projects", handlers.GetProjectsHandler)
	protectedRouter.HandleFunc("POST /api/v1/projects", handlers.CreateProjectHandler)
	protectedRouter.HandleFunc("GET /api/v1/projects/{projectSlug}", handlers.GetProjectByIDHandler)
	protectedRouter.HandleFunc("PATCH /api/v1/projects/{projectSlug}", handlers.UpdateProjectHandler)
	protectedRouter.HandleFunc("DELETE /api/v1/projects/{projectSlug}", handlers.DeleteProjectHandler)

	// Project Members
	protectedRouter.HandleFunc("POST /api/v1/projects/{projectSlug}/members", handlers.AddProjectMemberHandler)
	protectedRouter.HandleFunc("DELETE /api/v1/projects/{projectSlug}/members", handlers.RemoveProjectMemberHandler)
	protectedRouter.HandleFunc("PATCH /api/v1/projects/{projectSlug}/members/{userId}", handlers.ChangeProjectMemberRoleHandler)

	// Monitors
	protectedRouter.HandleFunc("GET /api/v1/monitors", handlers.GetMonitorByProjectSlugHandler)
	protectedRouter.HandleFunc("POST /api/v1/monitors", handlers.CreateMonitorHandler)
	protectedRouter.HandleFunc("GET /api/v1/monitors/{monitorId}", handlers.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("DELETE /api/v1/monitors/{monitorId}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /api/v1/monitors/{monitorId}", handlers.UpdateMonitorHandler)
	protectedRouter.HandleFunc("GET /api/v1/projects/{projectSlug}/monitors/{monitorSlug}", handlers.GetMonitorBySlugByProject)

	// Monitor Results
	protectedRouter.HandleFunc("GET /api/v1/monitors/{monitorId}/results/latest", handlers.GetLatestMonitorResultByMonitorIDHandler)
	protectedRouter.HandleFunc("GET /api/v1/monitors/{monitorId}/results", handlers.GetMonitorResultsByMonitorIDHandler)

	// WebSocket
	publicRouter.HandleFunc("GET /api/ws", handlers.WebSocketConnectionHandler)

	// SPA Handler for frontend
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}
