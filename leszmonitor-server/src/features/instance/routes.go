package instance

import "net/http"

// RegisterRoutes registers the instance metadata API routes on the protected router.
func RegisterRoutes(protectedRouter *http.ServeMux, c InstanceMetadataAPIController) {
	protectedRouter.HandleFunc(
		"GET /api/v1/instance-metadata",
		c.GetInstanceMetadataHandler,
	)
}
