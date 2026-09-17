package tags

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_TagService_CreateTag(t *testing.T) {
	t.Run("Successfully creates a tag", func(t *testing.T) {
		ctx, tagService, database, owner := setupTagIntegrationTest(t)

		created, svcErr := tagService.CreateTag(ctx, Tag{
			Name:        "  Production ",
			Description: "Production environment",
			ColorHex:    "#ABC",
		})
		require.Nil(t, svcErr)
		require.NotNil(t, created)

		assert.NotEqual(t, uuid.Nil, created.ID)
		assert.Equal(t, "Production", created.Name)
		assert.Equal(t, "#aabbcc", created.ColorHex, "shorthand hex should be normalized")
		assert.False(t, created.CreatedAt.IsZero())

		// Verify audit log was created
		filter := audit.AuditLogFilter{ResourceID: &created.ID}
		entries, dbErr := audit.NewAuditLogDAO(database.Querier()).
			GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionCreateTag {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				assert.Equal(t, created.ID.String(), entry.ResourceID.String())
				break
			}
		}
		assert.True(t, found, "expected a tag.create audit log entry")
	})

	t.Run("Assigns the ID server-side", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		clientID := uuid.New()
		created, svcErr := tagService.CreateTag(ctx, Tag{
			ID:       clientID,
			Name:     "Staging",
			ColorHex: "#001122",
		})
		require.Nil(t, svcErr)
		assert.NotEqual(t, clientID, created.ID)
	})

	t.Run("Rejects an invalid tag", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		_, svcErr := tagService.CreateTag(ctx, Tag{Name: "No color"})
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)

		_, svcErr = tagService.CreateTag(ctx, Tag{Name: "", ColorHex: "#aabbcc"})
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

	t.Run("Allows duplicate names", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		payload := Tag{Name: "Production", ColorHex: "#aabbcc"}
		first, svcErr := tagService.CreateTag(ctx, payload)
		require.Nil(t, svcErr)

		second, svcErr := tagService.CreateTag(ctx, payload)
		require.Nil(t, svcErr)
		assert.NotEqual(t, first.ID, second.ID)
	})
}

func TestIntegration_TagService_GetTagByID(t *testing.T) {
	t.Run("Returns the tag", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		created, svcErr := tagService.CreateTag(ctx, Tag{
			Name:        "Production",
			Description: "Production environment",
			ColorHex:    "#aabbcc",
		})
		require.Nil(t, svcErr)

		found, svcErr := tagService.GetTagByID(ctx, created.ID.String())
		require.Nil(t, svcErr)
		assert.Equal(t, created.ID, found.ID)
		assert.Equal(t, "Production", found.Name)
		assert.Equal(t, "Production environment", found.Description)
		assert.Equal(t, "#aabbcc", found.ColorHex)
	})

	t.Run("Returns not found for an unknown ID", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		_, svcErr := tagService.GetTagByID(ctx, uuid.New().String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Returns bad request for a malformed ID", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		_, svcErr := tagService.GetTagByID(ctx, "not-a-uuid")
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})
}

func TestIntegration_TagService_GetAllTags(t *testing.T) {
	t.Run("Returns every tag in the instance", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		t1, svcErr := tagService.CreateTag(ctx, Tag{Name: "Production", ColorHex: "#aabbcc"})
		require.Nil(t, svcErr)
		t2, svcErr := tagService.CreateTag(ctx, Tag{Name: "Staging", ColorHex: "#001122"})
		require.Nil(t, svcErr)

		all, svcErr := tagService.GetAllTags(ctx)
		require.Nil(t, svcErr)
		require.Len(t, all, 2)

		ids := make(map[uuid.UUID]bool)
		for _, tag := range all {
			ids[tag.ID] = true
		}
		assert.True(t, ids[t1.ID])
		assert.True(t, ids[t2.ID])
	})

	t.Run("Returns an empty list when there are no tags", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		all, svcErr := tagService.GetAllTags(ctx)
		require.Nil(t, svcErr)
		assert.Empty(t, all)
	})
}

