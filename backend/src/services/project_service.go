package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
)

type IProjectService interface {
	CreateProject(ctx context.Context, ownerUsername string, payload CreateProjectPayload) (*models.Project, *ServiceError)
	GetProjectBySlug(ctx context.Context, projectSlug string) (*models.Project, *ServiceError)
	GetProjects(ctx context.Context, requestorUsername string, usernameQuery string) ([]models.Project, *ServiceError)
	DeleteProject(ctx context.Context, projectSlug string) *ServiceError
	UpdateProject(ctx context.Context, projectSlug string, payload UpdateProjectPayload) (*models.Project, *ServiceError)

	AddUserToProject(ctx context.Context, projectSlug string, payload AddProjectMemberPayload) *ServiceError
	RemoveUserFromProject(ctx context.Context, projectSlug string, payload RemoveProjectMemberPayload) *ServiceError
	ChangeProjectMemberRole(ctx context.Context, projectSlug string, payload ChangeProjectMemberRolePayload) *ServiceError
}

// ProjectService handles project-related CRUD operations including membership management.
type ProjectService struct {
	db          db.DB
	UserService IUserService // public so that we can do the circular dependency
}

type ProjectServiceDeps struct {
	DB          db.DB
	UserService IUserService
}

func NewProjectService(deps ProjectServiceDeps) *ProjectService {
	return &ProjectService{
		db:          deps.DB,
		UserService: deps.UserService,
	}
}

const ProjectServiceName = "ProjectService"

type CreateProjectPayload struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
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
func (s *ProjectService) CreateProject(ctx context.Context, ownerUsername string, payload CreateProjectPayload) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "CreateProject")

	user, err := s.db.Users().GetUserByUsername(ctx, ownerUsername)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, NewNotFoundError("user %s not found", ownerUsername)
		}
		return nil, NewInternalError("failed to retrieve user: %w", err)
	}

	project, err := models.NewProject(payload.Name, payload.Description, user.ID)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create new project")
		return nil, NewBadRequestError("invalid project data: %w", err)
	}

	var created *models.Project
	if txErr := s.db.WithTx(ctx, func(tx db.DB) error {
		if err := tx.Projects().InsertProject(ctx, project); err != nil {
			return err
		}
		var err error
		created, err = tx.Projects().GetProjectBySlug(ctx, project.Slug)
		return err
	}); txErr != nil {
		if errors.Is(txErr, db.ErrAlreadyExists) {
			return nil, NewConflictError("project with slug %s already exists", project.Slug)
		}
		logger.Error().Err(txErr).Msg("Failed to insert project")
		return nil, NewInternalError("failed to create project: %w", txErr)
	}

	logger.Info().Str("projectId", project.Slug).Msg("Project created successfully")
	return created, nil
}

// GetProjectBySlug retrieves a project by its slug.
func (s *ProjectService) GetProjectBySlug(ctx context.Context, projectSlug string) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "GetProjectBySlug")

	project, err := s.internalGetProjectBySlug(ctx, projectSlug)
	if err != nil {
		return nil, err
	}

	logger.Trace().Str("projectID", project.Slug).Msg("Retrieved project successfully")
	return project, nil
}

// GetProjects returns all projects the authenticated user is a member of.
func (s *ProjectService) GetProjects(ctx context.Context, requestorUsername string, usernameQuery string) ([]models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "GetProjects")

	user, err := s.db.Users().GetUserByUsername(ctx, requestorUsername)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, NewNotFoundError("user %s not found", requestorUsername)
		}
		return nil, NewInternalError("failed to retrieve user: %w", err)
	}

	getProjectsQuery := db.GetProjectsQuery{
		RequestingUserID: user.ID,
		MemberUsername:   usernameQuery,
	}

	projects, err := s.db.Projects().GetProjectsByQuery(ctx, getProjectsQuery)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get projects for user")
		return nil, NewInternalError("failed to get projects: %w", err)
	}

	logger.Info().Int("count", len(projects)).Str("requestingUser", requestorUsername).Str("userQuery", usernameQuery).Msg("Retrieved projects for user")
	return projects, nil
}

// DeleteProject deletes a project.
func (s *ProjectService) DeleteProject(ctx context.Context, projectSlug string) *ServiceError {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "DeleteProject")

	project, getErr := s.internalGetProjectBySlug(ctx, projectSlug)
	if getErr != nil {
		return getErr
	}

	deleted, err := s.db.Projects().DeleteProject(ctx, project.Slug)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to delete project")
		return NewInternalError("failed to delete project: %w", err)
	}
	if !deleted {
		return NewNotFoundError("project %s not found", project.Slug)
	}

	logger.Info().Str("projectID", project.Slug).Msg("Project deleted successfully")
	return nil
}

