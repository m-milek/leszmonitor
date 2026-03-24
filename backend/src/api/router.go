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

	// Orgs
	protectedRouter.HandleFunc("GET /orgs", handlers.GetAllOrgsHandler)
	protectedRouter.HandleFunc("GET /orgs/{orgId}", handlers.GetOrgHandler)
	protectedRouter.HandleFunc("POST /orgs", handlers.CreateOrgHandler)
	protectedRouter.HandleFunc("DELETE /orgs/{orgId}", handlers.DeleteOrgHandler)
	protectedRouter.HandleFunc("PATCH /orgs/{orgId}", handlers.UpdateOrgHandler)

	// Org Members
	protectedRouter.HandleFunc("POST /orgs/{orgId}/members", handlers.AddOrgMemberHandler)
	protectedRouter.HandleFunc("DELETE /orgs/{orgId}/members", handlers.RemoveOrgMemberHandler)
	protectedRouter.HandleFunc("PATCH /orgs/{orgId}/members/{userId}", handlers.ChangeOrgMemberRoleHandler)

	// Monitors
	protectedRouter.HandleFunc("GET /orgs/{orgId}/monitors", handlers.GetAllMonitorsHandler)
	protectedRouter.HandleFunc("GET /orgs/{orgId}/monitors/{monitorId}", handlers.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("POST /orgs/{orgId}/monitors", handlers.CreateMonitorHandler)
	protectedRouter.HandleFunc("DELETE /orgs/{orgId}/monitors/{monitorId}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /orgs/{orgId}/monitors/{monitorId}", handlers.UpdateMonitorHandler)

	// Projects
	protectedRouter.HandleFunc("POST /orgs/{orgId}/projects", handlers.CreateProjectHandler)
	protectedRouter.HandleFunc("GET /orgs/{orgId}/projects", handlers.GetProjectsOfOrgHandler)
	protectedRouter.HandleFunc("GET /orgs/{orgId}/projects/{projectId}", handlers.GetProjectsByOrgID)
	protectedRouter.HandleFunc("PATCH /orgs/{orgId}/projects/{projectId}", handlers.UpdateProjectHandler)
	protectedRouter.HandleFunc("DELETE /orgs/{orgId}/projects/{projectId}", handlers.DeleteProjectHandler)

}
