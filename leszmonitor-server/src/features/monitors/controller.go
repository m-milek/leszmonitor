package monitors

import (
	"encoding/json"
	"io"
	"net/http"
	"slices"
	"strconv"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/httpx"
)

type MonitorAPIController struct {
	service IMonitorService
}

func NewMonitorAPIController(service IMonitorService) MonitorAPIController {
	return MonitorAPIController{
		service: service,
	}
}

const messageMonitorIDIsRequired = "Monitor ID is required"

const maxPushBodyBytes = 5000

// CreateMonitorHandler handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
func (c *MonitorAPIController) CreateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitor, err := decodeMonitorPayload(r)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if monitor.ProbeConfig == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "probeConfig is required")
		return
	}

	_, err = probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig)
	if err != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid probe config: "+err.Error())
		return
	}
	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	monitorCreateResponse, serviceErr := c.service.CreateMonitor(ctx, monitor)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusCreated, monitorCreateResponse)
}

func (c *MonitorAPIController) DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageMonitorIDIsRequired)
		return
	}

	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}
	err := c.service.DeleteMonitor(ctx, monitorID)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusOK, "Monitor deleted successfully")
}

func (c *MonitorAPIController) GetAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allMonitors, err := c.service.GetAllMonitors(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, allMonitors)
}

func (c *MonitorAPIController) GetMonitorByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageMonitorIDIsRequired)
		return
	}

	monitor, err := c.service.GetMonitorByID(ctx, monitorID)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, monitor)
}

// UpdateMonitorHandler handles the update of an existing monitor.
// TODO: Proper update mechanism, maybe custom payload so we can update only specific fields
func (c *MonitorAPIController) UpdateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageMonitorIDIsRequired)
		return
	}

	monitor, err := decodeMonitorPayload(r)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if monitor.ProbeConfig == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "probeConfig is required")
		return
	}

	_, err = probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig)
	if err != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	serviceErr := c.service.UpdateMonitor(ctx, monitor)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusOK, "monitor updated successfully")
}

// decodeMonitorPayload decodes the request body into Monitor, and probeConfig separately as string.
// FE sends probeConfig as JSON object, but we want to store it as string in the database, so we need to handle it separately.
func decodeMonitorPayload(r *http.Request) (Monitor, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return Monitor{}, err
	}

	var rawPayload map[string]json.RawMessage
	if err := json.Unmarshal(bodyBytes, &rawPayload); err != nil {
		return Monitor{}, err
	}

	probeConfigRaw := rawPayload["probeConfig"]
	delete(rawPayload, "probeConfig")

	payloadBytes, err := json.Marshal(rawPayload)
	if err != nil {
		return Monitor{}, err
	}

	var monitor Monitor
	if err := json.Unmarshal(payloadBytes, &monitor); err != nil {
		return Monitor{}, err
	}

	monitor.ProbeConfig = string(probeConfigRaw)
	return monitor, nil
}

type UpdateMonitorStatePayload struct {
	NewState string `json:"newState"`
}

func (c *MonitorAPIController) UpdateMonitorStateByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageMonitorIDIsRequired)
		return
	}

	monitorUUID, err := uuid.Parse(monitorID)
	if err != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid monitor ID format")
		return
	}

	payload := UpdateMonitorStatePayload{}
	ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload)
	if !ok {
		return
	}

	if payload.NewState == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "newState is required")
		return
	}

	_, ok = auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	serviceErr := c.service.UpdateMonitorStateByID(
		ctx,
		monitorUUID,
		MonitorRunState(payload.NewState),
	)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusOK, "Monitor state updated successfully")
}

func (c *MonitorAPIController) RunMonitorManuallyByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, messageMonitorIDIsRequired)
		return
	}

	monitorUUID, err := uuid.Parse(monitorID)
	if err != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid monitor ID format")
		return
	}

	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	serviceErr := c.service.RunMonitorManuallyByID(ctx, monitorUUID)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusAccepted, "Monitor run scheduled")
}

func (c *MonitorAPIController) ReceivePushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorUUID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid monitor ID format")
		return
	}

	status := r.URL.Query().Get("status")
	if !slices.Contains([]string{"up", "down", ""}, status) {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid status reported")
		return
	}

	var latency *int64
	if latencyStr := r.URL.Query().Get("latency"); latencyStr != "" {
		parsed, err := strconv.ParseInt(latencyStr, 10, 64)
		if err != nil || parsed < 0 {
			httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid latency reported")
			return
		}
		latency = &parsed
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxPushBodyBytes+1))
	if err != nil {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Failed to read body")
		return
	}
	if len(body) > maxPushBodyBytes {
		httpx.RespondMessage(ctx, w, http.StatusBadRequest, "Body too large")
		return
	}

	payload := probe.PushProbePayload{
		Status:  kind.MonitorStatus(status),
		Latency: latency,
		Body:    string(body),
	}

	svcErr := c.service.ReceivePush(ctx, monitorUUID, payload)
	if svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}

	httpx.RespondMessage(ctx, w, http.StatusAccepted, "")
}
