package controllers

import (
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
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

	util.RespondJSON(ctx, w, http.StatusOK, metadata)
}
