package workers

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/users"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) (context.Context, db.DB) {
	ctx := context.Background()

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "testdb.sqlite")
	dsn := "file:" + dbPath + "?_pragma=foreign_keys(1)"

	realDB, err := db.New(ctx, dsn)
	require.NoError(t, err)

	t.Cleanup(func() {
		realDB.Close()
	})

	return ctx, realDB
}

func setupFullDB(t *testing.T) (context.Context, db.DB, *monitors.Monitor) {
	ctx, realDB := setupDB(t)

	user, err := users.NewUser("testuser", "testpassword123")
	require.NoError(t, err)
	insertedUser, err := users.NewUserDAO(realDB.Querier()).InsertUser(ctx, user)
	require.NoError(t, err)

	payload := monitors.Monitor{
		Name:        "Test Monitor " + uuid.New().String(),
		Description: "Testing monitor results",
		Interval:    1,
		Type:        kind.HTTPConfigType,
		ProbeConfig: `{"method": "GET", "url": "http://localhost:8080", "expectedStatusCodes": [200]}`,
		RunState:    monitors.MonitorStateActive,
	}
	payload.GenerateSlug()
	monitor := monitors.InitializeFromPayload(payload, insertedUser.ID)

	insertedMonitor, err := monitors.NewMonitorDAO(realDB.Querier()).InsertMonitor(ctx, *monitor)
	require.NoError(t, err)

	return ctx, realDB, insertedMonitor
}
