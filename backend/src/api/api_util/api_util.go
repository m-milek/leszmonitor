package util

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/logging"
	"net/http"
)

func RespondJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logging.Api.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

type SimpleResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

func RespondMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := SimpleResponse{Message: message}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logging.Api.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func RespondError(w http.ResponseWriter, statusCode int, err error) {
	logging.Api.Error().Err(err).Msg("Responding with error")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	message := err.Error()

	// Obfuscate internal server error messages
	if statusCode == http.StatusInternalServerError {
		message = "Internal server error"
	}

	response := map[string]any{
		"status": statusCode,
		"error":  ErrorResponse{Message: message},
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logging.Api.Error().Err(encodeErr).Msg("Failed to encode JSON error response")
	}
}
