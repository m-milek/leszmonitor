package statuspages_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/statuspages"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_StatusPageService_CreateStatusPage(t *testing.T) {
	t.Run("Creates a status page with tags", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		insertMonitorWithTags(ctx, t, tag.ID)

		created := createStatusPage(ctx, t, service, newStatusPage("public", tag.ID))

		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, "public", created.Slug)
		assert.Equal(t, util.ResourceSourceUI, created.Source)
		assert.Equal(t, []uuid.UUID{tag.ID}, created.TagIDs)
		assert.Equal(t, 1, created.MonitorCount)
	})

	t.Run("Rejects an invalid status page", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)

		_, svcErr := service.CreateStatusPage(ctx, newStatusPage("public"))
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

	t.Run("Rejects a duplicate slug", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		createStatusPage(ctx, t, service, newStatusPage("public", tag.ID))

		_, svcErr := service.CreateStatusPage(ctx, newStatusPage("public", tag.ID))
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
	})
}

func TestIntegration_StatusPageService_GetStatusPages(t *testing.T) {
	t.Run("Gets status pages by ID, slug and all", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		created := createStatusPage(ctx, t, service, newStatusPage("public", tag.ID))

		byID, svcErr := service.GetStatusPageByID(ctx, created.ID)
		require.Nil(t, svcErr)
		assert.Equal(t, created.Slug, byID.Slug)

		bySlug, svcErr := service.GetStatusPageBySlug(ctx, "public")
		require.Nil(t, svcErr)
		assert.Equal(t, created.ID, bySlug.ID)

		all, svcErr := service.GetAllStatusPages(ctx)
		require.Nil(t, svcErr)
		require.Len(t, all, 1)
		assert.Equal(t, []uuid.UUID{tag.ID}, all[0].TagIDs)
	})

	t.Run("Returns not found for unknown ID and slug", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)

		_, svcErr := service.GetStatusPageByID(ctx, uuid.New())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)

		_, svcErr = service.GetStatusPageBySlug(ctx, "missing")
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})
}

func TestIntegration_StatusPageService_UpdateStatusPage(t *testing.T) {
	t.Run("Updates fields and replaces tags", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		prod := insertTag(ctx, t, "prod")
		staging := insertTag(ctx, t, "staging")
		created := createStatusPage(ctx, t, service, newStatusPage("public", prod.ID))

		update := newStatusPage("renamed", staging.ID)
		update.ID = created.ID
		update.RequireAuth = true
		update.TimeWindowSeconds = 7200

		updated, svcErr := service.UpdateStatusPage(ctx, update)
		require.Nil(t, svcErr)
		assert.Equal(t, "renamed", updated.Slug)
		assert.True(t, updated.RequireAuth)
		assert.Equal(t, 7200, updated.TimeWindowSeconds)
		assert.Equal(t, []uuid.UUID{staging.ID}, updated.TagIDs)
	})

	t.Run("Returns not found for unknown ID", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")

		update := newStatusPage("public", tag.ID)
		update.ID = uuid.New()

		_, svcErr := service.UpdateStatusPage(ctx, update)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Rejects updating a config-based status page", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		page := insertConfigStatusPage(ctx, t, tag.ID)

		update := newStatusPage("public", tag.ID)
		update.ID = page.ID

		_, svcErr := service.UpdateStatusPage(ctx, update)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
	})
}

func TestIntegration_StatusPageService_DeleteStatusPage(t *testing.T) {
	t.Run("Deletes a status page", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		created := createStatusPage(ctx, t, service, newStatusPage("public", tag.ID))

		require.Nil(t, service.DeleteStatusPage(ctx, created.ID))

		_, svcErr := service.GetStatusPageByID(ctx, created.ID)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Returns not found for unknown ID", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)

		svcErr := service.DeleteStatusPage(ctx, uuid.New())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Rejects deleting a config-based status page", func(t *testing.T) {
		ctx, service := setupStatusPageIntegrationTest(t)
		tag := insertTag(ctx, t, "prod")
		page := insertConfigStatusPage(ctx, t, tag.ID)

		svcErr := service.DeleteStatusPage(ctx, page.ID)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
	})
}

func insertConfigStatusPage(ctx context.Context, t *testing.T, tagIDs ...uuid.UUID) statuspages.StatusPage {
	t.Helper()
	page := newStatusPage("config-page", tagIDs...)
	page.ID = uuid.New()
	page.Source = util.ResourceSourceConfig
	require.NoError(t, statuspages.NewStatusPageDAO(db.Get().Querier()).UpsertConfigBasedStatusPage(ctx, &page))
	return page
}
