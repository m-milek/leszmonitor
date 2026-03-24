package services

import (
	"context"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/mock"
)

// IOrgActionAuthorizer defines the interface for authorization operations.
// This interface allows for easy mocking in tests.
type IOrgActionAuthorizer interface {
	authorizeOrgAction(ctx context.Context, orgAuth *middleware.OrgAuth, permissions ...models.Permission) (*models.Org, *ServiceError)
}

// MockAuthorizationService is a mock implementation of IOrgActionAuthorizer for testing.
type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) authorizeOrgAction(ctx context.Context, orgAuth *middleware.OrgAuth, permissions ...models.Permission) (*models.Org, *ServiceError) {
	args := m.Called(ctx, orgAuth, permissions)
	if args.Get(0) == nil {
		return nil, args.Get(1).(*ServiceError)
	}
	return args.Get(0).(*models.Org), nil
}
