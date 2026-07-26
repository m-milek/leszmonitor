package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
)

type IAuthzMiddlewareService interface {
	CheckProjectPermissionBySlug(ctx context.Context, username, projectSlug string, permission models.Permission) (bool, error)
	CheckProjectPermissionByID(ctx context.Context, username, projectID string, permission models.Permission) (bool, error)
	CheckMonitorPermissionByID(ctx context.Context, username, monitorID string, permission models.Permission) (bool, error)
}

type AuthzMiddlewareService struct {
	db db.DB
}

func NewAuthzMiddlewareService(db db.DB) IAuthzMiddlewareService {
	return &AuthzMiddlewareService{db: db}
}

func (s *AuthzMiddlewareService) CheckProjectPermissionBySlug(ctx context.Context, username, projectSlug string, permission models.Permission) (bool, error) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameAuthzMiddleware, "CheckProjectPermissionBySlug")

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", username).Msg("User not found")
			return false, nil
		}
		logger.Error().Err(err).Str("username", username).Msg("Failed to retrieve user")
		return false, err
	}

	project, err := s.db.Projects().GetProjectBySlug(ctx, projectSlug)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("projectSlug", projectSlug).Msg("Project not found")
			return false, nil
		}
		logger.Error().Err(err).Str("projectSlug", projectSlug).Msg("Failed to retrieve project")
		return false, err
	}

	if !project.IsMember(user.ID) {
		logger.Debug().Str("username", username).Str("projectSlug", projectSlug).Msg("User is not a member of the project")
		return false, nil
	}

	hasPermission := project.GetMember(user.ID).Role.HasPermissions(permission)
	logger.Debug().Str("username", username).Str("projectSlug", projectSlug).Interface("permission", permission).Bool("hasPermission", hasPermission).Msg("Checked project permission successfully")
	return hasPermission, nil
}

func (s *AuthzMiddlewareService) CheckProjectPermissionByID(ctx context.Context, username, projectID string, permission models.Permission) (bool, error) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameAuthzMiddleware, "CheckProjectPermissionByID")

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", username).Msg("User not found")
			return false, nil
		}
		logger.Error().Err(err).Str("username", username).Msg("Failed to retrieve user")
		return false, err
	}

	parsedUUID, err := uuid.Parse(projectID)
	if err != nil {
		logger.Error().Err(err).Str("projectID", projectID).Msg("Invalid project ID format")
		return false, nil
	}

	project, err := s.db.Projects().GetProjectByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("projectID", projectID).Msg("Project not found")
			return false, nil
		}
		logger.Error().Err(err).Str("projectID", projectID).Msg("Failed to retrieve project")
		return false, err
	}

	if !project.IsMember(user.ID) {
		logger.Debug().Str("username", username).Str("projectID", projectID).Msg("User is not a member of the project")
		return false, nil
	}

	hasPermission := project.GetMember(user.ID).Role.HasPermissions(permission)
	logger.Debug().Str("username", username).Str("projectID", projectID).Interface("permission", permission).Bool("hasPermission", hasPermission).Msg("Checked project permission successfully")
	return hasPermission, nil
}

func (s *AuthzMiddlewareService) CheckMonitorPermissionByID(ctx context.Context, username, monitorID string, permission models.Permission) (bool, error) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameAuthzMiddleware, "CheckMonitorPermissionByID")

	parsedUUID, err := uuid.Parse(monitorID)
	if err != nil {
		logger.Error().Err(err).Str("monitorID", monitorID).Msg("Invalid monitor ID format")
		return false, nil
	}
	monitor, err := s.db.Monitors().GetMonitorByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("monitorID", monitorID).Msg("Monitor not found")
			return false, nil
		}
		logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to retrieve monitor")
		return false, err
	}

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", username).Msg("User not found")
			return false, nil
		}
		logger.Error().Err(err).Str("username", username).Msg("Failed to retrieve user")
		return false, err
	}

	project, err := s.db.Projects().GetProjectByID(ctx, monitor.ProjectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("projectID", monitor.ProjectID.String()).Msg("Project not found")
			return false, nil
		}
		logger.Error().Err(err).Str("projectID", monitor.ProjectID.String()).Msg("Failed to retrieve project")
		return false, err
	}

	if !project.IsMember(user.ID) {
		logger.Debug().Str("username", username).Str("monitorID", monitorID).Msg("User is not a member of the project")
		return false, nil
	}

	hasPermission := project.GetMember(user.ID).Role.HasPermissions(permission)
	logger.Debug().Str("username", username).Str("monitorID", monitorID).Interface("permission", permission).Bool("hasPermission", hasPermission).Msg("Checked monitor permission successfully")
	return hasPermission, nil
}
