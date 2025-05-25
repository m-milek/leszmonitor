package api

import (
	"github.com/m-milek/leszmonitor/api/handlers"
	"net/http"
)

func SetupRouters(
	publicRouter *http.ServeMux,
	protectedRouter *http.ServeMux,
) {
	protectedRouter.HandleFunc("/health", handlers.GetHealthCheckHandler)

	publicRouter.HandleFunc("/auth/register", handlers.UserRegisterHandler)
	publicRouter.HandleFunc("/auth/login", handlers.LoginHandler)
}
