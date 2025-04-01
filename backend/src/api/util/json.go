package util

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/logger"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Api.Error().Err(err).Msg("Failed to encode JSON response")
	}
}
