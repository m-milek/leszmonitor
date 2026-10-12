package globalparameters

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	router *http.ServeMux,
	c GlobalParameterAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"GET /api/v1/global-parameters",
		requirePermission(auth.PermissionReader)(c.GetAllParametersHandler),
	)
	router.HandleFunc(
		"PATCH /api/v1/global-parameters",
		requirePermission(auth.PermissionInstanceAdmin)(c.SetParametersHandler),
	)
}
