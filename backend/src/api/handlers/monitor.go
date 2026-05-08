package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/services"
)

// CreateMonitorHandler handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
func CreateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitor, err := decodeMonitorPayload(r)
	if err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if monitor.ProbeConfig == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "probeConfig is required")
		return
	}

	_, err = monitors.ProbeFromJSON(monitor.ProbeConfig, monitor.Type)
	if err != nil {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid probe config: "+err.Error())
		return
	}

	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	monitorCreateResponse, serviceErr := services.MonitorService.CreateMonitor(ctx, projectAuth, monitor)
	if serviceErr != nil {
		util.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusCreated, monitorCreateResponse)
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Monitor slug is required")
		return
	}

	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	err := services.MonitorService.DeleteMonitor(ctx, projectAuth, monitorID)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Monitor deleted successfully")
}

func GetAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	monitorsList, err := services.MonitorService.GetMonitorsByProjectID(ctx, projectAuth)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, monitorsList)
}

func GetMonitorByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Monitor slug is required")
		return
	}

	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	monitor, err := services.MonitorService.GetMonitorByID(ctx, projectAuth, monitorID)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, monitor)
}

// UpdateMonitorHandler handles the update of an existing monitor.
// TODO: Proper update mechanism, maybe custom payload so we can update only specific fields
func UpdateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Monitor slug is required")
		return
	}

	monitor, err := decodeMonitorPayload(r)
	if err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if monitor.ProbeConfig == "" {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "probeConfig is required")
		return
	}

	_, err = monitors.ProbeFromJSON(monitor.ProbeConfig, monitor.Type)
	if err != nil {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	serviceErr := services.MonitorService.UpdateMonitor(ctx, projectAuth, monitor)
	if serviceErr != nil {
		util.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "monitor updated successfully")
}

// decodeMonitorPayload decodes the request body into monitors.Monitor, and probeConfig separately as string.
// FE sends probeConfig as JSON object, but we want to store it as string in the database, so we need to handle it separately.
func decodeMonitorPayload(r *http.Request) (monitors.Monitor, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return monitors.Monitor{}, err
	}

	var rawPayload map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &rawPayload); err != nil {
		return monitors.Monitor{}, err
	}

	probeConfigRaw := rawPayload["probeConfig"]
	delete(rawPayload, "probeConfig")

	payloadBytes, err := json.Marshal(rawPayload)
	if err != nil {
		return monitors.Monitor{}, err
	}

	var monitor monitors.Monitor
	if err := json.Unmarshal(payloadBytes, &monitor); err != nil {
		return monitors.Monitor{}, err
	}

	monitor.ProbeConfig = string(probeConfigRaw)
	return monitor, nil
}
