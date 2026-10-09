package monitors_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/features/users"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func httpConfigEntry(name string) map[string]any {
	return map[string]any{
		"name":     name,
		"type":     "http",
		"interval": 60,
		"probeConfig": map[string]any{
			"method":              "GET",
			"url":                 "https://example.com",
			"expectedStatusCodes": []any{200},
		},
	}
}

func pushConfigEntry(name string) map[string]any {
	return map[string]any{
		"name":        name,
		"type":        "push",
		"interval":    60,
		"probeConfig": map[string]any{"gracePeriodSeconds": 10},
	}
}

func syncConfigMonitors(ctx context.Context, owner *users.User, entries map[string]map[string]any) error {
	return monitors.SynchronizeConfigBasedMonitors(ctx, db.Get(), entries, owner.ID)
}

func getMonitorBySlug(ctx context.Context, t *testing.T, slug string) *monitors.Monitor {
	t.Helper()
	monitor, err := monitors.NewMonitorDAO(db.Get().Querier()).GetMonitorBySlug(ctx, slug)
	require.NoError(t, err)
	return monitor
}

func getAllMonitors(ctx context.Context, t *testing.T) []monitors.Monitor {
	t.Helper()
	all, err := monitors.NewMonitorDAO(db.Get().Querier()).GetAllMonitors(ctx)
	require.NoError(t, err)
	return all
}

func TestIntegration_SynchronizeConfigBasedMonitors(t *testing.T) {
	t.Run("Inserts monitors from config", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)

		withRetention := pushConfigEntry("Backup job")
		withRetention["resultRetentionSeconds"] = 1000

		err := syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health": httpConfigEntry("API health"),
			"backup-job": withRetention,
		})
		require.NoError(t, err)

		api := getMonitorBySlug(ctx, t, "api-health")
		assert.Equal(t, "API health", api.Name)
		assert.Equal(t, kind.HTTPConfigType, api.Type)
		assert.Equal(t, monitors.MonitorSourceConfig, api.Source)
		assert.Equal(t, monitors.MonitorRunStateActive, api.RunState)
		assert.Equal(t, owner.ID, api.OwnerID)
		assert.Equal(t, 172800, api.ResultRetentionSeconds)

		backup := getMonitorBySlug(ctx, t, "backup-job")
		assert.Equal(t, 1000, backup.ResultRetentionSeconds)
	})

	t.Run("Derives the same ID on every run", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)
		entries := map[string]map[string]any{"api-health": httpConfigEntry("API health")}

		require.NoError(t, syncConfigMonitors(ctx, owner, entries))
		first := getMonitorBySlug(ctx, t, "api-health")

		require.NoError(t, syncConfigMonitors(ctx, owner, entries))
		second := getMonitorBySlug(ctx, t, "api-health")

		assert.Equal(t, first.ID, second.ID)
	})

	t.Run("Updates a changed monitor and keeps its results", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)

		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health": httpConfigEntry("API health"),
		}))
		before := getMonitorBySlug(ctx, t, "api-health")

		result := results.NewMonitorResult(before.ID, kind.HTTPConfigType, kind.MonitorStatusUp, false, nil, nil)
		_, err := results.NewMonitorResultDAO(db.Get().Querier()).InsertMonitorResult(ctx, &result)
		require.NoError(t, err)

		renamed := httpConfigEntry("API health renamed")
		renamed["interval"] = 120
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{"api-health": renamed}))

		after := getMonitorBySlug(ctx, t, "api-health")
		assert.Equal(t, before.ID, after.ID)
		assert.Equal(t, "API health renamed", after.Name)
		assert.Equal(t, 120, after.Interval)

		latest, err := results.NewMonitorResultDAO(db.Get().Querier()).GetLatestMonitorResultByMonitorID(ctx, after.ID)
		require.NoError(t, err)
		assert.Equal(t, result.ID, latest.GetID())
	})

	t.Run("Replaces the monitor when its type changes", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)

		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"job": httpConfigEntry("Job"),
		}))
		before := getMonitorBySlug(ctx, t, "job")

		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"job": pushConfigEntry("Job"),
		}))
		after := getMonitorBySlug(ctx, t, "job")

		assert.NotEqual(t, before.ID, after.ID)
		assert.Equal(t, kind.PushConfigType, after.Type)
		assert.Len(t, getAllMonitors(ctx, t), 1)
	})

	t.Run("Deletes monitors missing from config and keeps UI monitors", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)
		uiMonitor := insertTestMonitor(ctx, t)

		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health": httpConfigEntry("API health"),
			"backup-job": pushConfigEntry("Backup job"),
		}))
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health": httpConfigEntry("API health"),
		}))

		_, err := monitors.NewMonitorDAO(db.Get().Querier()).GetMonitorBySlug(ctx, "backup-job")
		require.ErrorIs(t, err, db.ErrNotFound)
		assert.Len(t, getAllMonitors(ctx, t), 2)
		assert.Equal(t, monitors.MonitorSourceUI, getMonitorBySlug(ctx, t, uiMonitor.Slug).Source)
	})

	t.Run("Deletes all config monitors when config has none", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)
		uiMonitor := insertTestMonitor(ctx, t)

		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health": httpConfigEntry("API health"),
		}))
		require.NoError(t, syncConfigMonitors(ctx, owner, nil))

		all := getAllMonitors(ctx, t)
		require.Len(t, all, 1)
		assert.Equal(t, uiMonitor.ID, all[0].ID)
	})

	t.Run("Fails and writes nothing when slug is used by a UI monitor", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)
		uiMonitor := insertTestMonitor(ctx, t)

		err := syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health":   httpConfigEntry("API health"),
			uiMonitor.Slug: httpConfigEntry("Clash"),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), uiMonitor.Slug)

		all := getAllMonitors(ctx, t)
		require.Len(t, all, 1)
		assert.Equal(t, uiMonitor.ID, all[0].ID)
	})

	t.Run("Assigns config and UI tags by name", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)
		uiTagID := insertTestTag(ctx, t, db.Get(), "staging")
		require.NoError(t, tags.SynchronizeConfigBasedTags(ctx, db.Get(), map[string]map[string]any{
			"production": {"colorHex": "#aabbcc"},
		}))
		configTag := uuid.NewSHA1(uuid.MustParse("abadcafe-deaf-dead-beef-cafebabefeed"), []byte("production"))

		entry := httpConfigEntry("API health")
		entry["tags"] = []any{"production", "staging"}
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{"api-health": entry}))

		assert.ElementsMatch(t, []uuid.UUID{configTag, uiTagID}, getMonitorBySlug(ctx, t, "api-health").TagIDs)
	})

	t.Run("Replaces tags on every run", func(t *testing.T) {
		ctx, _, _, owner := setupMonitorIntegrationTest(t)
		stagingID := insertTestTag(ctx, t, db.Get(), "staging")
		insertTestTag(ctx, t, db.Get(), "critical")

		entry := httpConfigEntry("API health")
		entry["tags"] = []any{"critical"}
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{"api-health": entry}))

		entry["tags"] = []any{"staging"}
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{"api-health": entry}))
		assert.Equal(t, []uuid.UUID{stagingID}, getMonitorBySlug(ctx, t, "api-health").TagIDs)

		delete(entry, "tags")
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{"api-health": entry}))
		assert.Empty(t, getMonitorBySlug(ctx, t, "api-health").TagIDs)
	})

	t.Run("Rejects invalid entries without writing anything", func(t *testing.T) {
		withForbiddenKey := httpConfigEntry("API health")
		withForbiddenKey["runState"] = "paused"

		withUnknownField := httpConfigEntry("API health")
		withUnknownField["intervall"] = 60

		withUnknownProbeField := httpConfigEntry("API health")
		withUnknownProbeField["probeConfig"].(map[string]any)["expectedStatusCode"] = 200

		withInvalidProbeConfig := httpConfigEntry("API health")
		withInvalidProbeConfig["probeConfig"].(map[string]any)["url"] = ""

		withoutType := httpConfigEntry("API health")
		delete(withoutType, "type")

		withUnknownTag := httpConfigEntry("API health")
		withUnknownTag["tags"] = []any{"missing"}

		withInvalidTags := httpConfigEntry("API health")
		withInvalidTags["tags"] = "production"

		withForbiddenTagIDs := httpConfigEntry("API health")
		withForbiddenTagIDs["tagIds"] = []any{}

		tests := []struct {
			name        string
			slug        string
			entry       map[string]any
			errContains string
		}{
			{
				name:        "Invalid slug",
				slug:        "API Health",
				entry:       httpConfigEntry("API health"),
				errContains: "invalid monitor slug",
			},
			{name: "Forbidden key", slug: "api-health", entry: withForbiddenKey, errContains: "runState"},
			{name: "Unknown field", slug: "api-health", entry: withUnknownField, errContains: "intervall"},
			{
				name:        "Unknown probe field",
				slug:        "api-health",
				entry:       withUnknownProbeField,
				errContains: "expectedStatusCode",
			},
			{
				name:        "Invalid probe config",
				slug:        "api-health",
				entry:       withInvalidProbeConfig,
				errContains: "URL cannot be empty",
			},
			{name: "Missing type", slug: "api-health", entry: withoutType, errContains: "type cannot be empty"},
			{name: "Unknown tag", slug: "api-health", entry: withUnknownTag, errContains: "unknown tag"},
			{name: "Tags not a list", slug: "api-health", entry: withInvalidTags, errContains: "tags must be a list"},
			{name: "Forbidden tagIds", slug: "api-health", entry: withForbiddenTagIDs, errContains: "tagIds"},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx, _, _, owner := setupMonitorIntegrationTest(t)

				err := syncConfigMonitors(ctx, owner, map[string]map[string]any{tt.slug: tt.entry})
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Empty(t, getAllMonitors(ctx, t))
			})
		}
	})
}

