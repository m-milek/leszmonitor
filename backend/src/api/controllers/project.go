package controllers

import (
	"encoding/json"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/services"
)

type ProjectAPIController struct {
	service services.IProjectService
}

func NewProjectAPIController(service services.IProjectService) ProjectAPIController {
	return ProjectAPIController{
		service: service,
	}
}

func (h *ProjectAPIController) CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var payload services.CreateProjectPayload
	if !util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	userClaims, ok := authorization.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	project, err := h.service.CreateProject(ctx, userClaims.Username, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusCreated, project)
}

func (h *ProjectAPIController) GetProjectsHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userClaims, ok := authorization.ExtractUserOrRespond(ctx, w, r)
	if !ok {
		return
	}

	requestingUser := userClaims.Username
	usernameQuery := r.URL.Query().Get("username")

	projects, err := h.service.GetProjects(ctx, requestingUser, usernameQuery)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, projects)
}

func (h *ProjectAPIController) GetProjectByIDHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		ProjectSlug: r.PathValue("projectSlug"),
	})
	if !ok {
		return
	}

	project, err := h.service.GetProjectByID(ctx, projectAuth)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, project)
}

func (h *ProjectAPIController) DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		ProjectSlug: r.URL.Query().Get("projectSlug"),
	})
	if !ok {
		return
	}

	err := h.service.DeleteProject(ctx, projectAuth)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Project deleted successfully")
}

func (h *ProjectAPIController) UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		ProjectSlug: r.URL.Query().Get("projectSlug"),
	})
	if !ok {
		return
	}

	var payload services.UpdateProjectPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		util.RespondError(ctx, w, http.StatusBadRequest, err)
		return
	}

	project, err := h.service.UpdateProject(ctx, projectAuth, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondJSON(ctx, w, http.StatusOK, project)
}

func (h *ProjectAPIController) AddProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		ProjectSlug: r.PathValue("projectSlug"),
	})
	if !ok {
		return
	}

	var payload services.AddProjectMemberPayload
	if !util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	err := h.service.AddUserToProject(ctx, projectAuth, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Member added to project successfully")
}

func (h *ProjectAPIController) RemoveProjectMemberHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		ProjectSlug: r.URL.Query().Get("projectSlug"),
	})
	if !ok {
		return
	}

	var payload services.RemoveProjectMemberPayload
	if util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	err := h.service.RemoveUserFromProject(ctx, projectAuth, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Member removed from project successfully")
}

func (h *ProjectAPIController) ChangeProjectMemberRoleHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	projectAuth, ok := authorization.NewOrRespond(ctx, w, authorization.Payload{
		ProjectSlug: r.URL.Query().Get("projectSlug"),
	})
	if !ok {
		return
	}

	var payload services.ChangeProjectMemberRolePayload
	if util.DecodeJSONOrRespond(ctx, w, r, &payload) {
		return
	}

	err := h.service.ChangeProjectMemberRole(ctx, projectAuth, payload)
	if err != nil {
		util.RespondError(ctx, w, err.Code, err.Err)
		return
	}

	util.RespondMessage(ctx, w, http.StatusOK, "Member role updated successfully")
}
