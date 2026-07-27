package probes

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitors"
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

	user, err := models.NewUser("testuser", "testpassword123")
	require.NoError(t, err)
	insertedUser, err := realDB.Users().InsertUser(ctx, user)
	require.NoError(t, err)

	project, err := models.NewProject("Test Project", "Desc", insertedUser.ID)
	require.NoError(t, err)
	err = realDB.Projects().InsertProject(ctx, project)
	require.NoError(t, err)

	payload := monitors.Monitor{
		Name:        "Test Monitor " + uuid.New().String(),
		Description: "Testing monitor results",
		Interval:    1,
		Type:        consts.HTTPConfigType,
		ProbeConfig: `{"method": "GET", "url": "http://localhost:8080", "expectedStatusCodes": [200]}`,
		State:       monitors.MonitorStateActive,
	}
	payload.GenerateSlug()
	monitor := monitors.InitializeFromPayload(payload, project.ID)

	insertedMonitor, err := realDB.Monitors().InsertMonitor(ctx, *monitor)
	require.NoError(t, err)

	return ctx, realDB, insertedMonitor
}
