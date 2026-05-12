package handlers

import (
	"encoding/json"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/services"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload services.CreateProjectPayload
	if !util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	userClaims, ok := util.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	project, err := services.ProjectService.CreateProject(ctx, userClaims.Username, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusCreated, project)
}

func GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userClaims, ok := util.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	projects, err := services.ProjectService.GetProjectsForUser(ctx, userClaims.Username)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, projects)
}

func GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProjectSlug)
	if !ok {
		return
	}

	project, err := services.ProjectService.GetProjectByID(ctx, projectAuth)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, project)
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProjectSlug)
	if !ok {
		return
	}

	err := services.ProjectService.DeleteProject(ctx, projectAuth)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Project deleted successfully")
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProjectSlug)
	if !ok {
		return
	}

	var payload services.UpdateProjectPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	project, err := services.ProjectService.UpdateProject(ctx, projectAuth, &payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, project)
}

func AddProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProjectSlug)
	if !ok {
		return
	}

	var payload services.AddProjectMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.ProjectService.AddUserToProject(ctx, projectAuth, &payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Member added to project successfully")
}

func RemoveProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProjectSlug)
	if !ok {
		return
	}

	var payload services.RemoveProjectMemberPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.ProjectService.RemoveUserFromProject(ctx, projectAuth, &payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Member removed from project successfully")
}

func ChangeProjectMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := util.GetProjectAuthOrRespond(ctx, w, r, middleware.AuthSourceProjectSlug)
	if !ok {
		return
	}

	var payload services.ChangeProjectMemberRolePayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondMessage(ctx, w, http.StatusBadRequest, "Invalid request payload: "+err.Error())
		return
	}

	err := services.ProjectService.ChangeProjectMemberRole(ctx, projectAuth, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Member role updated successfully")
}
