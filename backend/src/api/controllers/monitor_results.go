package controllers

import (
	"errors"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
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

	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		MonitorID: monitorID,
	})
	if !ok {
		return
	}

	result, err := c.service.GetLatestMonitorResultByMonitorID(ctx, projectAuth, monitorID)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
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

	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		MonitorID: monitorID,
	})
	if !ok {
		return
	}

	results, err := c.service.GetMonitorResultsByMonitorID(ctx, projectAuth, monitorID, pagination)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, results)
}
