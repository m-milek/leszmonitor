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
	"github.com/rs/zerolog"
)

// authorizationServiceT handles authorization-related operations.
// It provides methods to authorize actions based on team membership and permissions.
type authorizationServiceT struct {
	baseService
}

// NewAuthorizationService creates a new instance of authorizationServiceT.
func newAuthorizationService() *authorizationServiceT {
	return &authorizationServiceT{
		baseService{
			serviceLogger: logging.NewServiceLogger("AuthorizationService"),
			methodLoggers: make(map[string]zerolog.Logger),
		},
	}
}

// Checks if a given user has given permissions in the context of a specific team.
func (s *authorizationServiceT) authorizeTeamAction(ctx context.Context, teamAuth *middleware.TeamAuth, permissions ...models.Permission) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("AuthorizeTeamAction")

	requestorUsername := teamAuth.Username

	// Does that team exist?
	team, err := s.internalGetTeamByID(ctx, teamAuth.TeamID)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.internalGetUserByUsername(ctx, requestorUsername)
	if err != nil {
		return nil, err
	}

	// Is the requestor a member of that team?
	if !team.IsMember(user.ID) {
		logger.Warn().Str("username", requestorUsername).Str("team", team.Name).Msg("User is not a member of the team")
		return nil, &ServiceError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("user %s is not a member of team %s", requestorUsername, team.Name),
		}
	}

	// What permissions does the requestor have in that team?
	permissionIDs := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// Does the requestor have the required permissions?
	if !team.GetMember(user.ID).Role.HasPermissions(permissions...) {
		logger.Warn().Str("username", requestorUsername).Str("team", team.Name).Strs("permissions", permissionIDs).Msg("User does not have required permissions for team")
		return nil, &ServiceError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("user %s does not have required permissions for team %s", requestorUsername, team.Name),
		}
	}

	logger.Trace().Str("username", requestorUsername).Str("team", team.Name).Strs("permissions", permissionIDs).Msg("User has required permissions for team")
	return team, nil
}

// internalGetTeamByID retrieves a team by its display ID without authorization checks.
func (s *authorizationServiceT) internalGetTeamByID(ctx context.Context, teamID string) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("internalGetTeamByID")

	team, err := db.Get().Teams().GetTeamByDisplayID(ctx, teamID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("teamID", teamID).Msg("Team not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("team %s not found", teamID),
			}
		}
		logger.Error().Err(err).Str("teamID", teamID).Msg("Error retrieving team")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving team %s: %w", teamID, err),
		}
	}

	return team, nil
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
