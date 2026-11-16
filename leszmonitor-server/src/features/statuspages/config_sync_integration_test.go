package statuspages_test

import (
	"context"
	"testing"

	"github.com/m-milek/leszmonitor/features/statuspages"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func statusPageConfigEntry(title string, tagNames ...any) map[string]any {
	return map[string]any{
		"title":              title,
		"description":        "Defined in config",
		"groupMonitorsByTag": true,
		"requireAuth":        false,
		"enabled":            true,
		"timeWindowSeconds":  3600,
		"tags":               tagNames,
	}
}

func syncConfigStatusPages(ctx context.Context, entries map[string]map[string]any) error {
	return statuspages.SynchronizeConfigBasedStatusPages(ctx, db.Get(), entries)
}

func TestIntegration_SynchronizeConfigBasedStatusPages(t *testing.T) {
	t.Run("Inserts status pages from config", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")

		require.NoError(t, syncConfigStatusPages(ctx, map[string]map[string]any{
			"public": statusPageConfigEntry("Public", "prod"),
		}))

		page, svcErr := service.GetStatusPageBySlug(ctx, "public")
		require.Nil(t, svcErr)
		assert.Equal(t, "Public", page.Title)
		assert.True(t, page.GroupMonitorsByTag)
		assert.Equal(t, 3600, page.TimeWindowSeconds)
		assert.Equal(t, util.ResourceSourceConfig, page.Source)
		assert.Len(t, page.TagIDs, 1)
		assert.Equal(t, tag.ID, page.TagIDs[0])
	})

	t.Run("Keeps the same ID and updates fields on re-sync", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		insertTag(ctx, t, "prod")

		require.NoError(t, syncConfigStatusPages(ctx, map[string]map[string]any{
			"public": statusPageConfigEntry("Public", "prod"),
		}))
		before, svcErr := service.GetStatusPageBySlug(ctx, "public")
		require.Nil(t, svcErr)

		require.NoError(t, syncConfigStatusPages(ctx, map[string]map[string]any{
			"public": statusPageConfigEntry("Renamed", "prod"),
		}))
		after, svcErr := service.GetStatusPageBySlug(ctx, "public")
		require.Nil(t, svcErr)

		assert.Equal(t, before.ID, after.ID)
		assert.Equal(t, "Renamed", after.Title)
	})

	t.Run("Deletes status pages missing from config and keeps UI pages", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		uiPage := createStatusPage(ctx, t, service, newStatusPage("ui-page", tag.ID))

		require.NoError(t, syncConfigStatusPages(ctx, map[string]map[string]any{
			"public": statusPageConfigEntry("Public", "prod"),
		}))
		require.NoError(t, syncConfigStatusPages(ctx, map[string]map[string]any{}))

		all, svcErr := service.GetAllStatusPages(ctx)
		require.Nil(t, svcErr)
		require.Len(t, all, 1)
		assert.Equal(t, uiPage.ID, all[0].ID)
	})

	invalidEntry := func(modify func(entry map[string]any)) map[string]any {
		entry := statusPageConfigEntry("Public", "prod")
		modify(entry)
		return entry
	}

	tests := []struct {
		name        string
		slug        string
		entry       map[string]any
		errContains string
	}{
		{"Invalid slug", "Bad Slug", statusPageConfigEntry("Public", "prod"), "invalid status page slug"},
		{"Forbidden key", "public", invalidEntry(func(e map[string]any) { e["source"] = "ui" }), "forbidden key"},
		{"Missing enabled", "public", invalidEntry(func(e map[string]any) { delete(e, "enabled") }), "missing key"},
		{"Missing time window", "public", invalidEntry(func(e map[string]any) { delete(e, "timeWindowSeconds") }), "missing key"},
		{"No tags", "public", invalidEntry(func(e map[string]any) { e["tags"] = []any{} }), "at least one tag"},
		{"Unknown tag", "public", statusPageConfigEntry("Public", "missing"), "unknown tag"},
		{"Unknown field", "public", invalidEntry(func(e map[string]any) { e["color"] = "red" }), "failed to decode"},
		{"Empty title", "public", invalidEntry(func(e map[string]any) { e["title"] = "" }), "title cannot be empty"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, _ := setupStatusPageIntegrationTest(t)
			insertTag(ctx, t, "prod")

			err := syncConfigStatusPages(ctx, map[string]map[string]any{tt.slug: tt.entry})
			assert.ErrorContains(t, err, tt.errContains)
		})
	}
}
