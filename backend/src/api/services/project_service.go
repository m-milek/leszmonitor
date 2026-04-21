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

// ProjectServiceT handles project-related CRUD operations including membership management.
type ProjectServiceT struct {
	baseService
}

func newProjectService(service baseService) *ProjectServiceT {
	service.serviceLogger = logging.NewServiceLogger("ProjectService")
	return &ProjectServiceT{baseService: service}
}

var ProjectService = newProjectService(newBaseService(newAuthorizationService(), "ProjectService"))

type CreateProjectPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateProjectPayload struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AddProjectMemberPayload struct {
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
}

type RemoveProjectMemberPayload struct {
	Username string `json:"username"`
}

type ChangeProjectMemberRolePayload struct {
	Username string      `json:"username"`
	Role     models.Role `json:"role"`
}

// CreateProject creates a new project owned by the authenticated user.
func (s *ProjectServiceT) CreateProject(ctx context.Context, ownerUsername string, payload CreateProjectPayload) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("CreateProject")

	user, err := db.Get().Users().GetUserByUsername(ctx, ownerUsername)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("user %s not found", ownerUsername)}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve user: %w", err)}
	}

	project, err := models.NewProject(payload.Name, payload.Description, user.ID)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create new project")
		return nil, &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid project data: %w", err)}
	}

	if err = s.getDB().Projects().InsertProject(ctx, project); err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return nil, &ServiceError{Code: http.StatusConflict, Err: fmt.Errorf("project with slug %s already exists", project.Slug)}
		}
		logger.Error().Err(err).Msg("Failed to insert project")
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to create project: %w", err)}
	}

	created, err := s.getDB().Projects().GetProjectBySlug(ctx, project.Slug)
	if err != nil {
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to fetch created project: %w", err)}
	}

	logger.Info().Str("projectId", project.Slug).Msg("Project created successfully")
	return created, nil
}

// GetProjectByID retrieves a project by its slug, authorizing the requesting user.
func (s *ProjectServiceT) GetProjectByID(ctx context.Context, projectAuth *middleware.ProjectAuth) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("GetProjectByID")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionProjectReader)
	if authErr != nil {
		return nil, authErr
	}

	logger.Trace().Str("projectID", project.Slug).Msg("Retrieved project successfully")
	return project, nil
}

// GetProjectsForUser returns all projects the authenticated user is a member of.
func (s *ProjectServiceT) GetProjectsForUser(ctx context.Context, username string) ([]models.Project, *ServiceError) {
	logger := s.getMethodLogger("GetProjectsForUser")

	user, err := db.Get().Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("user %s not found", username)}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve user: %w", err)}
	}

	projects, err := s.getDB().Projects().GetProjectsByUserID(ctx, user.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get projects for user")
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to get projects: %w", err)}
	}

	logger.Info().Int("count", len(projects)).Msg("Retrieved projects for user")
	return projects, nil
}

// DeleteProject deletes a project. Requires ProjectAdmin permission.
func (s *ProjectServiceT) DeleteProject(ctx context.Context, projectAuth *middleware.ProjectAuth) *ServiceError {
	logger := s.getMethodLogger("DeleteProject")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionProjectAdmin)
	if authErr != nil {
		return authErr
	}

	deleted, err := s.getDB().Projects().DeleteProject(ctx, project.Slug)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete project")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to delete project: %w", err)}
	}
	if !deleted {
		return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("project %s not found", project.Slug)}
	}

	logger.Info().Str("projectID", project.Slug).Msg("Project deleted successfully")
	return nil
}

// UpdateProject updates a project's name/description. Requires ProjectEditor permission.
func (s *ProjectServiceT) UpdateProject(ctx context.Context, projectAuth *middleware.ProjectAuth, payload *UpdateProjectPayload) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("UpdateProject")

	oldProject, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionProjectEditor)
	if authErr != nil {
		return nil, authErr
	}

	newProject := *oldProject
	newProject.Name = payload.Name
	newProject.Description = payload.Description
	newProject.SlugFromName.Init(newProject.Name)

	if _, err := s.getDB().Projects().UpdateProject(ctx, oldProject, &newProject); err != nil {
		logger.Error().Err(err).Msg("Failed to update project")
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to update project: %w", err)}
	}

	logger.Info().Str("projectID", oldProject.Slug).Msg("Project updated successfully")
	return &newProject, nil
}

