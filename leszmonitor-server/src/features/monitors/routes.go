package monitors

import (
	"net/http"

	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/platform/auth"
)

// RegisterRoutes registers the monitor, monitor result and monitor stats API routes.
func RegisterRoutes(
	protectedRouter *http.ServeMux,
	c MonitorAPIController,
	results results.MonitorResultsAPIController,
	stats stats.MonitorStatsAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	protectedRouter.HandleFunc(
		"POST /api/v1/monitors",
		requirePermission(auth.PermissionWriter)(c.CreateMonitorHandler),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors",
		requirePermission(auth.PermissionReader)(c.GetAllMonitorsHandler),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}",
		requirePermission(auth.PermissionReader)(c.GetMonitorByIDHandler),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/monitors/{monitorId}",
		requirePermission(auth.PermissionWriter)(c.DeleteMonitorHandler),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}",
		requirePermission(auth.PermissionWriter)(c.UpdateMonitorHandler),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}/state",
		requirePermission(auth.PermissionWriter)(c.UpdateMonitorStateByIDHandler),
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results/latest",
		requirePermission(auth.PermissionReader)(results.GetLatestMonitorResultByMonitorIDHandler),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results",
		requirePermission(auth.PermissionReader)(results.GetMonitorResultsByMonitorIDHandler),
	)

	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/stats",
		requirePermission(auth.PermissionReader)(stats.GetLatencyStatsByMonitorIDHandler),
	)
}
