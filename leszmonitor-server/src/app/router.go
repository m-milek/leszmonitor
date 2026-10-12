package app

import (
	"embed"
	"errors"
	"net/http"

	"github.com/m-milek/leszmonitor/apidocs"
	"github.com/m-milek/leszmonitor/features/auditlog"
	"github.com/m-milek/leszmonitor/features/globalparameters"
	"github.com/m-milek/leszmonitor/features/instance"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/monitors/statuschange"
	"github.com/m-milek/leszmonitor/features/realtime"
	"github.com/m-milek/leszmonitor/features/users"

	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/httpx"
)

func SetupRouters(
	router *http.ServeMux,
	staticFiles embed.FS,
	h Handlers,
) {
	// Users
	users.RegisterRoutes(router, h.User, auth.RequireAuth)

	// Monitors
	monitors.RegisterRoutes(
		router,
		h.Monitor,
		h.MonitorResults,
		requirePermission(h),
	)
	stats.RegisterRoutes(router, h.MonitorStats, requirePermission(h))
	statuschange.RegisterRoutes(router, h.MonitorStatusChange, requirePermission(h))

	// Tags
	tags.RegisterRoutes(router, h.Tag, requirePermission(h))

	// Audit log
	auditlog.RegisterRoutes(router, h.AuditLog, requirePermission(h))

	// Global parameters
	globalparameters.RegisterRoutes(router, h.GlobalParameter, requirePermission(h))

	// Instance metadata
	instance.RegisterRoutes(router, h.InstanceMetadata, auth.RequireAuth)

	// WebSocket
	realtime.RegisterRoutes(router)

	// Health
	router.HandleFunc("GET /api/v1/health", GetHealthCheckHandler)

	// API spec
	router.HandleFunc("GET /api/v1/swagger.json", auth.RequireAuth(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(constants.HTTPHeaderContentType, constants.HTTPContentTypeJSON)
		_, _ = w.Write(apidocs.Spec)
	}))

	// Unknown API paths must not fall through to the SPA
	router.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {
		httpx.RespondError(r.Context(), w, http.StatusNotFound, errors.New("not found"))
	})

	// SPA Handler for UI
	router.Handle("/", newSPAHandler(staticFiles))
}

// requirePermission builds the auth + permission middleware factory handed to feature routers.
func requirePermission(h Handlers) func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
	return func(perm auth.Permission) func(http.HandlerFunc) http.HandlerFunc {
		check := users.RequirePermission(h.AuthzMiddlewareService, perm)
		return func(next http.HandlerFunc) http.HandlerFunc {
			return auth.RequireAuth(check(next))
		}
	}
}
