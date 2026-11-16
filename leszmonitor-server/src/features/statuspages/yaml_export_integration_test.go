package statuspages_test

import (
	"bytes"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestIntegration_StatusPageService_ExportToYAML(t *testing.T) {
	t.Run("Exports config that syncs back to the same status page", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		page := newStatusPage("public", tag.ID)
		page.GroupMonitorsByTag = true
		created := createStatusPage(ctx, t, service, page)

		var buf bytes.Buffer
		require.Nil(t, service.ExportToYAML(ctx, created.ID, &buf))

		var entries map[string]map[string]any
		require.NoError(t, yaml.Unmarshal(buf.Bytes(), &entries))
		require.Contains(t, entries, "public")

		entry := entries["public"]
		assert.Equal(t, []any{"prod"}, entry["tags"])
		for _, key := range []string{"id", "slug", "tagIds", "source", "monitorCount", "createdAt", "updatedAt"} {
			assert.NotContains(t, entry, key)
		}

		require.Nil(t, service.DeleteStatusPage(ctx, created.ID))
		require.NoError(t, syncConfigStatusPages(ctx, entries))

		synced, svcErr := service.GetStatusPageBySlug(ctx, "public")
		require.Nil(t, svcErr)
		assert.Equal(t, created.Title, synced.Title)
		assert.Equal(t, created.Description, synced.Description)
		assert.Equal(t, created.GroupMonitorsByTag, synced.GroupMonitorsByTag)
		assert.Equal(t, created.TimeWindowSeconds, synced.TimeWindowSeconds)
		assert.Equal(t, created.TagIDs, synced.TagIDs)
	})

	t.Run("Returns not found for unknown ID", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)

		var buf bytes.Buffer
		svcErr := service.ExportToYAML(ctx, uuid.New(), &buf)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})
}
