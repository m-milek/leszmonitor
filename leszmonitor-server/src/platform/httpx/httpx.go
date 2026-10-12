package httpx

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/log"
)

func RespondJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	logger := log.FromContext(ctx)

	w.Header().Set(constants.HTTPHeaderContentType, constants.HTTPContentTypeJSON)
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorBody struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Status int       `json:"status"`
	Error  ErrorBody `json:"error"`
}

func RespondMessage(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	logger := log.FromContext(ctx)
	w.Header().Set(constants.HTTPHeaderContentType, constants.HTTPContentTypeJSON)
	w.WriteHeader(statusCode)

	response := MessageResponse{Message: message}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func RespondError(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	logger := log.FromContext(ctx)
	logger.Error().Err(err).Msg("Responding with error")

	w.Header().Set(constants.HTTPHeaderContentType, constants.HTTPContentTypeJSON)
	w.WriteHeader(statusCode)

	message := err.Error()

	// Obfuscate internal server error messages
	if statusCode == http.StatusInternalServerError {
		message = "Internal server error"
	}

	response := ErrorResponse{
		Status: statusCode,
		Error:  ErrorBody{Message: message},
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.Error().Err(encodeErr).Msg("Failed to encode JSON error response")
	}
}

// DecodeJSONOrRespond decodes JSON from the request body into v, or writes a 400 response and returns false.
func DecodeJSONOrRespond(ctx context.Context, w http.ResponseWriter, r *http.Request, v any) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		RespondError(ctx, w, http.StatusBadRequest, err)
		return false
	}
	return true
}
