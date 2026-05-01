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
	protectedRouter.HandleFunc("GET /api/projects/{projectId}", handlers.GetProjectByIDHandler)
	protectedRouter.HandleFunc("PATCH /api/projects/{projectId}", handlers.UpdateProjectHandler)
	protectedRouter.HandleFunc("DELETE /api/projects/{projectId}", handlers.DeleteProjectHandler)

	// Project Members
	protectedRouter.HandleFunc("POST /api/projects/{projectId}/members", handlers.AddProjectMemberHandler)
	protectedRouter.HandleFunc("DELETE /api/projects/{projectId}/members", handlers.RemoveProjectMemberHandler)
	protectedRouter.HandleFunc("PATCH /api/projects/{projectId}/members/{userId}", handlers.ChangeProjectMemberRoleHandler)

	// Monitors
	protectedRouter.HandleFunc("GET /api/projects/{projectId}/monitors", handlers.GetAllMonitorsHandler)
	protectedRouter.HandleFunc("GET /api/projects/{projectId}/monitors/{monitorId}", handlers.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("POST /api/projects/{projectId}/monitors", handlers.CreateMonitorHandler)
	protectedRouter.HandleFunc("DELETE /api/projects/{projectId}/monitors/{monitorId}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /api/projects/{projectId}/monitors/{monitorId}", handlers.UpdateMonitorHandler)

	publicRouter.HandleFunc("GET /api/ws", handlers.WebSocketConnectionHandler)

	publicRouter.Handle("/", newSPAHandler(staticFiles))
}
