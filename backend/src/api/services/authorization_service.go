package services

import (
	"context"
	"fmt"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
	"github.com/rs/zerolog"
	"net/http"
)

// AuthorizationServiceT handles authorization-related operations.
// It provides methods to authorize actions based on team membership and permissions.
type AuthorizationServiceT struct {
	BaseService
}

// NewAuthorizationService creates a new instance of AuthorizationServiceT.
func newAuthorizationService() *AuthorizationServiceT {
	return &AuthorizationServiceT{
		BaseService{
			serviceLogger: logging.NewServiceLogger("AuthorizationService"),
			methodLoggers: make(map[string]zerolog.Logger),
		},
	}
}

// Checks if a given user has given permissions in the context of a specific team.
func (s *AuthorizationServiceT) authorizeTeamAction(ctx context.Context, teamAuth *middleware.TeamAuth, permissions ...models.Permission) (*models.Team, *ServiceError) {
	logger := s.getMethodLogger("AuthorizeTeamAction")

	requestorUsername := teamAuth.Username

	// Does that team exist?
	team, err := TeamService.internalGetTeamById(ctx, teamAuth.TeamId)
	if err != nil {
		return nil, err
	}

	user, err := UserService.GetUserByUsername(ctx, requestorUsername)
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
