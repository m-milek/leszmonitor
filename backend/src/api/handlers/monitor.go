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
	groupID := r.PathValue("groupId")
	if groupID == "" {
		logging.Api.Trace().Msg("Group DisplayID is required for creating a monitor")
		util.RespondMessage(w, http.StatusBadRequest, "Group DisplayID is required")
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

	monitorCreateResponse, serviceErr := services.MonitorService.CreateMonitor(r.Context(), teamAuth, groupID, monitor)

	if serviceErr != nil {
		util.RespondError(w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, monitorCreateResponse)
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	monitorID := r.PathValue("monitorId")

	if monitorID == "" {
		logging.Api.Trace().Msg("Monitor DisplayID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "BaseMonitor DisplayID is required")
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	err := services.MonitorService.DeleteMonitor(r.Context(), teamAuth, monitorID)

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

func GetMonitorByIDHandler(w http.ResponseWriter, r *http.Request) {
	monitorID := r.PathValue("monitorId")

	if monitorID == "" {
		logging.Api.Trace().Msg("Monitor DisplayID is required")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor DisplayID is required")
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitor, err := services.MonitorService.GetMonitorByID(r.Context(), teamAuth, monitorID)
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
	monitorID := r.PathValue("monitorId")

	if monitorID == "" {
		logging.Api.Trace().Msg("Monitor DisplayID is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor DisplayID is required")
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
