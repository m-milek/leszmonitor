package tags

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

// RegisterRoutes registers the tag API routes.
func RegisterRoutes(
	router *http.ServeMux,
	c TagAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"POST /api/v1/tags",
		requirePermission(auth.PermissionWriter)(c.CreateTagHandler),
	)
	router.HandleFunc(
		"GET /api/v1/tags",
		requirePermission(auth.PermissionReader)(c.GetAllTagsHandler),
	)
	router.HandleFunc(
		"GET /api/v1/tags/{tagId}",
		requirePermission(auth.PermissionReader)(c.GetTagByIDHandler),
	)
	router.HandleFunc(
		"PATCH /api/v1/tags/{tagId}",
		requirePermission(auth.PermissionWriter)(c.UpdateTagHandler),
	)
	router.HandleFunc(
		"DELETE /api/v1/tags/{tagId}",
		requirePermission(auth.PermissionWriter)(c.DeleteTagHandler),
	)
}
