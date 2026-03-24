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
// It provides methods to authorize actions based on org membership and permissions.
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

// Checks if a given user has given permissions in the context of a specific org.
func (s *authorizationServiceT) authorizeOrgAction(ctx context.Context, orgAuth *middleware.OrgAuth, permissions ...models.Permission) (*models.Org, *ServiceError) {
	logger := s.getMethodLogger("authorizeOrgAction")

	requestorUsername := orgAuth.Username

	// Does that org exist?
	org, err := s.internalGetOrgByID(ctx, orgAuth.OrgID)
	if err != nil {
		return nil, err
	}

	// Get user
	user, err := s.internalGetUserByUsername(ctx, requestorUsername)
	if err != nil {
		return nil, err
	}

	// Is the requestor a member of that org?
	if !org.IsMember(user.ID) {
		logger.Warn().Str("username", requestorUsername).Str("org", org.Name).Msg("User is not a member of the org")
		return nil, &ServiceError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("user %s is not a member of org %s", requestorUsername, org.Name),
		}
	}

	// What permissions does the requestor have in that org?
	permissionIDs := make([]string, 0, len(permissions))
	for _, perm := range permissions {
		permissionIDs = append(permissionIDs, perm.ID)
	}

	// Does the requestor have the required permissions?
	if !org.GetMember(user.ID).Role.HasPermissions(permissions...) {
		logger.Warn().Str("username", requestorUsername).Str("org", org.Name).Strs("permissions", permissionIDs).Msg("User does not have required permissions for org")
		return nil, &ServiceError{
			Code: http.StatusForbidden,
			Err:  fmt.Errorf("user %s does not have required permissions for org %s", requestorUsername, org.Name),
		}
	}

	logger.Trace().Str("username", requestorUsername).Str("org", org.Name).Strs("permissions", permissionIDs).Msg("User has required permissions for org")
	return org, nil
}

// internalGetOrgByID retrieves an org by its display ID without authorization checks.
func (s *authorizationServiceT) internalGetOrgByID(ctx context.Context, orgID string) (*models.Org, *ServiceError) {
	logger := s.getMethodLogger("internalGetOrgByID")

	org, err := db.Get().Orgs().GetOrgByDisplayID(ctx, orgID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("orgID", orgID).Msg("Org not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("org %s not found", orgID),
			}
		}
		logger.Error().Err(err).Str("orgID", orgID).Msg("Error retrieving org")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving org %s: %w", orgID, err),
		}
	}

	return org, nil
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
