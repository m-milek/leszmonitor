package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
)

type IAuthzMiddlewareService interface {
	CheckPermissionByID(
		ctx context.Context,
		username string,
		permission models.Permission,
	) (bool, error)
}

type AuthzMiddlewareService struct {
	db db.DB
}

func NewAuthzMiddlewareService(db db.DB) IAuthzMiddlewareService {
	return &AuthzMiddlewareService{db: db}
}

func (s *AuthzMiddlewareService) CheckPermissionByID(
	ctx context.Context,
	username string,
	permission models.Permission,
) (bool, error) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameAuthzMiddleware, "CheckPermissionByID")

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
		Msg("Checked monitor permission successfully")
	return hasPermission, nil
}
