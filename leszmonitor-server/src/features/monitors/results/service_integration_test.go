package results_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_MonitorResultsService_GetLatest(t *testing.T) {
	t.Run("Successfully gets latest monitor result", func(t *testing.T) {
		ctx, service, _, _ := setupMonitorResultsIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		// Insert 2 results
		res1 := results.NewMonitorResult(monitor.ID, kind.HTTPConfigType, kind.MonitorStatusUp, false, 100, nil)
		res1.CreatedAt = time.Now().UTC().Add(-10 * time.Minute)
		_, err := results.NewMonitorResultDAO(db.Get().Querier()).InsertMonitorResult(ctx, &res1)
		require.NoError(t, err)

		res2 := results.NewMonitorResult(
			monitor.ID,
			kind.HTTPConfigType,
			kind.MonitorStatusDown,
			false,
			200,
			nil,
		)
		res2.CreatedAt = time.Now().UTC().Add(-5 * time.Minute)
		_, err = results.NewMonitorResultDAO(db.Get().Querier()).InsertMonitorResult(ctx, &res2)
		require.NoError(t, err)

		latest, svcErr := service.GetLatestMonitorResultByMonitorID(ctx, monitor.ID.String())
		require.Nil(t, svcErr)
		require.NotNil(t, latest)
		assert.Equal(t, res2.ID, latest.GetID())
		assert.Equal(t, kind.MonitorStatusDown, latest.GetStatus())
	})

	t.Run("Fails with 404 when no results exist", func(t *testing.T) {
		ctx, service, _, _ := setupMonitorResultsIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		latest, svcErr := service.GetLatestMonitorResultByMonitorID(ctx, monitor.ID.String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
		assert.Nil(t, latest)
	})
}

func TestIntegration_MonitorResultsService_GetAll(t *testing.T) {
	t.Run("Successfully gets paginated monitor results", func(t *testing.T) {
		ctx, service, _, _ := setupMonitorResultsIntegrationTest(t)

		monitor := insertTestMonitor(ctx, t)

		// Insert 3 results
		for i := range 3 {
			res := results.NewMonitorResult(
				monitor.ID,
				kind.HTTPConfigType,
				kind.MonitorStatusUp,
				false,
				int64(100+i),
				nil,
			)
			_, err := results.NewMonitorResultDAO(db.Get().Querier()).InsertMonitorResult(ctx, &res)
			require.NoError(t, err)
		}

		pag := &util.Pagination{Page: 1, PerPage: 10}
		results, svcErr := service.GetMonitorResultsByMonitorID(ctx, monitor.ID.String(), pag)
		require.Nil(t, svcErr)
		require.Len(t, results, 3)
	})
}
