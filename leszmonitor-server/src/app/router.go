package app

import (
	"embed"
	"net/http"

	"github.com/m-milek/leszmonitor/features/auditlog"
	"github.com/m-milek/leszmonitor/features/instance"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/realtime"
	"github.com/m-milek/leszmonitor/features/users"

	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/auth"
)

func SetupRouters(
	publicRouter *http.ServeMux,
	protectedRouter *http.ServeMux,
	staticFiles embed.FS,
	h Handlers,
) {
	// Users
	users.RegisterRoutes(publicRouter, protectedRouter, h.User)

	// Monitors
	monitors.RegisterRoutes(
		protectedRouter,
		h.Monitor,
		h.MonitorResults,
		h.MonitorStats,
		requirePermission(h),
	)

	// Tags
	tags.RegisterRoutes(protectedRouter, h.Tag, requirePermission(h))

	// Audit log
	auditlog.RegisterRoutes(protectedRouter, h.AuditLog, requirePermission(h))

	// Instance metadata
	instance.RegisterRoutes(protectedRouter, h.InstanceMetadata)

	// WebSocket
	realtime.RegisterRoutes(publicRouter)

	// Health
	protectedRouter.HandleFunc("GET /api/v1/health", GetHealthCheckHandler)

	// SPA Handler for UI
	publicRouter.Handle("/", newSPAHandler(staticFiles))
}

// requirePermission builds the permission middleware factory handed to feature routers.
func requirePermission(h Handlers) func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(perm auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
		return users.RequirePermission(h.AuthzMiddlewareService, perm)
	}
}