func TestIntegration_TagService_UpdateTag(t *testing.T) {
	t.Run("Successfully updates a tag", func(t *testing.T) {
		ctx, tagService, database, owner := setupTagIntegrationTest(t)

		created, svcErr := tagService.CreateTag(ctx, Tag{
			Name:        "Production",
			Description: "Production environment",
			ColorHex:    "#aabbcc",
		})
		require.Nil(t, svcErr)

		updated, svcErr := tagService.UpdateTag(ctx, Tag{
			ID:          created.ID,
			Name:        "Prod",
			Description: "Renamed",
			ColorHex:    "#123456",
		})
		require.Nil(t, svcErr)
		assert.Equal(t, "Prod", updated.Name)
		assert.Equal(t, "Renamed", updated.Description)
		assert.Equal(t, "#123456", updated.ColorHex)
		require.NotNil(t, updated.UpdatedAt)

		all, svcErr := tagService.GetAllTags(ctx)
		require.Nil(t, svcErr)
		require.Len(t, all, 1)
		assert.Equal(t, "Prod", all[0].Name)

		// Verify audit log records both the before and after state
		filter := audit.AuditLogFilter{ResourceID: &created.ID}
		entries, dbErr := audit.NewAuditLogDAO(database.Querier()).
			GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionUpdateTag {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				require.NotNil(t, entry.Before)
				require.NotNil(t, entry.After)
				assert.Contains(t, *entry.Before, "Production")
				assert.Contains(t, *entry.After, "Prod")
				break
			}
		}
		assert.True(t, found, "expected a tag.update audit log entry")
	})

	t.Run("Returns not found for an unknown ID", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		_, svcErr := tagService.UpdateTag(ctx, Tag{
			ID:       uuid.New(),
			Name:     "Ghost",
			ColorHex: "#aabbcc",
		})
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})

	t.Run("Rejects an invalid tag", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		created, svcErr := tagService.CreateTag(ctx, Tag{Name: "Production", ColorHex: "#aabbcc"})
		require.Nil(t, svcErr)

		_, svcErr = tagService.UpdateTag(ctx, Tag{ID: created.ID, Name: "Prod", ColorHex: "not-a-color"})
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

	t.Run("Allows renaming onto an existing name", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		_, svcErr := tagService.CreateTag(ctx, Tag{Name: "Production", ColorHex: "#aabbcc"})
		require.Nil(t, svcErr)
		staging, svcErr := tagService.CreateTag(ctx, Tag{Name: "Staging", ColorHex: "#001122"})
		require.Nil(t, svcErr)

		updated, svcErr := tagService.UpdateTag(ctx, Tag{
			ID:       staging.ID,
			Name:     "Production",
			ColorHex: "#001122",
		})
		require.Nil(t, svcErr)
		assert.Equal(t, "Production", updated.Name)
	})
}

func TestIntegration_TagService_DeleteTag(t *testing.T) {
	t.Run("Successfully deletes a tag", func(t *testing.T) {
		ctx, tagService, database, owner := setupTagIntegrationTest(t)

		created, svcErr := tagService.CreateTag(ctx, Tag{Name: "Production", ColorHex: "#aabbcc"})
		require.Nil(t, svcErr)

		svcErr = tagService.DeleteTag(ctx, created.ID.String())
		require.Nil(t, svcErr)

		all, svcErr := tagService.GetAllTags(ctx)
		require.Nil(t, svcErr)
		assert.Empty(t, all)

		// Verify audit log was created with the pre-delete state
		filter := audit.AuditLogFilter{ResourceID: &created.ID}
		entries, dbErr := audit.NewAuditLogDAO(database.Querier()).
			GetAuditLogEntries(ctx, filter, util.Pagination{Page: 1, PerPage: 10})
		require.NoError(t, dbErr)

		found := false
		for _, entry := range entries {
			if entry.Action == audit.ActionDeleteTag {
				found = true
				assert.Equal(t, owner.Username, *entry.Username)
				require.NotNil(t, entry.Before)
				assert.Contains(t, *entry.Before, "Production")
				break
			}
		}
		assert.True(t, found, "expected a tag.delete audit log entry")
	})

	t.Run("Returns bad request for a malformed ID", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		svcErr := tagService.DeleteTag(ctx, "not-a-uuid")
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusBadRequest, svcErr.Code)
	})

	t.Run("Returns not found for an unknown ID", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)

		svcErr := tagService.DeleteTag(ctx, uuid.New().String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusNotFound, svcErr.Code)
	})
}
