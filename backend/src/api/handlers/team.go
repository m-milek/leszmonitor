package handlers

import (
	"encoding/json"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
)

func TeamCreateHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.TeamCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	userClaims, ok := util.ExtractUserOrRespond(w, r)
	if !ok {
		return
	}

	teamCreateResponse, err := services.TeamService.CreateTeam(r.Context(), userClaims.Username, &payload)

	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, teamCreateResponse)
}

func TeamDeleteHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	err := services.TeamService.DeleteTeam(r.Context(), teamAuth)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "")
}

func TeamUpdateHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.TeamUpdatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	team, err := services.TeamService.UpdateTeam(r.Context(), teamAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, team)
}

func GetTeamHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	team, err := services.TeamService.GetTeamByID(r.Context(), teamAuth)

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
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.TeamAddMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.TeamService.AddUserToTeam(r.Context(), teamAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "User added to team successfully")
}

func TeamRemoveMemberHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.TeamRemoveMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.TeamService.RemoveUserFromTeam(r.Context(), teamAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "User removed from team successfully")
}

func TeamChangeMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	teamAuth, ok := util.GetTeamAuthOrRespond(w, r)
	if !ok {
		return
	}

	var payload services.TeamChangeMemberRolePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.TeamService.ChangeMemberRole(r.Context(), teamAuth, payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "User role updated successfully")
}
