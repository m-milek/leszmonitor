package stats

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	router *http.ServeMux,
	c MonitorStatsAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/stats",
		requirePermission(auth.PermissionReader)(c.GetStatsByMonitorIDHandler),
	)
}
