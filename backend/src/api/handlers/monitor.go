package handlers

import (
	"errors"
	"github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitor"
	"net/http"
)

type monitorIdResponse struct {
	MonitorId string `json:"monitorId"`
}

// AddMonitorHandler handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
func AddMonitorHandler(w http.ResponseWriter, r *http.Request) {

	monitor, err := monitors.FromReader(r.Body)
	if err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	monitor.GenerateId()

	if err := monitor.Validate(); err != nil {
		logger.Api.Trace().Err(err).Msg("Monitor validation failed")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor configuration: "+err.Error())
		return
	}

	logger.Api.Debug().Any("monitor", monitor).Msg("Parsed monitor configuration")

	_, err = db.AddMonitor(monitor)
	if err != nil {
		logger.Api.Error().Err(err).Msg("Failed to add monitor to database")
		util.RespondMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	util.RespondJSON(w, http.StatusCreated, monitorIdResponse{
		MonitorId: monitor.GetId(),
	})
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		logger.Api.Trace().Msg("Monitor ID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "BaseMonitor ID is required")
		return
	}

	wasDeleted, err := db.DeleteMonitor(id)

	if err != nil {
		msg := "Failed to delete monitor"
		logger.Api.Error().Err(err).Str("monitor_id", id).Msg(msg)
		util.RespondMessage(w, http.StatusInternalServerError, msg)
		return
	}

	if !wasDeleted {
		msg := "Monitor not found or already deleted"
		logger.Api.Warn().Str("monitor_id", id).Msg(msg)
		util.RespondMessage(w, http.StatusNotFound, msg)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Monitor deleted successfully")
}

func GetAllMonitorsHandler(w http.ResponseWriter, r *http.Request) {
	monitorsList, err := db.GetAllMonitors()
	if err != nil {
		logger.Api.Error().Err(err).Msg("Failed to retrieve monitors from database")
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve monitors")
		return
	}

	logger.Api.Debug().Int("count", len(monitorsList)).Msg("Retrieved all monitors")

	util.RespondJSON(w, http.StatusOK, monitorsList)
}

func GetMonitorHandler(writer http.ResponseWriter, request *http.Request) {
	id := request.PathValue("id")

	if id == "" {
		logger.Api.Trace().Msg("Monitor ID is required")
		util.RespondMessage(writer, http.StatusBadRequest, "Monitor ID is required")
		return
	}

	monitor, err := db.GetMonitorById(id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			msg := "Monitor not found"
			logger.Api.Warn().Str("monitor_id", id).Msg(msg)
			util.RespondMessage(writer, http.StatusNotFound, msg)
			return
		}
		logger.Api.Error().Err(err).Str("monitor_id", id).Msg("Failed to retrieve monitor from database")
		util.RespondMessage(writer, http.StatusInternalServerError, "Failed to retrieve monitor")
		return
	}

	if monitor == nil {
		logger.Api.Warn().Str("monitor_id", id).Msg("Monitor not found")
		util.RespondMessage(writer, http.StatusNotFound, "Monitor not found")
		return
	}

	util.RespondJSON(writer, http.StatusOK, monitor)
}

type editMonitorResponse struct {
	MonitorId  string `json:"monitorId"`
	WasUpdated bool   `json:"wasUpdated"`
}

// EditMonitorHandler handles the update of an existing monitor.
// // It expects a JSON payload with the updated monitor config of appropriate type.
func EditMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if id == "" {
		logger.Api.Trace().Msg("Monitor ID is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor ID is required")
		return
	}

	monitor, err := monitors.FromReader(r.Body)
	if err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor config: "+err.Error())
		return
	}

	if monitor.GetId() != id {
		logger.Api.Trace().Msgf("Monitor ID mismatch: expected %s, got %s", id, monitor.GetId())
		util.RespondMessage(w, http.StatusBadRequest, "Monitor ID mismatch")
		return
	}

	if err := monitor.Validate(); err != nil {
		logger.Api.Trace().Err(err).Msg("Monitor validation failed")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor configuration for update: "+err.Error())
		return
	}

	wasUpdated, err := db.UpdateMonitor(monitor)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			msg := "Monitor not found"
			logger.Api.Trace().Str("monitor_id", id).Msg(msg)
			util.RespondMessage(w, http.StatusNotFound, msg)
			return
		}
		logger.Api.Error().Err(err).Str("monitor_id", id).Msg("Failed to update monitor in database")
		util.RespondMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if !wasUpdated {
		msg := "Monitor not found or no changes made"
		logger.Api.Warn().Str("monitor_id", id).Msg(msg)
		util.RespondMessage(w, http.StatusNotFound, msg)
		return
	}

	util.RespondJSON(w, http.StatusOK, editMonitorResponse{
		MonitorId:  id,
		WasUpdated: wasUpdated,
	})
}
