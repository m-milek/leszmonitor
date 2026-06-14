package controllers

import (
	"errors"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/services"
	util2 "github.com/m-milek/leszmonitor/util"
)

type MonitorResultsAPIController struct {
	service services.IMonitorResultsService
}

func NewMonitorResultsAPIController(service services.IMonitorResultsService) MonitorResultsAPIController {
	return MonitorResultsAPIController{
		service: service,
	}
}

func (c *MonitorResultsAPIController) GetLatestMonitorResultByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondError(ctx, w, http.StatusBadRequest, errors.New("monitor ID is required"))
		return
	}

	result, svcErr := c.service.GetLatestMonitorResultByMonitorID(ctx, monitorID)
	if svcErr != nil {
		util.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, result)
}

func (c *MonitorResultsAPIController) GetMonitorResultsByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondError(ctx, w, http.StatusBadRequest, errors.New("monitor ID is required"))
		return
	}

	pagination, paginationErr := util2.PaginationFromRequest(r)
	if paginationErr != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, paginationErr)
		return
	}

	results, svcErr := c.service.GetMonitorResultsByMonitorID(ctx, monitorID, pagination)
	if svcErr != nil {
		util.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, results)
}