// AddUserToProject adds a user to a project with a specified role. Requires ProjectEditor permission.
func (s *ProjectServiceT) AddUserToProject(ctx context.Context, projectAuth *middleware.ProjectAuth, payload *AddProjectMemberPayload) *ServiceError {
	logger := s.getMethodLogger("AddUserToProject")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionProjectEditor)
	if authErr != nil {
		return authErr
	}

	user, err := db.Get().Users().GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("user %s not found", payload.Username)}
		}
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve user: %w", err)}
	}

	if err := payload.Role.Validate(); err != nil {
		return &ServiceError{Code: http.StatusBadRequest, Err: err}
	}

	member, err := models.NewProjectMember(user.ID, payload.Role)
	if err != nil {
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to create member: %w", err)}
	}

	_, err = s.getDB().Projects().AddMemberToProject(ctx, project.Slug, member)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return &ServiceError{Code: http.StatusConflict, Err: fmt.Errorf("user %s is already a member of project %s", payload.Username, project.Slug)}
		}
		logger.Error().Err(err).Msg("Failed to add user to project")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to add user to project: %w", err)}
	}

	return nil
}

// RemoveUserFromProject removes a user from a project. Requires ProjectEditor permission.
func (s *ProjectServiceT) RemoveUserFromProject(ctx context.Context, projectAuth *middleware.ProjectAuth, payload *RemoveProjectMemberPayload) *ServiceError {
	logger := s.getMethodLogger("RemoveUserFromProject")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionProjectEditor)
	if authErr != nil {
		return authErr
	}

	user, serviceErr := UserService.GetUserByUsername(ctx, payload.Username)
	if serviceErr != nil {
		return serviceErr
	}

	member := project.GetMember(user.ID)
	if member == nil {
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("user %s is not a member of project %s", payload.Username, project.Slug)}
	}
	if member.Role == models.RoleOwner {
		logger.Warn().Str("username", payload.Username).Msg("Cannot remove project owner")
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("cannot remove the project owner")}
	}

	removed, err := s.getDB().Projects().RemoveMemberFromProject(ctx, project.Slug, user.ID)
	if err != nil {
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to remove user from project: %w", err)}
	}
	if !removed {
		return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("user %s is not a member of project %s", payload.Username, project.Slug)}
	}

	return nil
}

// ChangeProjectMemberRole changes a member's role. Requires ProjectAdmin permission.
func (s *ProjectServiceT) ChangeProjectMemberRole(ctx context.Context, projectAuth *middleware.ProjectAuth, payload ChangeProjectMemberRolePayload) *ServiceError {
	logger := s.getMethodLogger("ChangeProjectMemberRole")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionProjectAdmin)
	if authErr != nil {
		return authErr
	}

	user, serviceErr := UserService.GetUserByUsername(ctx, payload.Username)
	if serviceErr != nil {
		return serviceErr
	}

	if !project.IsMember(user.ID) {
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("user %s is not a member of project %s", payload.Username, project.Slug)}
	}

	if err := payload.Role.Validate(); err != nil {
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid role: %w", err)}
	}

	if err := project.ChangeMemberRole(user.ID, payload.Role); err != nil {
		logger.Error().Err(err).Msg("Error changing role")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("error changing role: %w", err)}
	}

	_, err := s.getDB().Projects().UpdateProject(ctx, project, project)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update project with new role")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to update project with new role: %w", err)}
	}

	return nil
}

func (s *ProjectServiceT) internalGetProjectBySlug(ctx context.Context, projectID string) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("internalGetProjectBySlug")

	project, err := s.getDB().Projects().GetProjectBySlug(ctx, projectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("projectID", projectID).Msg("Project not found")
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("project %s not found", projectID)}
		}
		logger.Error().Err(err).Msg("Failed to get project")
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to get project: %w", err)}
	}

	return project, nil
}