func TestIntegration_MonitorService_ConfigBasedMonitorIsReadOnly(t *testing.T) {
	setup := func(t *testing.T) (context.Context, *monitors.MonitorService, *monitors.Monitor) {
		ctx, monitorService, _, owner := setupMonitorIntegrationTest(t)
		require.NoError(t, syncConfigMonitors(ctx, owner, map[string]map[string]any{
			"api-health": httpConfigEntry("API health"),
		}))
		return ctx, monitorService, getMonitorBySlug(ctx, t, "api-health")
	}

	t.Run("Update fails with 409", func(t *testing.T) {
		ctx, monitorService, monitor := setup(t)

		updated := *monitor
		updated.Name = "Changed"
		svcErr := monitorService.UpdateMonitor(ctx, updated)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
		assert.Equal(t, "API health", getMonitorBySlug(ctx, t, "api-health").Name)
	})

	t.Run("Delete fails with 409", func(t *testing.T) {
		ctx, monitorService, monitor := setup(t)

		svcErr := monitorService.DeleteMonitor(ctx, monitor.ID.String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
		getMonitorBySlug(ctx, t, "api-health")
	})

	t.Run("State update fails with 409", func(t *testing.T) {
		ctx, monitorService, monitor := setup(t)

		svcErr := monitorService.UpdateMonitorStateByID(ctx, monitor.ID, monitors.MonitorRunStateStopped)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
		assert.Equal(t, monitors.MonitorRunStateActive, getMonitorBySlug(ctx, t, "api-health").RunState)
	})
}
