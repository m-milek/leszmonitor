package tags

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/auth"
)

// RegisterRoutes registers the tag API routes on the protected router.
func RegisterRoutes(
	protectedRouter *http.ServeMux,
	c TagAPIController,
	requirePermission func(auth.Permission) func(http.HandlerFunc) http.HandlerFunc,
) {
	protectedRouter.HandleFunc(
		"POST /api/v1/tags",
		requirePermission(auth.PermissionWriter)(c.CreateTagHandler),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/tags",
		requirePermission(auth.PermissionReader)(c.GetAllTagsHandler),
	)
	protectedRouter.HandleFunc(
		"GET /api/v1/tags/{tagId}",
		requirePermission(auth.PermissionReader)(c.GetTagByIDHandler),
	)
	protectedRouter.HandleFunc(
		"PATCH /api/v1/tags/{tagId}",
		requirePermission(auth.PermissionWriter)(c.UpdateTagHandler),
	)
	protectedRouter.HandleFunc(
		"DELETE /api/v1/tags/{tagId}",
		requirePermission(auth.PermissionWriter)(c.DeleteTagHandler),
	)
}
