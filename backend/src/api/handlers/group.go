package handlers

import (
	"encoding/json"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
	"github.com/m-milek/leszmonitor/logging"
	"net/http"
)

type CreateMonitorGroupPayload struct {
	Name        string `json:"name"`        // Name of the monitor group
	Description string `json:"description"` // Description of the monitor group
}

func CreateMonitorGroupHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")

	var payload CreateMonitorGroupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to decode monitor group creation payload")
		util.RespondError(w, http.StatusBadRequest, err)
	}

	monitorGroup, err := services.GroupService.CreateMonitorGroup(r.Context(), teamId, payload.Name, payload.Description)

	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to create monitor group")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, monitorGroup)
}

func GetTeamMonitorGroups(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")

	if teamId == "" {
		logging.Api.Trace().Msg("Team ID is required to get monitor groups")
		util.RespondMessage(w, http.StatusBadRequest, "Team ID is required")
		return
	}

	monitorGroups, err := services.GroupService.GetTeamMonitorGroups(r.Context(), teamId)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to get monitor groups for team")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorGroups)
}

func GetTeamMonitorGroupById(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")
	groupId := r.PathValue("groupId")

	if groupId == "" {
		logging.Api.Trace().Msg("Monitor group ID is required to get monitor group")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor group ID is required")
		return
	}

	monitorGroup, err := services.GroupService.GetTeamMonitorGroupById(r.Context(), teamId, groupId)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to get monitor group by ID")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorGroup)
}

func DeleteMonitorGroupHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")
	groupId := r.PathValue("groupId")

	if groupId == "" {
		logging.Api.Trace().Msg("Monitor group ID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor group ID is required")
		return
	}

	err := services.GroupService.DeleteMonitorGroup(r.Context(), teamId, groupId)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to delete monitor group")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Monitor group deleted successfully")
}

func UpdateMonitorGroupHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("teamId")
	groupId := r.PathValue("groupId")

	if groupId == "" {
		logging.Api.Trace().Msg("Monitor group ID is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor group ID is required")
		return
	}

	var payload UpdateMonitorGroupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to decode monitor group update payload")
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	monitorGroup, err := services.GroupService.UpdateMonitorGroup(r.Context(), teamId, groupId, &payload)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to update monitor group")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorGroup)
}
