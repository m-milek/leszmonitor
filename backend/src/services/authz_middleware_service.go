package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
)

type IAuthzMiddlewareService interface {
	CheckProjectPermissionBySlug(ctx context.Context, username, projectSlug string, permission models.Permission) (bool, error)
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
