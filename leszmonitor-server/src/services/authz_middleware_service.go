package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/log"
)

type IAuthzMiddlewareService interface {
	// CheckUserPermission reports whether the user's role grants the given permission.
	CheckUserPermission(
		ctx context.Context,
		username string,
		permission auth.Permission,
	) (bool, error)
}

type AuthzMiddlewareService struct {
	db db.DB
}

func NewAuthzMiddlewareService(db db.DB) IAuthzMiddlewareService {
	return &AuthzMiddlewareService{db: db}
}

func (s *AuthzMiddlewareService) CheckUserPermission(
	ctx context.Context,
	username string,
	permission auth.Permission,
) (bool, error) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameAuthzMiddleware, "CheckUserPermission")

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", username).Msg("User not found")
			return false, nil
		}
		logger.Error().Err(err).Str("username", username).Msg("Failed to retrieve user")
		return false, err
	}

	hasPermission := user.Role.HasPermissions(permission)
	logger.Trace().
		Str("username", username).
		Interface("permission", permission).
		Bool("hasPermission", hasPermission).
		Msg("Checked user permission successfully")
	return hasPermission, nil
}
