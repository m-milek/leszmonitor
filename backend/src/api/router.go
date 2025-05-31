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
	protectedRouter.HandleFunc("POST /monitor", handlers.AddMonitor)

	publicRouter.HandleFunc("POST /auth/register", handlers.UserRegisterHandler)
	publicRouter.HandleFunc("POST /auth/login", handlers.LoginHandler)
}
