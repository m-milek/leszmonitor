package handlers

import (
	"encoding/json"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
	"github.com/m-milek/leszmonitor/logging"
	"net/http"
)

func CreateMonitorGroupHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.CreateMonitorGroupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to decode monitor group creation payload")
		util.RespondError(w, http.StatusBadRequest, err)
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitorGroup, err := services.GroupService.CreateMonitorGroup(r.Context(), teamAuth, payload)

	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to create monitor group")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, monitorGroup)
}

func GetTeamMonitorGroupsHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitorGroups, err := services.GroupService.GetTeamMonitorGroups(r.Context(), teamAuth)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to get monitor groups for team")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorGroups)
}

func GetTeamMonitorGroupById(w http.ResponseWriter, r *http.Request) {
	groupId := r.PathValue("groupId")

	if groupId == "" {
		logging.Api.Trace().Msg("Monitor group ID is required to get monitor group")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor group ID is required")
		return
	}

	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	monitorGroup, err := services.GroupService.GetTeamMonitorGroupById(r.Context(), teamAuth, groupId)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to get monitor group by ID")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorGroup)
}

func DeleteMonitorGroupHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	groupId := r.PathValue("groupId")

	if groupId == "" {
		logging.Api.Trace().Msg("Monitor group ID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor group ID is required")
		return
	}

	err := services.GroupService.DeleteMonitorGroup(r.Context(), teamAuth, groupId)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to delete monitor group")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Monitor group deleted successfully")
}

func UpdateMonitorGroupHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	groupId := r.PathValue("groupId")

	if groupId == "" {
		logging.Api.Trace().Msg("Monitor group ID is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Monitor group ID is required")
		return
	}

	var payload services.UpdateMonitorGroupPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to decode monitor group update payload")
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	monitorGroup, err := services.GroupService.UpdateMonitorGroup(r.Context(), teamAuth, groupId, &payload)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to update monitor group")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, monitorGroup)
}
