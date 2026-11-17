package notifications

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

func RegisterRoutes(
	router *http.ServeMux,
	c NotificationChannelAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"POST /api/v1/notification-channels",
		requirePermission(auth.PermissionWriter)(c.CreateChannelHandler),
	)
	router.HandleFunc(
		"GET /api/v1/notification-channels",
		requirePermission(auth.PermissionReader)(c.GetAllChannelsHandler),
	)
	router.HandleFunc(
		"GET /api/v1/notification-channels/{channelId}",
		requirePermission(auth.PermissionReader)(c.GetChannelByIDHandler),
	)
	router.HandleFunc(
		"PATCH /api/v1/notification-channels/{channelId}",
		requirePermission(auth.PermissionWriter)(c.UpdateChannelHandler),
	)
	router.HandleFunc(
		"DELETE /api/v1/notification-channels/{channelId}",
		requirePermission(auth.PermissionWriter)(c.DeleteChannelHandler),
	)
	router.HandleFunc(
		"GET /api/v1/notification-channels/{channelId}/export",
		requirePermission(auth.PermissionReader)(c.ExportToYAMLHandler),
	)
}
