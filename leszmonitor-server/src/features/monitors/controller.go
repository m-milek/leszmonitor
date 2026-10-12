package monitors

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
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

const messageMonitorIDIsRequired = "monitor ID is required"

const maxPushBodyBytes = 5000

// CreateMonitorHandler handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
//
//	@Summary	Create a monitor
//	@Tags		monitors
//	@Accept		json
//	@Produce	json
//	@Param		body	body		Monitor	true	"Monitor definition; probeConfig is sent as a JSON object"
//	@Success	201		{object}	Monitor
//	@Failure	400		{object}	httpx.ErrorResponse
//	@Failure	401		{object}	httpx.ErrorResponse
//	@Failure	403		{object}	httpx.ErrorResponse
//	@Failure	500		{object}	httpx.ErrorResponse
//	@Security	BearerAuth
//	@Router		/monitors [post]
func (c *MonitorAPIController) CreateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitor, err := decodeMonitorPayload(r)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	if monitor.ProbeConfig == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("probeConfig is required"))
		return
	}

	_, err = probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("invalid probe config: %w", err))
		return
	}
	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	createdMonitor, serviceErr := c.service.CreateMonitor(ctx, monitor)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusCreated, createdMonitor)
}

// @Summary	Delete a monitor
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path	string	true	"Monitor ID"
// @Success	204
// @Failure	400	{object}	httpx.ErrorResponse
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	404	{object}	httpx.ErrorResponse
// @Failure	409	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId} [delete]
func (c *MonitorAPIController) DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New(messageMonitorIDIsRequired))
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

	w.WriteHeader(http.StatusNoContent)
}

// @Summary	List monitors
// @Tags		monitors
// @Produce	json
// @Success	200	{array}		Monitor
// @Failure	401	{object}	httpx.ErrorResponse
// @Failure	403	{object}	httpx.ErrorResponse
// @Failure	500	{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors [get]
func (c *MonitorAPIController) GetAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	allMonitors, err := c.service.GetAllMonitors(ctx)
	if err != nil {
		httpx.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, allMonitors)
}

// @Summary	Get a monitor
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path		string	true	"Monitor ID"
// @Success	200			{object}	Monitor
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId} [get]
func (c *MonitorAPIController) GetMonitorByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New(messageMonitorIDIsRequired))
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
//
//	@Summary	Update a monitor
//	@Tags		monitors
//	@Accept		json
//	@Produce	json
//	@Param		monitorId	path		string	true	"Monitor ID"
//	@Param		body		body		Monitor	true	"Monitor definition; probeConfig is sent as a JSON object"
//	@Success	200			{object}	Monitor
//	@Failure	400			{object}	httpx.ErrorResponse
//	@Failure	401			{object}	httpx.ErrorResponse
//	@Failure	403			{object}	httpx.ErrorResponse
//	@Failure	404			{object}	httpx.ErrorResponse
//	@Failure	409			{object}	httpx.ErrorResponse
//	@Failure	500			{object}	httpx.ErrorResponse
//	@Security	BearerAuth
//	@Router		/monitors/{monitorId} [patch]
func (c *MonitorAPIController) UpdateMonitorHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorUUID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
		return
	}

	monitor, err := decodeMonitorPayload(r)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}
	monitor.ID = monitorUUID

	if monitor.ProbeConfig == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("probeConfig is required"))
		return
	}

	_, err = probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, fmt.Errorf("invalid probe config: %w", err))
		return
	}

	_, ok := auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	updatedMonitor, serviceErr := c.service.UpdateMonitor(ctx, monitor)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, updatedMonitor)
}

type UpdateMonitorStatePayload struct {
	NewState string `json:"newState"`
}

// @Summary	Pause or resume a monitor
// @Tags		monitors
// @Accept		json
// @Produce	json
// @Param		monitorId	path		string						true	"Monitor ID"
// @Param		body		body		UpdateMonitorStatePayload	true	"New run state (active or paused)"
// @Success	200			{object}	Monitor
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	409			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/state [patch]
func (c *MonitorAPIController) UpdateMonitorStateByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New(messageMonitorIDIsRequired))
		return
	}

	monitorUUID, err := uuid.Parse(monitorID)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
		return
	}

	payload := UpdateMonitorStatePayload{}
	ok := httpx.DecodeJSONOrRespond(ctx, w, r, &payload)
	if !ok {
		return
	}

	if payload.NewState == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("newState is required"))
		return
	}

	_, ok = auth.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	updatedMonitor, serviceErr := c.service.UpdateMonitorStateByID(
		ctx,
		monitorUUID,
		MonitorRunState(payload.NewState),
	)
	if serviceErr != nil {
		httpx.RespondError(ctx, w, serviceErr.Code, serviceErr.Err)
		return
	}

	httpx.RespondJSON(ctx, w, http.StatusOK, updatedMonitor)
}

