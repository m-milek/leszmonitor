package api

import (
	"github.com/m-milek/leszmonitor/api/handlers"
	"net/http"
)

func SetupRouter(
	router *http.ServeMux,
) {
	router.HandleFunc("/health", handlers.GetHealthCheckHandler)
	router.HandleFunc("/user/register", handlers.UserRegisterHandler)
}
