package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models"
	"github.com/rs/zerolog"
)

// authorizationServiceT handles authorization-related operations.
// It provides methods to authorize actions based on project membership and permissions.
type authorizationServiceT struct {
	baseService
}

// newAuthorizationService creates a new instance of authorizationServiceT.
func newAuthorizationService() *authorizationServiceT {
	return &authorizationServiceT{
		baseService{
			serviceLogger: log.NewServiceLogger("AuthorizationService"),
			methodLoggers: make(map[string]zerolog.Logger),
		},
	}
}

// authorizeProjectAction checks if a given user has the given permissions within the context of a specific project.
func (s *authorizationServiceT) authorizeProjectAction(ctx context.Context, projectAuth *middleware.ProjectAuth, permissions ...models.Permission) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("authorizeProjectAction")

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
func (s *authorizationServiceT) internalGetProjectByID(ctx context.Context, projectID string) (*models.Project, *ServiceError) {
	logger := s.getMethodLogger("internalGetProjectByID")

	project, err := db.Get().Projects().GetProjectBySlug(ctx, projectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("projectID", projectID).Msg("Project not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("project %s not found", projectID),
			}
		}
		logger.Error().Err(err).Str("projectID", projectID).Msg("Error retrieving project")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving project %s: %w", projectID, err),
		}
	}

	return project, nil
}

// internalGetUserByUsername retrieves a user by username without authorization checks.
func (s *authorizationServiceT) internalGetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError) {
	logger := s.getMethodLogger("internalGetUserByUsername")

	user, err := db.Get().Users().GetUserByUsername(ctx, username)
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
