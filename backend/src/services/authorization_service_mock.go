package services

import (
	"context"

	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/mock"
)

// MockAuthorizationService is a mock implementation of IAuthorizer for testing.
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
