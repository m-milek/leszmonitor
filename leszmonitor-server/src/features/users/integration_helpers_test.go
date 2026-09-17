package users

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/require"
)

// setupIntegrationTest initializes a temporary SQLite DB, sets up the user service, and registers a test user.
func setupIntegrationTest(t *testing.T) (context.Context, *UserService, *User) {
	ctx := context.Background()

	t.Setenv("JWT_SECRET", "test_secret_key_1234567890123456")
	t.Setenv("JWT_EXPIRY_HOURS", "24")

	// Use a temporary directory for the sqlite file
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "testdb.sqlite")

	// DSN format for SQLite
	dsn := "file:" + dbPath + "?_pragma=foreign_keys(1)"

	database, err := db.New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(database.Close)

	userService := NewUserService(UserServiceDeps{
		DB: database,
	})

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

	ctx = auth.SetUserInContext(ctx, &auth.UserClaims{
		Username: user.Username,
	})

	return ctx, userService, user
}
