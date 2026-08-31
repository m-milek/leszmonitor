package services

import (
	"net/http"
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/models/shared"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MonitorStatsService_GetLatencyStatsByMonitorID(t *testing.T) {
	t.Run("Returns correct avg/min/max latency for results in range", func(t *testing.T) {
		ctx, service, projectService, _, owner := setupMonitorStatsIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Project 1"})
		monitor := insertTestMonitor(t, ctx, project.ID)

		now := time.Now().UTC()

		// Insert 3 results with known latencies: 100, 200, 300 → avg=200, min=100, max=300
		for _, latency := range []int64{100, 200, 300} {
			res := monitorresult.NewMonitorResult(monitor.ID, consts.HTTPConfigType, shared.MonitorStatusUp, false, latency, "", nil)
			// Store created_at in UTC RFC3339 so SQLite string comparisons work correctly with the DAO's UTC-formatted query bounds
			res.CreatedAt = now.Add(-30 * time.Minute).Format(time.RFC3339)
			_, err := db.Get().MonitorResults().InsertMonitorResult(ctx, &res)
			require.NoError(t, err)
		}

		from := now.Add(-1 * time.Hour)
		to := now.Add(1 * time.Hour)

		stats, svcErr := service.GetLatencyStatsByMonitorID(ctx, monitor.ID.String(), from, to)
		require.Nil(t, svcErr)
		assert.InDelta(t, 200.0, stats.Avg, 0.001)
		assert.InDelta(t, 100.0, stats.Min, 0.001)
		assert.InDelta(t, 300.0, stats.Max, 0.001)
	})

	t.Run("Returns 404 when no results exist for monitor", func(t *testing.T) {
		ctx, service, projectService, _, owner := setupMonitorStatsIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Project 1"})
		monitor := insertTestMonitor(t, ctx, project.ID)

		from := time.Now().UTC().Add(-1 * time.Hour)
		to := time.Now().UTC().Add(1 * time.Hour)

		stats, svcErr := service.GetLatencyStatsByMonitorID(ctx, monitor.ID.String(), from, to)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
		assert.Equal(t, 0.0, stats.Avg)
		assert.Equal(t, 0.0, stats.Min)
		assert.Equal(t, 0.0, stats.Max)
	})

	t.Run("Only includes results within the specified time range", func(t *testing.T) {
		ctx, service, projectService, _, owner := setupMonitorStatsIntegrationTest(t)

		project, _ := projectService.CreateProject(ctx, owner.Username, CreateProjectPayload{Name: "Project 1"})
		monitor := insertTestMonitor(t, ctx, project.ID)

		now := time.Now().UTC()

		// Result inside the query range: latency 500
		resIn := monitorresult.NewMonitorResult(monitor.ID, consts.HTTPConfigType, shared.MonitorStatusUp, false, 500, "", nil)
		resIn.CreatedAt = now.Add(-30 * time.Minute).Format(time.RFC3339)
		_, err := db.Get().MonitorResults().InsertMonitorResult(ctx, &resIn)
		require.NoError(t, err)

		// Result outside the query range (too old): latency 9999, should be excluded
		resOut := monitorresult.NewMonitorResult(monitor.ID, consts.HTTPConfigType, shared.MonitorStatusUp, false, 9999, "", nil)
		resOut.CreatedAt = now.Add(-3 * time.Hour).Format(time.RFC3339)
		_, err = db.Get().MonitorResults().InsertMonitorResult(ctx, &resOut)
		require.NoError(t, err)

		from := now.Add(-1 * time.Hour)
		to := now.Add(1 * time.Hour)

		stats, svcErr := service.GetLatencyStatsByMonitorID(ctx, monitor.ID.String(), from, to)
		require.Nil(t, svcErr)
		assert.InDelta(t, 500.0, stats.Avg, 0.001)
		assert.InDelta(t, 500.0, stats.Min, 0.001)
		assert.InDelta(t, 500.0, stats.Max, 0.001)
	})
}
