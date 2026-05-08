package util

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/log"
)

func RespondJSON(ctx context.Context, w http.ResponseWriter, statusCode int, data any) {
	logger := log.FromContext(ctx)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

type simpleResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func RespondMessage(ctx context.Context, w http.ResponseWriter, statusCode int, message string) {
	logger := log.FromContext(ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := simpleResponse{Message: message}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Error().Err(err).Msg("Failed to encode JSON response")
	}
}

func RespondError(ctx context.Context, w http.ResponseWriter, statusCode int, err error) {
	logger := log.FromContext(ctx)
	logger.Error().Err(err).Msg("Responding with error")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	message := err.Error()

	// Obfuscate internal server error messages
	if statusCode == http.StatusInternalServerError {
		message = "Internal server error"
	}

	response := map[string]any{
		"status": statusCode,
		"error":  errorResponse{Message: message},
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logger.Error().Err(encodeErr).Msg("Failed to encode JSON error response")
	}
}

// GetProjectAuthOrRespond extracts project auth from the request or writes a 401 and returns nil, false.
func GetProjectAuthOrRespond(ctx context.Context, w http.ResponseWriter, r *http.Request, authSource middleware.AuthSourceKind) (*middleware.ProjectAuth, bool) {
	projectAuth, err := middleware.ProjectAuthFromRequest(r, authSource)
	if err != nil {
		RespondError(ctx, w, http.StatusUnauthorized, err)
		return nil, false
	}
	return projectAuth, true
}

// ExtractUserOrRespond returns the user from context or writes a 401 response and returns nil, false.
func ExtractUserOrRespond(ctx context.Context, w http.ResponseWriter, r *http.Request) (*auth.UserClaims, bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondMessage(ctx, w, http.StatusUnauthorized, "Unauthorized")
		return nil, false
	}
	return user, true
}

// DecodeJSONOrRespond decodes JSON from the request body into v, or writes a 400 response and returns false.
func DecodeJSONOrRespond(ctx context.Context, w http.ResponseWriter, r *http.Request, v interface{}) bool {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		RespondError(ctx, w, http.StatusBadRequest, err)
		return false
	}
	return true
}
