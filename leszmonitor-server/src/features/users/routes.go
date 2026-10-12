package users

import "net/http"

// RegisterRoutes registers the user and auth API routes.
func RegisterRoutes(
	router *http.ServeMux,
	c UserAPIController,
	requireAuth func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc("GET /api/v1/users", requireAuth(c.GetAllUsersHandler))
	router.HandleFunc("GET /api/v1/users/{username}", requireAuth(c.GetUserHandler))
	router.HandleFunc(
		"PATCH /api/v1/users/{username}/role",
		requireAuth(RequireInstanceAdmin()(c.SetUserRoleHandler)),
	)
	router.HandleFunc("POST /api/v1/auth/register", c.UserRegisterHandler)
	router.HandleFunc("POST /api/v1/auth/login", c.UserLoginHandler)
}
