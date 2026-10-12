package monitors

import (
	"net/http"

	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/auth"
)

// RegisterRoutes registers the monitor and monitor result API routes.
func RegisterRoutes(
	router *http.ServeMux,
	c MonitorAPIController,
	results results.MonitorResultsAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"POST /api/v1/monitors",
		requirePermission(auth.PermissionWriter)(c.CreateMonitorHandler),
	)
	router.HandleFunc(
		"GET /api/v1/monitors",
		requirePermission(auth.PermissionReader)(c.GetAllMonitorsHandler),
	)
	router.HandleFunc(
		"GET /api/v1/monitors/{monitorId}",
		requirePermission(auth.PermissionReader)(c.GetMonitorByIDHandler),
	)
	router.HandleFunc(
		"DELETE /api/v1/monitors/{monitorId}",
		requirePermission(auth.PermissionWriter)(c.DeleteMonitorHandler),
	)
	router.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}",
		requirePermission(auth.PermissionWriter)(c.UpdateMonitorHandler),
	)
	router.HandleFunc(
		"PATCH /api/v1/monitors/{monitorId}/state",
		requirePermission(auth.PermissionWriter)(c.UpdateMonitorStateByIDHandler),
	)
	router.HandleFunc(
		"POST /api/v1/monitors/{monitorId}/run",
		requirePermission(auth.PermissionWriter)(c.RunMonitorManuallyByIDHandler),
	)
	router.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/export",
		requirePermission(auth.PermissionReader)(c.ExportToYAMLHandler),
	)

	router.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results/latest",
		requirePermission(auth.PermissionReader)(results.GetLatestMonitorResultByMonitorIDHandler),
	)
	router.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/results",
		requirePermission(auth.PermissionReader)(results.GetMonitorResultsByMonitorIDHandler),
	)

	router.HandleFunc("GET /api/v1/push/{monitorId}", c.ReceivePushHandler)
	router.HandleFunc("POST /api/v1/push/{monitorId}", c.ReceivePushHandler)
}
