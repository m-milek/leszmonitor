package auditlog

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/require"
)

// setupAuditLogIntegrationTest initializes a temporary SQLite DB, sets up services, and registers a test user.
func setupAuditLogIntegrationTest(
	t *testing.T,
) (context.Context, AuditLogService, *users.UserService, *users.User) {
	ctx := context.Background()

	t.Setenv("JWT_SECRET", "test_secret_key_1234567890123456")
	t.Setenv("JWT_EXPIRY_HOURS", "24")

	dsn := "file:" + filepath.Join(t.TempDir(), "testdb.sqlite") + "?_pragma=foreign_keys(1)"

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

	registerPayload := &users.UserRegisterPayload{
		Username:        "integration_user",
		Password:        "Password123!",
		PasswordConfirm: "Password123!",
	}
	require.Nil(t, userService.RegisterUser(ctx, registerPayload))

	user, svcErr := userService.GetUserByUsername(ctx, "integration_user")
	require.Nil(t, svcErr)
	require.NotNil(t, user)

	ctx = auth.SetUserInContext(ctx, &auth.UserClaims{
		Username: user.Username,
	})

	auditLogService := NewAuditLogService(AuditLogServiceDeps{
		DB: realDB,
	})

	return ctx, auditLogService, userService, user
}
