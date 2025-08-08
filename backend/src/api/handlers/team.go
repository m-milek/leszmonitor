package handlers

import (
	"encoding/json"
	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/common"
	"github.com/m-milek/leszmonitor/db"
	"net/http"
)

type TeamCreatePayload struct {
	Name        string `json:"name"`        // The name of the team
	Description string `json:"description"` // A brief description of the team
}

func TeamCreateHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload TeamCreatePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	team := common.NewTeam(payload.Name, payload.Description, user.Username)

	_, err := db.AddTeam(team)

	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to create team: "+err.Error())
		return
	}

	util.RespondJSON(w, http.StatusCreated, map[string]string{
		"teamId": team.Id,
	})
}

func TeamDeleteHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	if teamId == "" {
		util.RespondMessage(w, http.StatusBadRequest, "Team ID is required")
		return
	}

	team, err := db.GetTeamById(teamId)

	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve team: "+err.Error())
		return
	}

	if team == nil {
		util.RespondMessage(w, http.StatusNotFound, "Team not found")
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !team.IsAdmin(requestingUser.Username) {
		util.RespondMessage(w, http.StatusForbidden, "You are not authorized to delete this team")
		return
	}

	_, err = db.DeleteTeam(teamId)
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to delete team: "+err.Error())
		return
	}

	util.RespondMessage(w, http.StatusOK, "Team deleted successfully")
}

func TeamUpdateHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")

	if teamId == "" {
		util.RespondMessage(w, http.StatusBadRequest, "Team ID is required")
		return
	}

	var payload common.Team
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	team, err := db.GetTeamById(teamId)

	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve team: "+err.Error())
		return
	}

	if team == nil {
		util.RespondMessage(w, http.StatusNotFound, "Team not found")
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !team.IsAdmin(requestingUser.Username) {
		util.RespondMessage(w, http.StatusForbidden, "You are not authorized to update this team")
		return
	}

	team.Name = payload.Name
	team.Description = payload.Description

	_, err = db.UpdateTeam(team)
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to update team: "+err.Error())
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

	team, err := db.GetTeamById(teamId)

	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve team: "+err.Error())
		return
	}

	if team == nil {
		util.RespondMessage(w, http.StatusNotFound, "Team not found")
		return
	}

	util.RespondJSON(w, http.StatusOK, team)
}

func GetAllTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := db.GetAllTeams()
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve teams: "+err.Error())
		return
	}

	util.RespondJSON(w, http.StatusOK, teams)
}

func TeamAddUserHandler(w http.ResponseWriter, r *http.Request) {
	teamId := r.PathValue("id")
	username := r.URL.Query().Get("username")
	roleStr := r.URL.Query().Get("role")

	if teamId == "" || username == "" || roleStr == "" {
		util.RespondMessage(w, http.StatusBadRequest, "Team ID, Username and Role are required")
		return
	}

	team, err := db.GetTeamById(teamId)
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve team: "+err.Error())
		return
	}

	requestingUser, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		util.RespondMessage(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if !team.IsAdmin(requestingUser.Username) {
		util.RespondMessage(w, http.StatusForbidden, "You are not authorized to add users to this team")
		return
	}

	// Convert to TeamRole and validate
	role := common.TeamRole(roleStr)
	if err := role.Validate(); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, err.Error())
		return
	}

	_, err = db.GetUser(username)
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to retrieve user: "+err.Error())
		return
	}

	err = team.AddMember(username, role)
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to add user to team: "+err.Error())
		return
	}

	_, err = db.UpdateTeam(team)
	if err != nil {
		util.RespondMessage(w, http.StatusInternalServerError, "Failed to add user to team: "+err.Error())
		return
	}

	util.RespondMessage(w, http.StatusOK, "User added to team successfully")
}
