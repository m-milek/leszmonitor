package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/security"
)

type IProjectService interface {
	CreateProject(
		ctx context.Context,
		ownerUsername string,
		payload CreateProjectPayload,
	) (*models.Project, *ServiceError)

	GetProjectBySlug(ctx context.Context, projectSlug string) (*models.Project, *ServiceError)
	GetProjects(ctx context.Context, requestorUsername string, usernameQuery string) ([]models.Project, *ServiceError)
	DeleteProject(ctx context.Context, projectSlug string) *ServiceError
	UpdateProject(
		ctx context.Context,
		projectSlug string,
		payload UpdateProjectPayload,
	) (*models.Project, *ServiceError)

	AddUserToProject(ctx context.Context, projectSlug string, payload AddProjectMemberPayload) *ServiceError
	RemoveUserFromProject(ctx context.Context, projectSlug string, payload RemoveProjectMemberPayload) *ServiceError
	ChangeProjectMemberRole(
		ctx context.Context,
		projectSlug string,
		payload ChangeProjectMemberRolePayload,
	) *ServiceError
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
func (s *ProjectService) CreateProject(
	ctx context.Context,
	ownerUsername string,
	payload CreateProjectPayload,
) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "CreateProject")

	user, err := s.db.Users().GetUserByUsername(ctx, ownerUsername)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("ownerUsername", ownerUsername).Msg("User not found")
			return nil, NewNotFoundError("user %s not found", ownerUsername)
		}
		logger.Error().Err(err).Str("ownerUsername", ownerUsername).Msg("Failed to retrieve user")
		return nil, NewInternalError("failed to retrieve user: %w", err)
	}

	project, err := models.NewProject(payload.Name, payload.Description, user.ID)
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to create new project")
		return nil, NewBadRequestError("invalid project data: %w", err)
	}

	created, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (*models.Project, *security.AuditLogParams, error) {
		if err := tx.Projects().InsertProject(ctx, project); err != nil {
			return nil, nil, err
		}

		createdProject, err := tx.Projects().GetProjectBySlug(ctx, project.Slug)
		if err != nil {
			return nil, nil, err
		}

		params := &security.AuditLogParams{
			Username:  &ownerUsername,
			ProjectID: &createdProject.ID,
			Action:    security.ActionCreateProject,
			IsSuccess: true,
			Summary:   fmt.Sprintf("Project %s created", createdProject.Name),
			After:     createdProject,
		}

		return createdProject, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrAlreadyExists) {
			logger.Error().Str("slug", project.Slug).Msg("Project with slug already exists")
			return nil, NewConflictError("project with slug %s already exists", project.Slug)
		}
		logger.Error().Err(txErr).Msg("Failed to insert project")
		return nil, NewInternalError("failed to create project: %w", txErr)
	}

	logger.Debug().Str("projectId", project.Slug).Msg("Project created successfully")
	return created, nil
}

// GetProjectBySlug retrieves a project by its slug.
func (s *ProjectService) GetProjectBySlug(ctx context.Context, projectSlug string) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "GetProjectBySlug")

	project, err := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if err != nil {
		logger.Error().Err(err).Str("projectSlug", projectSlug).Msg("Failed to get project by slug")
		return nil, err
	}

	logger.Debug().Str("projectID", project.Slug).Msg("Retrieved project successfully")
	return project, nil
}

// GetProjects returns all projects the authenticated user is a member of.
func (s *ProjectService) GetProjects(
	ctx context.Context,
	requestorUsername string,
	usernameQuery string,
) ([]models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "GetProjects")

	user, err := s.db.Users().GetUserByUsername(ctx, requestorUsername)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("requestorUsername", requestorUsername).Msg("User not found")
			return nil, NewNotFoundError("user %s not found", requestorUsername)
		}
		logger.Error().Err(err).Str("requestorUsername", requestorUsername).Msg("Failed to retrieve user")
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

	logger.Debug().
		Int("count", len(projects)).
		Str("requestingUser", requestorUsername).
		Str("userQuery", usernameQuery).
		Msg("Retrieved projects for user successfully")
	return projects, nil
}

