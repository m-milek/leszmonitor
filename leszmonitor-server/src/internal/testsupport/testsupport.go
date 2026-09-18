// Package testsupport provides shared fixtures for integration tests.
package testsupport

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/require"
)

// Setup initializes a temporary SQLite DB, sets up services, and registers a test user.
func Setup(t *testing.T) (context.Context, *users.UserService, *users.User) {
	t.Helper()
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

// testMonitorIntervalSeconds is the check interval used by monitors created in tests.
const testMonitorIntervalSeconds = 60

// InsertTestMonitor inserts a monitor owned by the test user and returns it.
func InsertTestMonitor(ctx context.Context, t *testing.T) *monitors.Monitor {
	t.Helper()
	payload := monitors.Monitor{
		Name:        "Test Monitor " + uuid.New().String(),
		Description: "Testing monitor results",
		Interval:    testMonitorIntervalSeconds,
		Type:        kind.HTTPConfigType,
		ProbeConfig: "{}",
	}
	payload.GenerateSlug()

	owner, err := users.NewUserDAO(db.Get().Querier()).GetUserByUsername(ctx, "integration_user")
	require.NoError(t, err)
	monitor := monitors.InitializeFromPayload(payload, owner.ID)

	inserted, dbErr := monitors.NewMonitorDAO(db.Get().Querier()).InsertMonitor(ctx, *monitor)
	require.NoError(t, dbErr)
	return inserted
}
