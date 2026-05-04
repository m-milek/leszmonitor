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
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceKindMonitor)
	if !ok {
		return
	}

	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondError(w, http.StatusBadRequest, errors.New("monitor ID is required"))
		return
	}

	result, err := services.MonitorResultsService.GetLatestMonitorResultByMonitorID(r.Context(), projectAuth, monitorID)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, result)
}

func GetMonitorResultsByMonitorIDHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceKindMonitor)
	if !ok {
		return
	}

	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondError(w, http.StatusBadRequest, errors.New("monitor ID is required"))
		return
	}

	pagination, paginationErr := util2.PaginationFromRequest(r)
	if paginationErr != nil {
		util.RespondError(w, http.StatusBadRequest, paginationErr)
		return
	}

	results, err := services.MonitorResultsService.GetMonitorResultsByMonitorID(r.Context(), projectAuth, monitorID, pagination)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, results)
}
