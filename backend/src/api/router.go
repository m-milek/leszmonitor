package api

import (
	"net/http"

	"github.com/m-milek/leszmonitor/api/handlers"
)

func SetupRouters(
	publicRouter *http.ServeMux,
	protectedRouter *http.ServeMux,
) {
	protectedRouter.HandleFunc("GET /health", handlers.GetHealthCheckHandler)

	// Users
	protectedRouter.HandleFunc("GET /users", handlers.GetAllUsersHandler)
	protectedRouter.HandleFunc("GET /users/{username}", handlers.GetUserHandler)
	publicRouter.HandleFunc("POST /auth/register", handlers.UserRegisterHandler)
	publicRouter.HandleFunc("POST /auth/login", handlers.UserLoginHandler)

	// Projects
	protectedRouter.HandleFunc("GET /projects", handlers.GetProjectsHandler)
	protectedRouter.HandleFunc("POST /projects", handlers.CreateProjectHandler)
	protectedRouter.HandleFunc("GET /projects/{projectId}", handlers.GetProjectByIDHandler)
	protectedRouter.HandleFunc("PATCH /projects/{projectId}", handlers.UpdateProjectHandler)
	protectedRouter.HandleFunc("DELETE /projects/{projectId}", handlers.DeleteProjectHandler)

	// Project Members
	protectedRouter.HandleFunc("POST /projects/{projectId}/members", handlers.AddProjectMemberHandler)
	protectedRouter.HandleFunc("DELETE /projects/{projectId}/members", handlers.RemoveProjectMemberHandler)
	protectedRouter.HandleFunc("PATCH /projects/{projectId}/members/{userId}", handlers.ChangeProjectMemberRoleHandler)

	// Monitors
	protectedRouter.HandleFunc("GET /projects/{projectId}/monitors", handlers.GetAllMonitorsHandler)
	protectedRouter.HandleFunc("GET /projects/{projectId}/monitors/{monitorId}", handlers.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("POST /projects/{projectId}/monitors", handlers.CreateMonitorHandler)
	protectedRouter.HandleFunc("DELETE /projects/{projectId}/monitors/{monitorId}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /projects/{projectId}/monitors/{monitorId}", handlers.UpdateMonitorHandler)

	publicRouter.HandleFunc("GET /ws", handlers.WebSocketConnectionHandler)
}
