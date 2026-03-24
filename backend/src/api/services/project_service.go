package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
)

type ProjectServiceT struct {
	baseService
}

// newProjectService creates a new instance of ProjectServiceT.
func newProjectService(service baseService) *ProjectServiceT {
	service.serviceLogger = logging.NewServiceLogger("ProjectService")
	return &ProjectServiceT{
		baseService: service,
	}
}

var ProjectService = newProjectService(newBaseService(newAuthorizationService(), "ProjectService"))

type UpdateProjectPayload struct {
	Name        string `json:"name"`        // Name of the project
	Description string `json:"description"` // Description of the project
}

type CreateProjectPayload struct {
	Name        string `json:"name"`        // Name of the project
	Description string `json:"description"` // Description of the project
}

// CreateProject creates a new project for the org in the provided OrgAuth.
func (s *ProjectServiceT) CreateProject(context context.Context, orgAuth *middleware.OrgAuth, payload CreateProjectPayload) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("CreateProject")

	org, authErr := s.authService.authorizeOrgAction(context, orgAuth, models.PermissionOrgEditor)
	if authErr != nil {
		return nil, authErr
	}

	project, err := models.NewProject(payload.Name, payload.Description, org)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create new project")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid project data: %w", err),
		}
	}

	err = s.getDB().Projects().InsertProject(context, project)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			logger.Warn().Str("projectId", project.DisplayID).Msg("Project already exists")
			return nil, &ServiceError{
				Code: http.StatusConflict,
				Err:  fmt.Errorf("project with DisplayID %s already exists", project.DisplayID),
			}
		}
		logger.Error().Err(err).Msg("Failed to create project")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create project: %w", err),
		}
	}

	createdProject, err := s.getDB().Projects().GetProjectByDisplayID(context, project.DisplayID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to fetch created project")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to fetch created project: %w", err),
		}
	}

	logger.Info().Str("projectId", project.DisplayID).Msg("Project created successfully")
	return createdProject, nil
}

// GetProjects retrieves all projects for the org in the provided OrgAuth.
func (s *ProjectServiceT) GetProjects(context context.Context, orgAuth *middleware.OrgAuth) ([]models.Project, *ServiceError) {
	logger := ProjectService.getMethodLogger("GetProjects")

	org, authErr := s.authService.authorizeOrgAction(context, orgAuth, models.PermissionOrgReader)
	if authErr != nil {
		return nil, authErr
	}

	projects, err := s.getDB().Projects().GetProjectsByOrgID(context, org)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get projects for org")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to get projects for org %s: %w", org.DisplayID, err),
		}
	}

	logger.Info().Int("count", len(projects)).Msg("Retrieved projects for org")
	return projects, nil
}

// GetProjectsByOrgID retrieves a specific project by its DisplayID for the org in the provided OrgAuth.
func (s *ProjectServiceT) GetProjectsByOrgID(context context.Context, orgAuth *middleware.OrgAuth, projectID string) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("GetProjectsByOrgID")

	if projectID == "" {
		logger.Warn().Msg("Project DisplayID is required to get project")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("project DisplayID is required"),
		}
	}

	_, authErr := s.authService.authorizeOrgAction(context, orgAuth, models.PermissionOrgReader)
	if authErr != nil {
		return nil, authErr
	}

	project, err := s.internalGetProjectByDisplayID(context, projectID)
	if err != nil {
		return nil, err
	}

	logger.Info().Str("projectID", project.DisplayID).Msg("Retrieved project by DisplayID")
	return project, nil
}

// DeleteProject deletes a specific project by its DisplayID for the org in the provided OrgAuth.
func (s *ProjectServiceT) DeleteProject(context context.Context, orgAuth *middleware.OrgAuth, projectID string) *ServiceError {
	logger := s.getMethodLogger("DeleteProject")

	org, authErr := s.authService.authorizeOrgAction(context, orgAuth, models.PermissionOrgEditor)
	if authErr != nil {
		return authErr
	}

	if projectID == "" {
		logger.Warn().Msg("Project DisplayID is required for deletion")
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("project DisplayID is required"),
		}
	}

	deleted, err := s.getDB().Projects().DeleteProject(context, org, projectID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete project")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete project with DisplayID %s: %w", projectID, err),
		}
	}

	if !deleted {
		logger.Warn().Str("projectID", projectID).Msg("Project not found for deletion")
		return &ServiceError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("project with DisplayID %s not found", projectID),
		}
	}

	logger.Info().Str("projectID", projectID).Msg("Project deleted successfully")
	return nil
}

// UpdateProject updates the details of a specific project by its DisplayID for the org in the provided OrgAuth.
func (s *ProjectServiceT) UpdateProject(ctx context.Context, orgAuth *middleware.OrgAuth, projectID string, payload *UpdateProjectPayload) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("UpdateProject")

	org, authErr := s.authService.authorizeOrgAction(ctx, orgAuth, models.PermissionOrgEditor)
	if authErr != nil {
		return nil, authErr
	}

	if projectID == "" {
		logger.Warn().Msg("Project DisplayID is required for update")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  errors.New("project DisplayID is required"),
		}
	}

	oldProject, err := s.internalGetProjectByDisplayID(ctx, projectID)
	if err != nil {
		return nil, err
	}

	newProject := *oldProject
	newProject.Name = payload.Name
	newProject.Description = payload.Description
	newProject.DisplayIDFromName.Init(newProject.Name)

	_, updateErr := s.getDB().Projects().UpdateProject(ctx, org, oldProject, &newProject)
	if updateErr != nil {
		logger.Error().Err(updateErr).Msg("Failed to update project")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update project with DisplayID %s: %w", projectID, updateErr),
		}
	}

	logger.Info().Str("projectID", oldProject.DisplayID).Msg("Project updated successfully")
	return &newProject, nil
}

func (s *ProjectServiceT) internalGetProjectByDisplayID(ctx context.Context, projectID string) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("internalGetProjectByDisplayID")

	project, err := s.getDB().Projects().GetProjectByDisplayID(ctx, projectID)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("projectID", projectID).Msg("Project not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("project with DisplayID %s not found", projectID),
			}
		}
		logger.Error().Err(err).Msg("Failed to get project")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to get project: %w", err),
		}
	}

	return project, nil
}
