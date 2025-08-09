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

	protectedRouter.HandleFunc("GET /monitor", handlers.GetAllMonitorsHandler)
	protectedRouter.HandleFunc("POST /monitor", handlers.AddMonitorHandler)
	protectedRouter.HandleFunc("DELETE /monitor/{id}", handlers.DeleteMonitorHandler)
	protectedRouter.HandleFunc("PATCH /monitor/{id}", handlers.EditMonitorHandler)
	protectedRouter.HandleFunc("GET /monitor/{id}", handlers.GetMonitorHandler)

	protectedRouter.HandleFunc("GET /user", handlers.GetAllUsersHandler)
	protectedRouter.HandleFunc("GET /user/{username}", handlers.GetUserHandler)

	protectedRouter.HandleFunc("GET /team", handlers.GetAllTeamsHandler)
	protectedRouter.HandleFunc("GET /team/{id}", handlers.GetTeamHandler)
	protectedRouter.HandleFunc("POST /team", handlers.TeamCreateHandler)
	protectedRouter.HandleFunc("DELETE /team/{id}", handlers.TeamDeleteHandler)
	protectedRouter.HandleFunc("PATCH /team/{id}", handlers.TeamUpdateHandler)
	protectedRouter.HandleFunc("POST /team/{id}/add-member", handlers.TeamAddMemberHandler)
	protectedRouter.HandleFunc("POST /team/{id}/remove-member", handlers.TeamRemoveMemberHandler)
	protectedRouter.HandleFunc("POST /team/{id}/change-member-role", handlers.TeamChangeMemberRoleHandler)

	publicRouter.HandleFunc("POST /auth/register", handlers.UserRegisterHandler)
	publicRouter.HandleFunc("POST /auth/login", handlers.UserLoginHandler)

}
