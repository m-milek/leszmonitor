package handlers

import (
	"github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/uptime/monitor"
	"net/http"
)

// CreateMonitorHandler handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
func CreateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("groupId")
	if groupId == "" {
		logging.Api.Trace().Msg("Group DisplayId is required for creating a monitor")
		util.RespondMessage(w, http.StatusBadRequest, "Group DisplayId is required")
		return
	}

	monitor, err := monitors.FromReader(r.Body)
	if err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitorCreateResponse, serviceErr := services.MonitorService.CreateMonitor(r.Context(), teamAuth, groupId, monitor)

	if serviceErr != nil {
		util.RespondError(w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, monitorCreateResponse)
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitorId := r.PathValue("monitorId")

	if monitorId == "" {
		logging.Api.Trace().Msg("Monitor DisplayId is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "BaseMonitor DisplayId is required")
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	err := services.MonitorService.DeleteMonitor(r.Context(), teamAuth, monitorId)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Monitor deleted successfully")
}

func GetAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitorsList, err := services.MonitorService.GetAllMonitors(r.Context(), teamAuth)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorsList)
}

func GetMonitorByIdHandler(w http.ResponseWriter, r *http.Request) {
	monitorId := r.PathValue("monitorId")

	if monitorId == "" {
		logging.Api.Trace().Msg("Monitor DisplayId is required")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor DisplayId is required")
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitor, err := services.MonitorService.GetMonitorById(r.Context(), teamAuth, monitorId)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitor)
}

// UpdateMonitorHandler handles the update of an existing monitor.
// // It expects a JSON payload with the updated monitor config of appropriate type.
// TODO: Proper update mechanism, maybe custom payload so we can update only specific fields
func UpdateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitorId := r.PathValue("monitorId")

	if monitorId == "" {
		logging.Api.Trace().Msg("Monitor DisplayId is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor DisplayId is required")
		return
	}

	monitor, err := monitors.FromReader(r.Body)
	if err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	serviceErr := services.MonitorService.UpdateMonitor(r.Context(), teamAuth, monitor)
	if serviceErr != nil {
		util.RespondError(w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "monitor updated successfully")
}
