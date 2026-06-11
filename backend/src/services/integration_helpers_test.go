package services

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTest initializes a temporary SQLite DB, sets up services, and registers a test user.
func setupIntegrationTest(t *testing.T) (context.Context, *ProjectService, *UserService, *models.User) {
	ctx := context.Background()

	os.Setenv("JWT_SECRET", "test_secret_key_1234567890123456")
	os.Setenv("JWT_EXPIRY_HOURS", "24")

	// Use a temporary directory for the sqlite file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "testdb.sqlite")

	// DSN format for SQLite
	dsn := "file:" + dbPath + "?_pragma=foreign_keys(1)"

	realDB, err := db.New(ctx, dsn)
	require.NoError(t, err)

	// Set globally because parts of the system (like authorization) might still rely on db.Get()
	db.Set(realDB)

	t.Cleanup(func() {
		realDB.Close()
		db.Set(nil)
	})

	userService := NewUserService(UserServiceDeps{
		DB: realDB,
	})

	projectService := NewProjectService(ProjectServiceDeps{
		DB: realDB,

		UserService: userService,
	})

	userService.projectService = projectService

	// Setup Phase: Create a real user in the DB
	registerPayload := &UserRegisterPayload{
		Username:        "integration_user",
		Password:        "Password123!",
		PasswordConfirm: "Password123!",
	}
	svcErr := userService.RegisterUser(ctx, registerPayload)
	require.Nil(t, svcErr)

	user, svcErr := userService.GetUserByUsername(ctx, "integration_user")
	require.Nil(t, svcErr)
	require.NotNil(t, user)

	return ctx, projectService, userService, user
}

func setupAuditLogIntegrationTest(t *testing.T) (context.Context, AuditLogService, *ProjectService, *UserService, *models.User) {
	ctx, projectService, userService, user := setupIntegrationTest(t)

	auditLogService := NewAuditLogService(AuditLogServiceDeps{
		DB:          db.Get(),
		AuthService: NewAuthorizationService(AuthorizationServiceDeps{DB: db.Get()}),
	})

	return ctx, auditLogService, projectService, userService, user
}

func setupMonitorResultsIntegrationTest(t *testing.T) (context.Context, *MonitorResultsService, *ProjectService, *UserService, *models.User) {
	ctx, projectService, userService, user := setupIntegrationTest(t)

	service := NewMonitorResultsService(MonitorResultsServiceDeps{
		DB:   db.Get(),
		Auth: NewAuthorizationService(AuthorizationServiceDeps{DB: db.Get()}),
	})

	return ctx, service, projectService, userService, user
}

func setupMonitorIntegrationTest(t *testing.T) (context.Context, *MonitorService, *ProjectService, *UserService, *models.User) {
	ctx, projectService, userService, user := setupIntegrationTest(t)

	monitorService := NewMonitorService(MonitorServiceDeps{
		DB:   db.Get(),
		Auth: NewAuthorizationService(AuthorizationServiceDeps{DB: db.Get()}),
	})

	return ctx, monitorService, projectService, userService, user
}

// insertTestMonitor is a helper to directly insert a monitor and return it
func insertTestMonitor(t *testing.T, ctx context.Context, projectID uuid.UUID) *monitors.Monitor {
	payload := monitors.Monitor{
		Name:        "Test Monitor " + uuid.New().String(),
		Description: "Testing monitor results",
		Interval:    60,
		Type:        consts.HttpConfigType,
		ProbeConfig: "{}",
	}
	payload.GenerateSlug()
	monitor := monitors.InitializeFromPayload(payload, projectID)

	inserted, dbErr := db.Get().Monitors().InsertMonitor(ctx, *monitor)
	require.NoError(t, dbErr)
	return inserted
}
