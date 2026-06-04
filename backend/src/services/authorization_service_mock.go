package services

import (
	"context"

	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/mock"
)

// IProjectActionAuthorizer defines the interface for authorization operations.
// This interface allows for easy mocking in tests.
type IProjectActionAuthorizer interface {
	authorizeProjectAction(ctx context.Context, projectAuth *authorization.ProjectAuthorization, permissions ...models.Permission) (*models.Project, *ServiceError)
	isInstanceAdmin(ctx context.Context, username string) (bool, error)
}

// MockAuthorizationService is a mock implementation of IProjectActionAuthorizer for testing.
type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) authorizeProjectAction(ctx context.Context, projectAuth *authorization.ProjectAuthorization, permissions ...models.Permission) (*models.Project, *ServiceError) {
	args := m.Called(ctx, projectAuth, permissions)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*ServiceError)
	}
	return args.Get(0).(*models.Project), nil
}
