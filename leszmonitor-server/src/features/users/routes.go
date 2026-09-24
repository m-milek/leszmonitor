package users

import "net/http"

// RegisterRoutes registers the user and auth API routes.
func RegisterRoutes(publicRouter *http.ServeMux, protectedRouter *http.ServeMux, c UserAPIController) {
	protectedRouter.HandleFunc("GET /api/v1/users", c.GetAllUsersHandler)
	protectedRouter.HandleFunc(
		"GET /api/v1/users/{username}", c.GetUserHandler,
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/users/{username}/role",
		RequireInstanceAdmin()(c.SetUserRoleHandler),
	)
	publicRouter.HandleFunc("POST /api/v1/auth/register", c.UserRegisterHandler)
	publicRouter.HandleFunc("POST /api/v1/auth/login", c.UserLoginHandler)
}
