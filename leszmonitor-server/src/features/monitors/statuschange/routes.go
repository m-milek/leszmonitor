package statuschange

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	protectedRouter *http.ServeMux,
	c MonitorStatusChangeAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	protectedRouter.HandleFunc(
		"GET /api/v1/monitors/{monitorId}/statusHistory",
		requirePermission(auth.PermissionReader)(c.GetStatusHistoryByMonitorIDHandler),
	)
}
