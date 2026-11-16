package statuspages_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/statuspages"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/internal/testsupport"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/require"
)

func setupStatusPageIntegrationTest(t *testing.T) (context.Context, *statuspages.StatusPageService) {
	ctx, _, _ := testsupport.Setup(t)

	statsService := stats.NewMonitorStatsService(stats.MonitorStatsServiceDeps{DB: db.Get()})
	service := statuspages.NewStatusPageService(statuspages.StatusPageServiceDeps{
		DB:    db.Get(),
		Stats: &statsService,
	})

	return ctx, service
}

func insertTag(ctx context.Context, t *testing.T, name string) *tags.Tag {
	t.Helper()
	tag, err := tags.NewTagDAO(db.Get().Querier()).
		InsertTag(ctx, tags.Tag{Name: name, ColorHex: "#aabbcc", Source: util.ResourceSourceUI})
	require.NoError(t, err)
	return tag
}

func insertMonitorWithTags(ctx context.Context, t *testing.T, tagIDs ...uuid.UUID) *monitors.Monitor {
	t.Helper()
	monitor := testsupport.InsertTestMonitor(ctx, t)
	for _, tagID := range tagIDs {
		_, err := db.Get().Querier().ExecContext(
			ctx,
			`INSERT INTO monitor_tags (monitor_id, tag_id) VALUES ($1, $2)`,
			monitor.ID,
			tagID,
		)
		require.NoError(t, err)
	}
	return monitor
}

func insertResult(ctx context.Context, t *testing.T, monitorID uuid.UUID, status kind.MonitorStatus, latencyMs int64) {
	t.Helper()
	res := results.NewMonitorResult(monitorID, kind.HTTPConfigType, status, false, &latencyMs, nil)
	_, err := results.NewMonitorResultDAO(db.Get().Querier()).InsertMonitorResult(ctx, &res)
	require.NoError(t, err)
}

func newStatusPage(slug string, tagIDs ...uuid.UUID) statuspages.StatusPage {
	return statuspages.StatusPage{
		Slug:              slug,
		Title:             "Status " + slug,
		Description:       "Test status page",
		Enabled:           true,
		TimeWindowSeconds: 3600,
		TagIDs:            tagIDs,
	}
}

func createStatusPage(
	ctx context.Context,
	t *testing.T,
	service *statuspages.StatusPageService,
	page statuspages.StatusPage,
) *statuspages.StatusPage {
	t.Helper()
	created, svcErr := service.CreateStatusPage(ctx, page)
	require.Nil(t, svcErr)
	return created
}
