package handlers

import (
	"encoding/json"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/api/services"
	"net/http"
)

func TeamCreateHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload services.TeamCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	teamCreateResponse, err := services.TeamService.CreateTeam(r.Context(), &payload, user.Username)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, teamCreateResponse)
}

func TeamDeleteHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	if teamId == "" {
		util.RespondMessage(w, http.StatusBadRequest, "Team ID is required")
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := services.TeamService.DeleteTeam(r.Context(), teamId, requestingUser.Username)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "")
}

func TeamUpdateHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	if teamId == "" {
		util.RespondMessage(w, http.StatusBadRequest, "Team ID is required")
		return
	}

	var payload services.TeamUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	team, err := services.TeamService.UpdateTeam(r.Context(), teamId, &payload, requestingUser.Username)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, team)
}

func GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	if teamId == "" {
		util.RespondMessage(w, http.StatusBadRequest, "Team ID is required")
		return
	}

	team, err := services.TeamService.GetTeamById(r.Context(), teamId)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, team)
}

func GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := services.TeamService.GetAllTeams(r.Context())

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, teams)
}

func TeamAddMemberHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	var payload services.TeamAddMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := services.TeamService.AddUserToTeam(r.Context(), teamId, &payload, requestingUser.Username)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "User added to team successfully")
}

func TeamRemoveMemberHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	var payload services.TeamRemoveMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := services.TeamService.RemoveUserFromTeam(r.Context(), teamId, &payload, requestingUser.Username)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "User removed from team successfully")
}

func TeamChangeMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	var payload services.TeamChangeMemberRolePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	err := services.TeamService.ChangeMemberRole(r.Context(), teamId, payload, requestingUser.Username)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "User role updated successfully")
}
