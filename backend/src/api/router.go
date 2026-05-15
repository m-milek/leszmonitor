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
	protectedRouter.HandleFunc("GET /api/health", handlers.GetHealthCheckHandler)

	// Users
	protectedRouter.HandleFunc("GET /api/users", handlers.GetAllUsersHandler)
	protectedRouter.HandleFunc("GET /api/users/{username}", handlers.GetUserHandler)
	publicRouter.HandleFunc("POST /api/auth/register", handlers.UserRegisterHandler)
	publicRouter.HandleFunc("POST /api/auth/login", handlers.UserLoginHandler)

	// Projects
	protectedRouter.HandleFunc("GET /api/projects", handlers.GetProjectsHandler)
	protectedRouter.HandleFunc("POST /api/projects", handlers.CreateProjectHandler)
	protectedRouter.HandleFunc("GET /api/projects/{projectSlug}", handlers.GetProjectByIDHandler)
	protectedRouter.HandleFunc("PATCH /api/projects/{projectSlug}", handlers.UpdateProjectHandler)
	protectedRouter.HandleFunc("DELETE /api/projects/{projectSlug}", handlers.DeleteProjectHandler)

	// Project Members
	protectedRouter.HandleFunc("POST /api/projects/{projectSlug}/members", handlers.AddProjectMemberHandler)
	protectedRouter.HandleFunc("DELETE /api/projects/{projectSlug}/members", handlers.RemoveProjectMemberHandler)
	protectedRouter.HandleFunc("PATCH /api/projects/{projectSlug}/members/{userId}", handlers.ChangeProjectMemberRoleHandler)

	// Monitors
	protectedRouter.HandleFunc("GET /api/monitors", handlers.GetMonitorByProjectSlugHandler)
	protectedRouter.HandleFunc("POST /api/monitors", handlers.CreateMonitorHandler)
	protectedRouter.HandleFunc("GET /api/monitors/{monitorId}", handlers.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("DELETE /api/monitors/{monitorId}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /api/monitors/{monitorId}", handlers.UpdateMonitorHandler)
	protectedRouter.HandleFunc("GET /api/projects/{projectSlug}/monitors/{monitorSlug}", handlers.GetMonitorBySlugByProject)

	// Monitor Results
	protectedRouter.HandleFunc("GET /api/monitors/{monitorId}/results/latest", handlers.GetLatestMonitorResultByMonitorIDHandler)
	protectedRouter.HandleFunc("GET /api/monitors/{monitorId}/results", handlers.GetMonitorResultsByMonitorIDHandler)

	// WebSocket
	publicRouter.HandleFunc("GET /api/ws", handlers.WebSocketConnectionHandler)

	// SPA Handler for frontend
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}