// DeleteProject deletes a project.
func (s *ProjectService) DeleteProject(ctx context.Context, projectSlug string) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "DeleteProject")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	project, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		return getErr
	}

	deleted, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (bool, *security.AuditLogParams, error) {
		deleted, err := tx.Projects().DeleteProject(ctx, project.Slug)
		if err != nil {
			return false, nil, err
		}
		if !deleted {
			return false, nil, db.ErrNotFound
		}

		params := &security.AuditLogParams{
			Username:  &userClaims.Username,
			ProjectID: &project.ID,
			Action:    security.ActionDeleteProject,
			IsSuccess: true,
			Summary:   fmt.Sprintf("Project %s deleted", project.Name),
			Before:    project,
		}

		return true, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrNotFound) {
			logger.Error().Str("projectID", project.Slug).Msg("Project not found for deletion")
			return NewNotFoundError("project %s not found", project.Slug)
		}
		logger.Error().Err(txErr).Msg("Failed to delete project")
		return NewInternalError("failed to delete project: %w", txErr)
	}
	if !deleted {
		logger.Error().Str("projectID", project.Slug).Msg("Project not found for deletion")
		return NewNotFoundError("project %s not found", project.Slug)
	}

	logger.Debug().Str("projectID", project.Slug).Msg("Project deleted successfully")
	return nil
}

// UpdateProject updates a project's name/description.
func (s *ProjectService) UpdateProject(
	ctx context.Context,
	projectSlug string,
	payload UpdateProjectPayload,
) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "UpdateProject")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, NewUnauthorizedError("user claims not found in context")
	}

	oldProject, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		return nil, getErr
	}

	newProject := *oldProject
	newProject.Name = payload.Name
	newProject.Description = payload.Description
	newProject.SlugFromName.Init(newProject.Name)

	_, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (*models.Project, *security.AuditLogParams, error) {
		if _, err := tx.Projects().UpdateProject(ctx, oldProject, &newProject); err != nil {
			return nil, nil, err
		}

		updatedProject, err := tx.Projects().GetProjectBySlug(ctx, newProject.Slug)
		if err != nil {
			return nil, nil, err
		}

		params := &security.AuditLogParams{
			Username:  &userClaims.Username,
			ProjectID: &updatedProject.ID,
			Action:    security.ActionUpdateProject,
			IsSuccess: true,
			Summary:   fmt.Sprintf("Project %s updated", updatedProject.Name),
			Before:    oldProject,
			After:     updatedProject,
		}

		return updatedProject, params, nil
	})
	if txErr != nil {
		logger.Error().Err(txErr).Msg("Failed to update project")
		return nil, NewInternalError("failed to update project: %w", txErr)
	}

	logger.Debug().Str("projectID", oldProject.Slug).Msg("Project updated successfully")
	return &newProject, nil
}

// AddUserToProject adds a user to a project with a specified role.
func (s *ProjectService) AddUserToProject(
	ctx context.Context,
	projectSlug string,
	payload AddProjectMemberPayload,
) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "AddUserToProject")

	project, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		return getErr
	}

	user, err := s.db.Users().GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", payload.Username).Msg("User not found")
			return NewNotFoundError("user %s not found", payload.Username)
		}
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to retrieve user")
		return NewInternalError("failed to retrieve user: %w", err)
	}

	if err := payload.Role.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid role")
		return NewBadRequestError("%w", err)
	}

	member, err := models.NewProjectMember(user.ID, payload.Role)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create member")
		return NewInternalError("failed to create member: %w", err)
	}

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	err = db.WithAuditedVoidTx(ctx, s.db, func(tx db.DB) (*security.AuditLogParams, error) {
		_, err := tx.Projects().AddMemberToProject(ctx, project.Slug, member)
		if err != nil {
			return nil, err
		}

		params := &security.AuditLogParams{
			Username:  &userClaims.Username,
			ProjectID: &project.ID,
			Action:    security.ActionAddProjectMember,
			IsSuccess: true,
			Summary:   fmt.Sprintf("User %s added to project %s with role %s", payload.Username, project.Name, payload.Role),
			After:     member,
		}
		return params, nil
	})
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			logger.Error().
				Str("username", payload.Username).
				Str("projectSlug", project.Slug).
				Msg("User is already a member of project")
			return NewConflictError("user %s is already a member of project %s", payload.Username, project.Slug)
		}
		logger.Error().Err(err).Msg("Failed to add user to project")
		return NewInternalError("failed to add user to project: %w", err)
	}

	logger.Debug().
		Str("username", payload.Username).
		Str("projectSlug", project.Slug).
		Msg("User added to project successfully")
	return nil
}

