package instance

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
)

type InstanceMetadataAPIController struct {
	service IInstanceMetadataService
}

func NewInstanceMetadataAPIController(service IInstanceMetadataService) InstanceMetadataAPIController {
	return InstanceMetadataAPIController{
		service: service,
	}
}

// @Summary	Get instance metadata
// @Tags		instance
// @Produce	json
// @Success	200	{object}	InstanceMetadata
// @Failure	401	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/instance-metadata [get]
func (c *InstanceMetadataAPIController) GetInstanceMetadataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metadata := c.service.GetInstanceMetadata()

	httpx.RespondJSON(ctx, w, http.StatusOK, metadata)
}
