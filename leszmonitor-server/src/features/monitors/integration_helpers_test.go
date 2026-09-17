package monitors

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTest initializes a temporary SQLite DB, sets up services, and registers a test user.
func setupIntegrationTest(t *testing.T) (context.Context, *users.UserService, *users.User) {
	ctx := context.Background()

	t.Setenv("JWT_SECRET", "test_secret_key_1234567890123456")
	t.Setenv("JWT_EXPIRY_HOURS", "24")

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

	userService := users.NewUserService(users.UserServiceDeps{
		DB: realDB,
	})

	// Setup Phase: Create a real user in the DB
	registerPayload := &users.UserRegisterPayload{
		Username:        "integration_user",
		Password:        "Password123!",
		PasswordConfirm: "Password123!",
	}
	svcErr := userService.RegisterUser(ctx, registerPayload)
	require.Nil(t, svcErr)

	user, svcErr := userService.GetUserByUsername(ctx, "integration_user")
	require.Nil(t, svcErr)
	require.NotNil(t, user)

	ctx = auth.SetUserInContext(ctx, &auth.UserClaims{
		Username: user.Username,
	})

	return ctx, userService, user
}

func setupMonitorIntegrationTest(
	t *testing.T,
) (context.Context, *MonitorService, *users.UserService, *users.User) {
	ctx, userService, user := setupIntegrationTest(t)

	monitorService := NewMonitorService(MonitorServiceDeps{
		DB: db.Get(),
	})

	return ctx, monitorService, userService, user
}

func setupMonitorResultsIntegrationTest(
	t *testing.T,
) (context.Context, *results.MonitorResultsService, *users.UserService, *users.User) {
	ctx, userService, user := setupIntegrationTest(t)

	service := results.NewMonitorResultsService(results.MonitorResultsServiceDeps{
		DB: db.Get(),
	})

	return ctx, service, userService, user
}

func setupMonitorStatsIntegrationTest(
	t *testing.T,
) (context.Context, *stats.MonitorStatsService, *users.UserService, *users.User) {
	ctx, userService, user := setupIntegrationTest(t)

	monitorStatsService := stats.NewMonitorStatsService(stats.MonitorStatsServiceDeps{
		DB: db.Get(),
	})

	return ctx, &monitorStatsService, userService, user
}

// insertTestMonitor is a helper to directly insert a monitor and return it.
func insertTestMonitor(t *testing.T, ctx context.Context) *Monitor {
	payload := Monitor{
		Name:        "Test Monitor " + uuid.New().String(),
		Description: "Testing monitor results",
		Interval:    60,
		Type:        kind.HTTPConfigType,
		ProbeConfig: "{}",
	}
	payload.GenerateSlug()

	owner, err := users.NewUserDAO(db.Get().Querier()).GetUserByUsername(ctx, "integration_user")
	require.NoError(t, err)
	monitor := InitializeFromPayload(payload, owner.ID)

	inserted, dbErr := NewMonitorDAO(db.Get().Querier()).InsertMonitor(ctx, *monitor)
	require.NoError(t, dbErr)
	return inserted
}
