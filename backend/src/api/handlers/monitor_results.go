package handlers

import (
	"errors"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/services"
	util2 "github.com/m-milek/leszmonitor/util"
)

func GetLatestMonitorResultByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceMonitorID)
	if !ok {
		return
	}

	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondError(ctx, w, http.StatusBadRequest, errors.New("monitor ID is required"))
		return
	}

	result, err := services.MonitorResultsService.GetLatestMonitorResultByMonitorID(ctx, projectAuth, monitorID)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, result)
}

func GetMonitorResultsByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceMonitorID)
	if !ok {
		return
	}

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

	results, err := services.MonitorResultsService.GetMonitorResultsByMonitorID(ctx, projectAuth, monitorID, pagination)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, results)
}
