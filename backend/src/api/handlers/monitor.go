package handlers

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
	"net/http"
)

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

	logger.Api.Debug().RawJSON("raw_data", rawData).Msg("Received monitor configuration")

	var monitorTypeExtractor monitors.MonitorTypeExtractor
	if err := json.Unmarshal(rawData, &monitorTypeExtractor); err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to unmarshal monitor type")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor type")
		return
	}

	// Map the monitor type to the appropriate config type
	monitor := monitors.MapMonitorType(monitorTypeExtractor.Type)
	if monitor == nil {
		logger.Api.Trace().Msgf("Unknown monitor type: %s", monitorTypeExtractor.Type)
		util.RespondMessage(w, http.StatusBadRequest, "Unknown monitor type: "+string(monitorTypeExtractor.Type))
		return
	}

	// unmarshal the raw data into a monitor instance
	if err := json.Unmarshal(rawData, &monitor); err != nil {
		logger.Api.Trace().Err(err).Msg("Failed to unmarshal monitor data")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor data")
		return
	}

	monitor.GenerateId()

	if err := monitor.Validate(); err != nil {
		logger.Api.Trace().Err(err).Msg("BaseMonitor validation failed")
		util.RespondMessage(w, http.StatusBadRequest, "Invalid monitor configuration: "+err.Error())
		return
	}

	logger.Api.Debug().Any("monitor", monitor).Msg("Parsed monitor configuration")

	_, err := db.AddMonitor(monitor)
	if err != nil {
		logger.Api.Error().Err(err).Msg("Failed to add monitor to database")
		util.RespondMessage(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	util.RespondJSON(w, http.StatusCreated, AddMonitorResponse{
		MonitorId: monitor.GetId(),
	})
}

func DeleteMonitorHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		logger.Api.Trace().Msg("BaseMonitor ID is required for deletion")
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
		msg := "BaseMonitor not found or already deleted"
		logger.Api.Warn().Str("monitor_id", id).Msg(msg)
		util.RespondMessage(w, http.StatusNotFound, msg)
		return
	}

	util.RespondMessage(w, http.StatusOK, "BaseMonitor deleted successfully")
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
		logger.Api.Trace().Msg("BaseMonitor ID is required")
		util.RespondMessage(writer, http.StatusBadRequest, "BaseMonitor ID is required")
		return
	}

	monitor, err := db.GetMonitorById(id)
	if err != nil {
		logger.Api.Error().Err(err).Str("monitor_id", id).Msg("Failed to retrieve monitor from database")
		util.RespondMessage(writer, http.StatusInternalServerError, "Failed to retrieve monitor")
		return
	}

	if monitor == nil {
		logger.Api.Warn().Str("monitor_id", id).Msg("BaseMonitor not found")
		util.RespondMessage(writer, http.StatusNotFound, "BaseMonitor not found")
		return
	}

	util.RespondJSON(writer, http.StatusOK, monitor)
}
