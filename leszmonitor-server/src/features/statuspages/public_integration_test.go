package statuspages_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/statuspages"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_StatusPageService_GetPublicMonitors(t *testing.T) {
	t.Run("Returns monitors with page tags and stats", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		prod := insertTag(ctx, t, "prod")
		other := insertTag(ctx, t, "other")

		monitor := insertMonitorWithTags(ctx, t, prod.ID, other.ID)
		insertMonitorWithTags(ctx, t, other.ID)
		insertResult(ctx, t, monitor.ID, kind.MonitorStatusUp, 100)
		insertResult(ctx, t, monitor.ID, kind.MonitorStatusDown, 300)

		page := createStatusPage(ctx, t, service, newStatusPage("public", prod.ID))

		now := time.Now().UTC()
		publicMonitors, svcErr := service.GetPublicMonitors(ctx, page.ID, now.Add(-time.Hour), now.Add(time.Minute))
		require.Nil(t, svcErr)
		require.Len(t, publicMonitors, 1)

		pm := publicMonitors[0]
		assert.Equal(t, monitor.ID, pm.ID)
		require.Len(t, pm.Tags, 1, "only tags on the status page are returned")
		assert.Equal(t, prod.ID, pm.Tags[0].ID)
		assert.Len(t, pm.Results, 2)
		assert.InDelta(t, 50.0, pm.Stats.Uptime.StatusToPercentage[kind.MonitorStatusUp], 0.001)
		assert.InDelta(t, 200.0, pm.Stats.Latency.Avg, 0.001)
	})

	t.Run("Reports paused monitors as paused", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		monitor := insertMonitorWithTags(ctx, t, tag.ID)
		insertResult(ctx, t, monitor.ID, kind.MonitorStatusUp, 100)
		_, err := db.Get().Querier().ExecContext(ctx, `UPDATE monitors SET run_state = 'paused' WHERE id = $1`, monitor.ID)
		require.NoError(t, err)

		page := createStatusPage(ctx, t, service, newStatusPage("public", tag.ID))

		now := time.Now().UTC()
		publicMonitors, svcErr := service.GetPublicMonitors(ctx, page.ID, now.Add(-time.Hour), now)
		require.Nil(t, svcErr)
		require.Len(t, publicMonitors, 1)
		assert.Equal(t, kind.MonitorStatusPaused, publicMonitors[0].Status)
	})

	t.Run("Reports unknown status without results", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		insertMonitorWithTags(ctx, t, tag.ID)

		page := createStatusPage(ctx, t, service, newStatusPage("public", tag.ID))

		now := time.Now().UTC()
		publicMonitors, svcErr := service.GetPublicMonitors(ctx, page.ID, now.Add(-time.Hour), now)
		require.Nil(t, svcErr)
		require.Len(t, publicMonitors, 1)
		assert.Equal(t, kind.MonitorStatusUnknown, publicMonitors[0].Status)
		assert.Empty(t, publicMonitors[0].Results)
	})
}

func TestIntegration_StatusPageAPIController_PublicAccess(t *testing.T) {
	denyReader := func(http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}

	request := func(t *testing.T, service *statuspages.StatusPageService, slug string) int {
		t.Helper()
		controller := statuspages.NewStatusPageAPIController(service)
		req := httptest.NewRequest(http.MethodGet, "/api/v1/public/status-pages/"+slug, nil)
		req.SetPathValue("slug", slug)
		rec := httptest.NewRecorder()
		controller.GetPublicStatusPageHandler(denyReader)(rec, req)
		return rec.Code
	}

	tests := []struct {
		name         string
		modify       func(p *statuspages.StatusPage)
		expectedCode int
	}{
		{"Public page is accessible", func(*statuspages.StatusPage) {}, http.StatusOK},
		{"Page requiring auth needs a reader", func(p *statuspages.StatusPage) { p.RequireAuth = true }, http.StatusUnauthorized},
		{"Disabled page is not found", func(p *statuspages.StatusPage) { p.Enabled = false }, http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, service := setupStatusPageIntegrationTest(t)
			tag := insertTag(ctx, t, "prod")
			page := newStatusPage("public", tag.ID)
			tt.modify(&page)
			createStatusPage(ctx, t, service, page)

			assert.Equal(t, tt.expectedCode, request(t, service, "public"))
		})
	}

	t.Run("Unknown slug is not found", func(t *testing.T) {
		_, service := setupStatusPageIntegrationTest(t)
		assert.Equal(t, http.StatusNotFound, request(t, service, "missing"))
	})
}
