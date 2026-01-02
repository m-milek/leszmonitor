package services

import (
	"context"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/mock"
)

// IAuthorizationService defines the interface for authorization operations.
// This interface allows for easy mocking in tests.
type IAuthorizationService interface {
	authorizeTeamAction(ctx context.Context, teamAuth *middleware.TeamAuth, permissions ...models.Permission) (*models.Team, *ServiceError)
}

// MockAuthorizationService is a mock implementation of IAuthorizationService for testing.
type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) authorizeTeamAction(ctx context.Context, teamAuth *middleware.TeamAuth, permissions ...models.Permission) (*models.Team, *ServiceError) {
	args := m.Called(ctx, teamAuth, permissions)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*ServiceError)
	}
	return args.Get(0).(*models.Team), nil
}
