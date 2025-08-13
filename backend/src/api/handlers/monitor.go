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

	monitor, err := monitors.FromReader(r.Body)
	if err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	monitorCreateResponse, serviceErr := services.MonitorService.CreateMonitor(r.Context(), monitor)

	if serviceErr != nil {
		util.RespondError(w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, monitorCreateResponse)
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		logging.Api.Trace().Msg("Monitor ID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "BaseMonitor ID is required")
		return
	}

	err := services.MonitorService.DeleteMonitor(r.Context(), id)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Monitor deleted successfully")
}

func GetAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	monitorsList, err := services.MonitorService.GetAllMonitors(r.Context())
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorsList)
}

func GetMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		logging.Api.Trace().Msg("Monitor ID is required")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor ID is required")
		return
	}

	monitor, err := services.MonitorService.GetMonitorById(r.Context(), id)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitor)
}

// UpdateMonitorHandler handles the update of an existing monitor.
// // It expects a JSON payload with the updated monitor config of appropriate type.
func UpdateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		logging.Api.Trace().Msg("Monitor ID is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor ID is required")
		return
	}

	monitor, err := monitors.FromReader(r.Body)
	if err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	if monitor.GetId() != id {
		logging.Api.Trace().Msgf("Monitor ID mismatch: expected %s, got %s", id, monitor.GetId())
		util.RespondMessage(w, http.StatusBadRequest, "Monitor ID mismatch")
		return
	}

	serviceErr := services.MonitorService.UpdateMonitor(r.Context(), monitor)
	if serviceErr != nil {
		util.RespondError(w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, "monitor updated successfully")
}
