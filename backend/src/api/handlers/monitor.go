package handlers

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/api/util"
	"github.com/m-milek/leszmonitor/logger"
	"github.com/m-milek/leszmonitor/uptime/monitors"
	"net/http"
	"reflect"
)

type MonitorTypeExtractor struct {
	Type monitors.MonitorType `json:"type"`
}

var monitorTypeMap = map[string]reflect.Type{
	string(monitors.Http): reflect.TypeOf(monitors.HttpMonitor{}),
	string(monitors.Ping): reflect.TypeOf(monitors.PingMonitor{}),
}

func mapMonitorType(typeTag string) reflect.Type {
	monitorType := monitors.MonitorType(typeTag)
	if monitorType == "" {
		return nil
	}
	if monitorType, ok := monitorTypeMap[string(monitorType)]; ok {
		return monitorType
	}
	return nil
}

// AddMonitor handles the addition of a new monitor.
// It expects a JSON payload with the monitor config of appropriate type.
func AddMonitor(w http.ResponseWriter, r *http.Request) {
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
	monitorType := mapMonitorType(string(typeExtractor.Type))
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

	logger.Api.Debug().Any("monitor", monitorInstance).Msg("Parsed monitor configuration")
}
