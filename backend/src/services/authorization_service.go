package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
)

// IAuthorizer defines the interface for authorization operations.
// This interface allows for easy mocking in tests.
type IAuthorizer interface {
	authorizeProjectAction(ctx context.Context, projectAuth *authorization.ProjectAuthorization, permissions ...models.Permission) (*models.Project, *ServiceError)
	isInstanceAdmin(ctx context.Context, username string) (bool, error)
}

// AuthorizationService handles authorization-related operations.
// It provides methods to authorize actions based on project membership and permissions.
type AuthorizationService struct {
	db db.DB
}

type AuthorizationServiceDeps struct {
	DB db.DB
}

// NewAuthorizationService creates a new instance of authorizationServiceT.
func NewAuthorizationService(deps AuthorizationServiceDeps) *AuthorizationService {
	return &AuthorizationService{
		db: deps.DB,
	}
}

// AuthorizeProjectAction checks if a given user has the given permissions within the context of a specific project.
func (s *AuthorizationService) authorizeProjectAction(ctx context.Context, projectAuth *authorization.ProjectAuthorization, permissions ...models.Permission) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "AuthorizationService", "AuthorizeProjectAction")

	requestorUsername := projectAuth.Username

	// Does that project exist?
	project, err := s.internalGetProjectByID(ctx, projectAuth.ProjectID)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.internalGetUserByUsername(ctx, requestorUsername)
	if err != nil {
		return nil, err
	}

	// Is the requestor a member of that project?
	if !project.IsMember(user.ID) {
		logger.Warn().Str("username", requestorUsername).Str("project", project.Name).Msg("User is not a member of the project")
		return nil, &ServiceError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("user %s is not a member of project %s", requestorUsername, project.Name),
		}
	}

	// What permissions does the requestor have in that project?
	permissionIDs := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// Does the requestor have the required permissions?
	if !project.GetMember(user.ID).Role.HasPermissions(permissions...) {
		logger.Warn().Str("username", requestorUsername).Str("project", project.Name).Strs("permissions", permissionIDs).Msg("User does not have required permissions for project")
		return nil, &ServiceError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("user %s does not have required permissions for project %s", requestorUsername, project.Name),
		}
	}

	logger.Trace().Str("username", requestorUsername).Str("project", project.Name).Strs("permissions", permissionIDs).Msg("User has required permissions for project")
	return project, nil
}

// internalGetProjectByID retrieves a project by its display ID without authorization checks.
func (s *AuthorizationService) internalGetProjectByID(ctx context.Context, projectID uuid.UUID) (*models.Project, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "AuthorizationService", "internalGetProjectByID")

	project, err := s.db.Projects().GetProjectByID(ctx, projectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("projectID", projectID.String()).Msg("Project not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("project %s not found", projectID),
			}
		}
		logger.Error().Err(err).Str("projectID", projectID.String()).Msg("Error retrieving project")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving project %s: %w", projectID, err),
		}
	}

	return project, nil
}

// internalGetUserByUsername retrieves a user by username without authorization checks.
func (s *AuthorizationService) internalGetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "AuthorizationService", "internalGetUserByUsername")

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("username", username).Msg("User not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("user %s not found", username),
			}
		}
		logger.Error().Err(err).Str("username", username).Msg("Error retrieving user")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving user %s: %w", username, err),
		}
	}

	return user, nil
}

func (s *AuthorizationService) isInstanceAdmin(ctx context.Context, username string) (bool, error) {
	logger := MethodLoggerFromContext(ctx, "AuthorizationService", "isInstanceAdmin")

	logger.Trace().Str("username", username).Msg("Checking if user is instance admin")

	return true, nil
}
