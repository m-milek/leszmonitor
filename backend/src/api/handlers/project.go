package handlers

import (
	"encoding/json"
	"net/http"

	util "github.com/m-milek/leszmonitor/api/api_util"
	"github.com/m-milek/leszmonitor/api/services"
	"github.com/m-milek/leszmonitor/logging"
)

func CreateProjectHandler(w http.ResponseWriter, r *http.Request) {
	var payload services.CreateProjectPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to decode project creation payload")
		util.RespondError(w, http.StatusBadRequest, err)
	}

	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	project, err := services.ProjectService.CreateProject(r.Context(), orgAuth, payload)

	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to create project")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusCreated, project)
}

func GetProjectsOfOrgHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	projects, err := services.ProjectService.GetProjects(r.Context(), orgAuth)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to get projects")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, projects)
}

func GetProjectsByOrgID(w http.ResponseWriter, r *http.Request) {
	projectID := r.PathValue("projectId")

	if projectID == "" {
		logging.Api.Trace().Msg("Org's DisplayID is required")
		util.RespondMessage(w, http.StatusBadRequest, "Org's DisplayID is required")
		return
	}

	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	project, err := services.ProjectService.GetProjectsByOrgID(r.Context(), orgAuth, projectID)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to get projects by DisplayID")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, project)
}

func DeleteProjectHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	projectID := r.PathValue("projectId")

	if projectID == "" {
		logging.Api.Trace().Msg("Project's DisplayID is required for deletion")
		util.RespondMessage(w, http.StatusBadRequest, "Project's DisplayID is required")
		return
	}

	err := services.ProjectService.DeleteProject(r.Context(), orgAuth, projectID)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to delete project")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondMessage(w, http.StatusOK, "Project deleted successfully")
}

func UpdateProjectHandler(w http.ResponseWriter, r *http.Request) {
	orgAuth, ok := util.GetOrgAuthOrRespond(w, r)
	if !ok {
		return
	}

	projectID := r.PathValue("projectId")

	if projectID == "" {
		logging.Api.Trace().Msg("Project's DisplayID is required for update")
		util.RespondMessage(w, http.StatusBadRequest, "Project's DisplayID is required")
		return
	}

	var payload services.UpdateProjectPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		logging.Api.Trace().Err(err).Msg("Failed to decode project update payload")
		util.RespondError(w, http.StatusBadRequest, err)
		return
	}

	project, err := services.ProjectService.UpdateProject(r.Context(), orgAuth, projectID, &payload)
	if err != nil {
		logging.Api.Error().Err(err).Msg("Failed to update project")
		util.RespondError(w, err.Code, err.Err)
		return
	}

	util.RespondJSON(w, http.StatusOK, project)
}
