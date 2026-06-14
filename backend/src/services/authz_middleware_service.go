package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
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
	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	project, err := s.db.Projects().GetProjectBySlug(ctx, projectSlug)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	if !project.IsMember(user.ID) {
		return false, nil
	}

	return project.GetMember(user.ID).Role.HasPermissions(permission), nil
}

func (s *AuthzMiddlewareService) CheckProjectPermissionByID(ctx context.Context, username, projectID string, permission models.Permission) (bool, error) {
	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	parsedUUID, err := uuid.Parse(projectID)
	if err != nil {
		return false, nil
	}

	project, err := s.db.Projects().GetProjectByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	if !project.IsMember(user.ID) {
		return false, nil
	}

	return project.GetMember(user.ID).Role.HasPermissions(permission), nil
}

func (s *AuthzMiddlewareService) CheckMonitorPermissionByID(ctx context.Context, username, monitorID string, permission models.Permission) (bool, error) {
	parsedUUID, err := uuid.Parse(monitorID)
	if err != nil {
		return false, nil
	}
	monitor, err := s.db.Monitors().GetMonitorByID(ctx, parsedUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	project, err := s.db.Projects().GetProjectByID(ctx, monitor.ProjectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return false, nil
		}
		return false, err
	}

	if !project.IsMember(user.ID) {
		return false, nil
	}

	return project.GetMember(user.ID).Role.HasPermissions(permission), nil
}
