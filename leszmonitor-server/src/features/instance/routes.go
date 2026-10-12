package instance

import "net/http"

// RegisterRoutes registers the instance metadata API routes.
func RegisterRoutes(
	router *http.ServeMux,
	c InstanceMetadataAPIController,
	requireAuth func(http.HandlerFunc) http.HandlerFunc,
) {
	router.HandleFunc(
		"GET /api/v1/instance-metadata",
		requireAuth(c.GetInstanceMetadataHandler),
	)
}