// RemoveUserFromProject removes a user from a project.
func (s *ProjectService) RemoveUserFromProject(
	ctx context.Context,
	projectSlug string,
	payload RemoveProjectMemberPayload,
) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "RemoveUserFromProject")

	project, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		return getErr
	}

	user, serviceErr := s.UserService.GetUserByUsername(ctx, payload.Username)
	if serviceErr != nil {
		return serviceErr
	}

	member := project.GetMember(user.ID)
	if member == nil {
		logger.Error().
			Str("username", payload.Username).
			Str("projectSlug", project.Slug).
			Msg("User is not a member of project")
		return NewBadRequestError(FormatUserIsNotAMemberOfProject, payload.Username, project.Slug)
	}
	if member.Role == models.RoleOwner {
		logger.Error().Str("username", payload.Username).Msg("Cannot remove project owner")
		return NewBadRequestError("cannot remove the project owner")
	}

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	removed, err := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (bool, *security.AuditLogParams, error) {
		removed, err := tx.Projects().RemoveMemberFromProject(ctx, project.Slug, user.ID)
		if err != nil {
			return false, nil, err
		}
		if !removed {
			return false, nil, nil
		}

		params := &security.AuditLogParams{
			Username:  &userClaims.Username,
			ProjectID: &project.ID,
			Action:    security.ActionRemoveProjectMember,
			IsSuccess: true,
			Summary:   fmt.Sprintf("User %s removed from project %s", payload.Username, project.Name),
			Before:    *member,
		}
		return true, params, nil
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to remove user from project")
		return NewInternalError("failed to remove user from project: %w", err)
	}
	if !removed {
		logger.Error().
			Str("username", payload.Username).
			Str("projectSlug", project.Slug).
			Msg("User is not a member of project")
		return NewNotFoundError(FormatUserIsNotAMemberOfProject, payload.Username, project.Slug)
	}

	logger.Debug().
		Str("username", payload.Username).
		Str("projectSlug", project.Slug).
		Msg("User removed from project successfully")
	return nil
}

// ChangeProjectMemberRole changes a member's role.
func (s *ProjectService) ChangeProjectMemberRole(
	ctx context.Context,
	projectSlug string,
	payload ChangeProjectMemberRolePayload,
) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "ChangeProjectMemberRole")

	project, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		return getErr
	}

	user, serviceErr := s.UserService.GetUserByUsername(ctx, payload.Username)
	if serviceErr != nil {
		return serviceErr
	}

	if !project.IsMember(user.ID) {
		logger.Error().
			Str("username", payload.Username).
			Str("projectSlug", project.Slug).
			Msg("User is not a member of project")
		return NewBadRequestError(FormatUserIsNotAMemberOfProject, payload.Username, project.Slug)
	}

	if err := payload.Role.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid role")
		return NewBadRequestError("invalid role: %w", err)
	}

	oldMember := *project.GetMember(user.ID)

	if err := project.ChangeMemberRole(user.ID, payload.Role); err != nil {
		logger.Error().Err(err).Msg("Error changing role")
		return NewInternalError("error changing role: %w", err)
	}

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	err := db.WithAuditedVoidTx(ctx, s.db, func(tx db.DB) (*security.AuditLogParams, error) {
		_, err := tx.Projects().ChangeMemberRole(ctx, project.Slug, user.ID, payload.Role)
		if err != nil {
			return nil, err
		}

		params := &security.AuditLogParams{
			Username:  &userClaims.Username,
			ProjectID: &project.ID,
			Action:    security.ActionUpdateProjectMember,
			IsSuccess: true,
			Summary:   fmt.Sprintf("User %s role changed to %s in project %s", payload.Username, payload.Role, project.Name),
			Before:    oldMember,
			After:     *project.GetMember(user.ID),
		}
		return params, nil
	})
	if err != nil {
		logger.Error().Err(err).Msg("Failed to update project with new role")
		return NewInternalError("failed to update project with new role: %w", err)
	}

	logger.Debug().
		Str("username", payload.Username).
		Str("projectSlug", project.Slug).
		Str("role", string(payload.Role)).
		Msg("User role changed successfully")
	return nil
}

func internalGetProjectBySlug(ctx context.Context, dbClient db.DB, projectID string) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameProject, "internalGetProjectBySlug")
	project, err := dbClient.Projects().GetProjectBySlug(ctx, projectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("projectID", projectID).Msg("Project not found")
			return nil, NewNotFoundError("project %s not found", projectID)
		}
		logger.Error().Err(err).Str("projectID", projectID).Msg("Failed to get project")
		return nil, NewInternalError("failed to get project: %w", err)
	}

	return project, nil
}
