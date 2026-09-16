package controllers

import (
	"net/http"

	"github.com/m-milek/leszmonitor/platform/httpx"
	"github.com/m-milek/leszmonitor/services"
)

type InstanceMetadataAPIController struct {
	service services.IInstanceMetadataService
}

func NewInstanceMetadataAPIController(service services.IInstanceMetadataService) InstanceMetadataAPIController {
	return InstanceMetadataAPIController{
		service: service,
	}
}

func (c *InstanceMetadataAPIController) GetInstanceMetadataHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metadata := c.service.GetInstanceMetadata()

	httpx.RespondJSON(ctx, w, http.StatusOK, metadata)
}
