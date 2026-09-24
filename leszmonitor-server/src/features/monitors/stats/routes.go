package stats

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	protectedRouter *http.ServeMux,
	c MonitorStatsAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/stats",
		requirePermission(auth.PermissionReader)(c.GetStatsByMonitorIDHandler),
	)
}
