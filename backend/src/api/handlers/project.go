package handlers

import (
	"encoding/json"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/services"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.CreateProjectPayload
	if !util.DecodeJSONOrRespond(w, r, &payload) {
		return
	}

	userClaims, ok := util.ExtractUserOrRespond(w, r)
	if !ok {
		return
	}

	project, err := services.ProjectService.CreateProject(r.Context(), userClaims.Username, payload)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to create project")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, project)
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	userClaims, ok := util.ExtractUserOrRespond(w, r)
	if !ok {
		return
	}

	projects, err := services.ProjectService.GetProjectsForUser(r.Context(), userClaims.Username)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to get projects")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, projects)
}

func GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	project, err := services.ProjectService.GetProjectByID(r.Context(), projectAuth)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, project)
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	err := services.ProjectService.DeleteProject(r.Context(), projectAuth)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Project deleted successfully")
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	var payload services.UpdateProjectPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	project, err := services.ProjectService.UpdateProject(r.Context(), projectAuth, &payload)
	if err != nil {
		log.Api.Error().Err(err).Msg("Failed to update project")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, project)
}

func AddProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	var payload services.AddProjectMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.ProjectService.AddUserToProject(r.Context(), projectAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Member added to project successfully")
}

func RemoveProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	var payload services.RemoveProjectMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.ProjectService.RemoveUserFromProject(r.Context(), projectAuth, &payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Member removed from project successfully")
}

func ChangeProjectMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	projectAuth, ok := util.GetProjectAuthOrRespond(w, r, middleware.AuthSourceProject)
	if !ok {
		return
	}

	var payload services.ChangeProjectMemberRolePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.ProjectService.ChangeProjectMemberRole(r.Context(), projectAuth, payload)
	if err != nil {
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Member role updated successfully")
}