// @Summary	Run a monitor now
// @Tags		monitors
// @Produce	json
// @Param		monitorId	path		string	true	"Monitor ID"
// @Success	202			{object}	httpx.MessageResponse
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/run [post]
func (c *MonitorAPIController) RunMonitorManuallyByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	monitorID := r.PathValue("monitorId")
	if monitorID == "" {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New(messageMonitorIDIsRequired))
		return
	}

	monitorUUID, err := uuid.Parse(monitorID)
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
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

	httpx.RespondMessage(ctx, w, http.StatusAccepted, "monitor run scheduled")
}

// @Summary		Report a push monitor heartbeat
// @Description	POST may carry a raw body of up to 5000 bytes, stored with the result.
// @Tags			push
// @Produce		json
// @Param			monitorId	path	string	true	"Monitor ID"
// @Param			status		query	string	false	"Reported status"					Enums(up, down)
// @Param			latency		query	int		false	"Reported latency in milliseconds"	minimum(0)
// @Success		202
// @Failure		400	{object}	httpx.ErrorResponse
// @Failure		404	{object}	httpx.ErrorResponse
// @Failure		409	{object}	httpx.ErrorResponse
// @Failure		500	{object}	httpx.ErrorResponse
// @Router			/push/{monitorId} [get]
// @Router			/push/{monitorId} [post]
func (c *MonitorAPIController) ReceivePushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorUUID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
		return
	}

	status := r.URL.Query().Get("status")
	if !slices.Contains([]string{"up", "down", ""}, status) {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid status reported"))
		return
	}

	var latency *int64
	if latencyStr := r.URL.Query().Get("latency"); latencyStr != "" {
		parsed, err := strconv.ParseInt(latencyStr, 10, 64)
		if err != nil || parsed < 0 {
			httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid latency reported"))
			return
		}
		latency = &parsed
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, maxPushBodyBytes+1))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("failed to read body"))
		return
	}
	if len(body) > maxPushBodyBytes {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("body too large"))
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

	w.WriteHeader(http.StatusAccepted)
}

// @Summary	Export a monitor as YAML config
// @Tags		monitors
// @Produce	plain
// @Param		monitorId	path		string	true	"Monitor ID"
// @Success	200			{string}	string	"Monitor config in YAML"
// @Failure	400			{object}	httpx.ErrorResponse
// @Failure	401			{object}	httpx.ErrorResponse
// @Failure	403			{object}	httpx.ErrorResponse
// @Failure	404			{object}	httpx.ErrorResponse
// @Failure	500			{object}	httpx.ErrorResponse
// @Security	BearerAuth
// @Router		/monitors/{monitorId}/export [get]
func (c *MonitorAPIController) ExportToYAMLHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	monitorID, err := uuid.Parse(r.PathValue("monitorId"))
	if err != nil {
		httpx.RespondError(ctx, w, http.StatusBadRequest, errors.New("invalid monitor ID format"))
		return
	}

	if svcErr := c.service.ExportToYAML(ctx, monitorID, w); svcErr != nil {
		httpx.RespondError(ctx, w, svcErr.Code, svcErr.Err)
		return
	}
}

// decodeMonitorPayload decodes the request body into Monitor, and probeConfig separately as string.
// FE sends probeConfig as JSON object, but we want to store it as string in the database, so we need to handle it separately.
func decodeMonitorPayload(r *http.Request) (Monitor, error) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return Monitor{}, err
	}

	return decodeMonitorJSON(bodyBytes)
}

func decodeMonitorJSON(data []byte) (Monitor, error) {
	var rawPayload map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawPayload); err != nil {
		return Monitor{}, err
	}

	probeConfigRaw := rawPayload["probeConfig"]
	delete(rawPayload, "probeConfig")

	payloadBytes, err := json.Marshal(rawPayload)
	if err != nil {
		return Monitor{}, err
	}

	var monitor Monitor
	dec := json.NewDecoder(bytes.NewReader(payloadBytes))
	dec.DisallowUnknownFields()
	if err := dec.Decode(&monitor); err != nil {
		return Monitor{}, err
	}

	monitor.ProbeConfig = string(probeConfigRaw)
	return monitor, nil
}
