package statuschange

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	router *http.ServeMux,
	c MonitorStatusChangeAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/status-history",
		requirePermission(auth.PermissionReader)(c.GetStatusHistoryByMonitorIDHandler),
	)
}
