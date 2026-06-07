package services

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/db"
	"github.com/stretchr/testify/assert"
)

func setupTestAuditLogService() (context.Context, AuditLogService, *db.MockDB) {
	ctx := context.Background()
	mockDB := &db.MockDB{
		ProjectsRepo: new(db.MockProjectRepository),
	}
	db.Set(mockDB)

	authService := NewAuthorizationService(AuthorizationServiceDeps{DB: mockDB})

	return ctx, NewAuditLogService(AuditLogServiceDeps{
		DB:          mockDB,
		AuthService: authService,
	}), mockDB
}

func TestAuditLogService_GetEntries(t *testing.T) {
	t.Run("GetEntries success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestAuditLogService()
		// defer db.Set(nil)

		// TODO: Add mock setup and call the service

		assert.True(t, true) // Replace with actual assertions
	})
}

func TestAuditLogService_Record(t *testing.T) {
	t.Run("Record success", func(t *testing.T) {
		// ctx, svc, mockDB := setupTestAuditLogService()
		// defer db.Set(nil)

		// TODO: Add mock setup and call the service

		assert.True(t, true) // Replace with actual assertions
	})
}
