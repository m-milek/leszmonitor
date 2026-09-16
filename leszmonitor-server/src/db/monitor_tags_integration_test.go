package db

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/stretchr/testify/require"
)

func setupMonitorTagsTest(t *testing.T) (context.Context, *Client) {
	t.Helper()
	ctx := context.Background()

	dsn := "file:" + filepath.Join(t.TempDir(), "testdb.sqlite") + "?_pragma=foreign_keys(1)"
	client, err := New(ctx, dsn)
	require.NoError(t, err)
	t.Cleanup(client.Close)

	return ctx, client
}

func insertTestTag(ctx context.Context, t *testing.T, client *Client, name string) uuid.UUID {
	t.Helper()
	tag, err := client.Tags().InsertTag(ctx, models.Tag{Name: name, ColorHex: "#aabbcc"})
	require.NoError(t, err)
	return tag.ID
}

func testMonitor(tagIDs []uuid.UUID) monitors.Monitor {
	return monitors.Monitor{
		ID:                     uuid.New(),
		Slug:                   "test-monitor",
		Name:                   "Test Monitor",
		Interval:               60,
		Type:                   consts.HTTPConfigType,
		ProbeConfig:            `{"method":"GET","url":"http://example.com"}`,
		ResultRetentionSeconds: 3600,
		RunState:               monitors.MonitorStateActive,
		OwnerID:                uuid.New(),
		TagIDs:                 tagIDs,
	}
}

func TestInsertMonitorPersistsTagIDs(t *testing.T) {
	ctx, client := setupMonitorTagsTest(t)
	first := insertTestTag(ctx, t, client, "prod")
	second := insertTestTag(ctx, t, client, "critical")

	created, err := client.Monitors().
		InsertMonitor(ctx, testMonitor([]uuid.UUID{first, second, first}))
	require.NoError(t, err)
	require.ElementsMatch(t, []uuid.UUID{first, second}, created.TagIDs)

	bySlug, err := client.Monitors().GetMonitorBySlug(ctx, "test-monitor")
	require.NoError(t, err)
	require.ElementsMatch(t, []uuid.UUID{first, second}, bySlug.TagIDs)

	all, err := client.Monitors().GetAllMonitors(ctx)
	require.NoError(t, err)
	require.Len(t, all, 1)
	require.ElementsMatch(t, []uuid.UUID{first, second}, all[0].TagIDs)
}

func TestInsertMonitorWithoutTagsReturnsEmptyList(t *testing.T) {
	ctx, client := setupMonitorTagsTest(t)

	created, err := client.Monitors().InsertMonitor(ctx, testMonitor(nil))
	require.NoError(t, err)
	require.NotNil(t, created.TagIDs)
	require.Empty(t, created.TagIDs)
}

func TestUpdateMonitorReplacesTagIDs(t *testing.T) {
	ctx, client := setupMonitorTagsTest(t)
	first := insertTestTag(ctx, t, client, "prod")
	second := insertTestTag(ctx, t, client, "critical")

	created, err := client.Monitors().InsertMonitor(ctx, testMonitor([]uuid.UUID{first}))
	require.NoError(t, err)

	updated := *created
	updated.TagIDs = []uuid.UUID{second}
	_, err = client.Monitors().UpdateMonitor(ctx, updated)
	require.NoError(t, err)

	reloaded, err := client.Monitors().GetMonitorByID(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, []uuid.UUID{second}, reloaded.TagIDs)

	updated.TagIDs = nil
	_, err = client.Monitors().UpdateMonitor(ctx, updated)
	require.NoError(t, err)

	reloaded, err = client.Monitors().GetMonitorByID(ctx, created.ID)
	require.NoError(t, err)
	require.Empty(t, reloaded.TagIDs)
}
