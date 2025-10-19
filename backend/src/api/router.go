package api

import (
	"github.com/m-milek/leszmonitor/api/handlers"
	"net/http"
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

	// Teams
	protectedRouter.HandleFunc("GET /teams", handlers.GetAllTeamsHandler)
	protectedRouter.HandleFunc("GET /teams/{teamId}", handlers.GetTeamHandler)
	protectedRouter.HandleFunc("POST /teams", handlers.TeamCreateHandler)
	protectedRouter.HandleFunc("DELETE /teams/{teamId}", handlers.TeamDeleteHandler)
	protectedRouter.HandleFunc("PATCH /teams/{teamId}", handlers.TeamUpdateHandler)

	// Team Members
	protectedRouter.HandleFunc("POST /teams/{teamId}/members", handlers.TeamAddMemberHandler)
	protectedRouter.HandleFunc("DELETE /teams/{teamId}/members", handlers.TeamRemoveMemberHandler)
	protectedRouter.HandleFunc("PATCH /teams/{teamId}/{userId}", handlers.TeamChangeMemberRoleHandler)

	// Monitor Groups
	protectedRouter.HandleFunc("POST /teams/{teamId}/groups", handlers.CreateMonitorGroupHandler)
	protectedRouter.HandleFunc("GET /teams/{teamId}/groups", handlers.GetTeamMonitorGroupsHandler)
	protectedRouter.HandleFunc("GET /teams/{teamId}/groups/{groupId}", handlers.GetTeamMonitorGroupByID)
	protectedRouter.HandleFunc("PATCH /teams/{teamId}/groups/{groupId}", handlers.UpdateMonitorGroupHandler)
	protectedRouter.HandleFunc("DELETE /teams/{teamId}/groups/{groupId}", handlers.DeleteMonitorGroupHandler)

	// Monitors
	protectedRouter.HandleFunc("GET /teams/{teamId}/groups/{groupId}/monitors", handlers.GetAllMonitorsHandler)
	protectedRouter.HandleFunc("GET /teams/{teamId}/groups/{groupId}/monitors/{monitorId}", handlers.GetMonitorByIDHandler)
	protectedRouter.HandleFunc("POST /teams/{teamId}/groups/{groupId}/monitors", handlers.CreateMonitorHandler)
	protectedRouter.HandleFunc("DELETE /teams/{teamId}/groups/{groupId}/monitors/{monitorId}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /teams/{teamId}/groups/{groupId}/monitors/{monitorId}", handlers.UpdateMonitorHandler)

}
