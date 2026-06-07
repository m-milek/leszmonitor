package services

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/db"
	"github.com/stretchr/testify/assert"
)

func setupTestMonitorResultsService() (context.Context, *MonitorResultsService, *db.MockDB) {
	ctx := context.Background()
	mockDB := &db.MockDB{
		ProjectsRepo: new(db.MockProjectRepository),
	}
	db.Set(mockDB)

	authService := NewAuthorizationService(AuthorizationServiceDeps{DB: mockDB})

	return ctx, NewMonitorResultsService(MonitorResultsServiceDeps{
		DB:   mockDB,
		Auth: authService,
	}), mockDB
}

func TestMonitorResultsService_GetLatestMonitorResultByMonitorID(t *testing.T) {
	t.Run("GetLatestMonitorResultByMonitorID success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestMonitorResultsService()
		// defer db.Set(nil)

		// TODO: Add mock setup and call the service

		assert.True(t, true) // Replace with actual assertions
	})
}
