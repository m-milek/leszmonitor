package util

import (
	"encoding/json"
	"github.com/m-milek/leszmonitor/api/middleware"
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

type simpleResponse struct {
	Message string `json:"message"`
}

type errorResponse struct {
	Message string `json:"message"`
}

func RespondMessage(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := simpleResponse{Message: message}

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
		"error":  errorResponse{Message: message},
	}

	if encodeErr := json.NewEncoder(w).Encode(response); encodeErr != nil {
		logging.Api.Error().Err(encodeErr).Msg("Failed to encode JSON error response")
	}
}

func GetTeamAuthOrRespond(w http.ResponseWriter, r *http.Request) (*middleware.TeamAuth, bool) {
	teamAuth, err := middleware.TeamAuthFromRequest(r)
	if err != nil {
		logging.Api.Warn().Err(err).Msg("Failed to authenticate")
		RespondError(w, http.StatusUnauthorized, err)
		return nil, false
	}
	return teamAuth, true
}

// ExtractUserOrRespond returns the user from context or writes a 401 response and returns nil, false.
func ExtractUserOrRespond(w http.ResponseWriter, r *http.Request) (*middleware.UserClaims, bool) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return nil, false
	}
	return user, true
}

// DecodeJSONOrRespond decodes JSON from the request body into v, or writes a 400 response and returns false.
func DecodeJSONOrRespond(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		RespondError(w, http.StatusBadRequest, err)
		return false
	}
	return true
}
