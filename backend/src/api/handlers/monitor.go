package handlers

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
	"net/http"
	"reflect"
)

type MonitorTypeExtractor struct {
	Type monitors.MonitorType `json:"type"`
}

type AddMonitorResponse struct {
	MonitorId string `json:"monitor_id"`
}

// AddMonitorHandler handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
func AddMonitorHandler(w http.ResponseWriter, r *http.Request) {
	var rawData json.RawMessage
	if err := json.NewDecoder(r.Body).Decode(&rawData); err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to decode request body")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	var typeExtractor MonitorTypeExtractor
	if err := json.Unmarshal(rawData, &typeExtractor); err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to extract monitor type from request payload")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor type in request payload")
		return
	}

	// Map the monitor type to its concrete type
	monitorType := monitors.MapMonitorType(string(typeExtractor.Type))
	if monitorType == nil {
		logger.Api.Trace().Str("type", string(typeExtractor.Type)).Msg("Unknown monitor type")
		util.RespondMessage(w, http.StatusBadRequest, "Unknown monitor type")
		return
	}

	monitorInstance := reflect.New(monitorType).Interface()

	// Check if it implements IMonitor
	monitor, ok := monitorInstance.(monitors.IMonitor)
	if !ok {
		logger.Api.Trace().Str("type", string(typeExtractor.Type)).Msg("Monitor type does not implement IMonitor interface")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor type")
		return
	}

	// Unmarshal the monitor configuration
	err := monitors.UnmarshalMonitor(rawData, monitor)
	if err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to parse monitor configuration")
		util.RespondMessage(w, http.StatusBadRequest, "Failed to parse monitor configuration")
		return
	}

	logger.Api.Debug().Any("monitor", monitor).Msg("Parsed monitor configuration")

	id, err := db.AddMonitor(monitor)
	if err != nil {
		logger.Api.Error().Err(err).Msg("Failed to add monitor to database")
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to add monitor")
		return
	}
	logger.Api.Debug().Str("monitor_id", id).Msg("Monitor added successfully")

	util.RespondJSON(w, http.StatusCreated, AddMonitorResponse{
		MonitorId: monitor.GetId(),
	})
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		logger.Api.Trace().Msg("Monitor ID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor ID is required")
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
