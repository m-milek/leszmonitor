package services

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/db"
	"github.com/stretchr/testify/assert"
)

func setupTestMonitorService() (context.Context, *MonitorService, *db.MockDB) {
	ctx := context.Background()
	mockDB := &db.MockDB{
		ProjectsRepo: new(db.MockProjectRepository),
	}
	db.Set(mockDB)

	authService := NewAuthorizationService(AuthorizationServiceDeps{DB: mockDB})

	return ctx, NewMonitorService(MonitorServiceDeps{
		DB:   mockDB,
		Auth: authService,
	}), mockDB
}

func TestMonitorService_CreateMonitor(t *testing.T) {
	t.Run("CreateMonitor success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestMonitorService()
		// defer db.Set(nil)

		// TODO: Add mock setup and call the service

		assert.True(t, true) // Replace with actual assertions
	})
}

func TestMonitorService_GetMonitorByID(t *testing.T) {
	t.Run("GetMonitorByID success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestMonitorService()
		// defer db.Set(nil)

		// TODO: Add mock setup and call the service

		assert.True(t, true) // Replace with actual assertions
	})
}