// UpdateProject updates a project's name/description.
func (s *ProjectService) UpdateProject(ctx context.Context, projectSlug string, payload UpdateProjectPayload) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "UpdateProject")

	oldProject, getErr := s.internalGetProjectBySlug(ctx, projectSlug)
	if getErr != nil {
		return nil, getErr
	}

	newProject := *oldProject
	newProject.Name = payload.Name
	newProject.Description = payload.Description
	newProject.SlugFromName.Init(newProject.Name)

	if _, err := s.db.Projects().UpdateProject(ctx, oldProject, &newProject); err != nil {
		logger.Error().Err(err).Msg("Failed to update project")
		return nil, NewInternalError("failed to update project: %w", err)
	}

	logger.Info().Str("projectID", oldProject.Slug).Msg("Project updated successfully")
	return &newProject, nil
}

// AddUserToProject adds a user to a project with a specified role.
func (s *ProjectService) AddUserToProject(ctx context.Context, projectSlug string, payload AddProjectMemberPayload) *ServiceError {

	project, getErr := s.internalGetProjectBySlug(ctx, projectSlug)
	if getErr != nil {
		return getErr
	}

	user, err := s.db.Users().GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return NewNotFoundError("user %s not found", payload.Username)
		}
		return NewInternalError("failed to retrieve user: %w", err)
	}

	if err := payload.Role.Validate(); err != nil {
		return NewBadRequestError("%w", err)
	}

	member, err := models.NewProjectMember(user.ID, payload.Role)
	if err != nil {
		return NewInternalError("failed to create member: %w", err)
	}

	_, err = s.db.Projects().AddMemberToProject(ctx, project.Slug, member)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return NewConflictError("user %s is already a member of project %s", payload.Username, project.Slug)
		}
		return NewInternalError("failed to add user to project: %w", err)
	}

	return nil
}

// RemoveUserFromProject removes a user from a project.
func (s *ProjectService) RemoveUserFromProject(ctx context.Context, projectSlug string, payload RemoveProjectMemberPayload) *ServiceError {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "RemoveUserFromProject")

	project, getErr := s.internalGetProjectBySlug(ctx, projectSlug)
	if getErr != nil {
		return getErr
	}

	user, serviceErr := s.UserService.GetUserByUsername(ctx, payload.Username)
	if serviceErr != nil {
		return serviceErr
	}

	member := project.GetMember(user.ID)
	if member == nil {
		return NewBadRequestError("user %s is not a member of project %s", payload.Username, project.Slug)
	}
	if member.Role == models.RoleOwner {
		logger.Warn().Str("username", payload.Username).Msg("Cannot remove project owner")
		return NewBadRequestError("cannot remove the project owner")
	}

	removed, err := s.db.Projects().RemoveMemberFromProject(ctx, project.Slug, user.ID)
	if err != nil {
		return NewInternalError("failed to remove user from project: %w", err)
	}
	if !removed {
		return NewNotFoundError("user %s is not a member of project %s", payload.Username, project.Slug)
	}

	return nil
}

// ChangeProjectMemberRole changes a member's role.
func (s *ProjectService) ChangeProjectMemberRole(ctx context.Context, projectSlug string, payload ChangeProjectMemberRolePayload) *ServiceError {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "ChangeProjectMemberRole")

	project, getErr := s.internalGetProjectBySlug(ctx, projectSlug)
	if getErr != nil {
		return getErr
	}

	user, serviceErr := s.UserService.GetUserByUsername(ctx, payload.Username)
	if serviceErr != nil {
		return serviceErr
	}

	if !project.IsMember(user.ID) {
		return NewBadRequestError("user %s is not a member of project %s", payload.Username, project.Slug)
	}

	if err := payload.Role.Validate(); err != nil {
		return NewBadRequestError("invalid role: %w", err)
	}

	if err := project.ChangeMemberRole(user.ID, payload.Role); err != nil {
		logger.Error().Err(err).Msg("Error changing role")
		return NewInternalError("error changing role: %w", err)
	}

	_, err := s.db.Projects().ChangeMemberRole(ctx, project.Slug, user.ID, payload.Role)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update project with new role")
		return NewInternalError("failed to update project with new role: %w", err)
	}

	return nil
}

func (s *ProjectService) internalGetProjectBySlug(ctx context.Context, projectID string) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, ProjectServiceName, "internalGetProjectBySlug")

	project, err := s.db.Projects().GetProjectBySlug(ctx, projectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("projectID", projectID).Msg("Project not found")
			return nil, NewNotFoundError("project %s not found", projectID)
		}
		return nil, NewInternalError("failed to get project: %w", err)
	}

	return project, nil
}
