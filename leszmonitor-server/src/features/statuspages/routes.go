package statuspages

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	router *http.ServeMux,
	c StatusPageAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"POST /api/v1/status-pages",
		requirePermission(auth.PermissionWriter)(c.CreateStatusPageHandler),
	)
	router.HandleFunc(
		"GET /api/v1/status-pages",
		requirePermission(auth.PermissionReader)(c.GetAllStatusPagesHandler),
	)
	router.HandleFunc(
		"GET /api/v1/status-pages/{statusPageId}",
		requirePermission(auth.PermissionReader)(c.GetStatusPageByIDHandler),
	)
	router.HandleFunc(
		"GET /api/v1/public/status-pages/{slug}",
		c.GetPublicStatusPageHandler(requirePermission(auth.PermissionReader)),
	)
	router.HandleFunc(
		"GET /api/v1/public/status-pages/{slug}/monitors",
		c.GetPublicMonitorsHandler(requirePermission(auth.PermissionReader)),
	)
	router.HandleFunc(
		"PATCH /api/v1/status-pages/{statusPageId}",
		requirePermission(auth.PermissionWriter)(c.UpdateStatusPageHandler),
	)
	router.HandleFunc(
		"DELETE /api/v1/status-pages/{statusPageId}",
		requirePermission(auth.PermissionWriter)(c.DeleteStatusPageHandler),
	)
}
