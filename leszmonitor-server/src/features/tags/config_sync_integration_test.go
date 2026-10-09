package tags_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tagConfigEntry(colorHex string) map[string]any {
	return map[string]any{
		"description": "Defined in config",
		"colorHex":    colorHex,
	}
}

func syncConfigTags(ctx context.Context, entries map[string]map[string]any) error {
	return tags.SynchronizeConfigBasedTags(ctx, db.Get(), entries)
}

func getAllTags(ctx context.Context, t *testing.T) []tags.Tag {
	t.Helper()
	all, err := tags.NewTagDAO(db.Get().Querier()).GetAllTags(ctx)
	require.NoError(t, err)
	return all
}

func findTagByName(t *testing.T, all []tags.Tag, name string) tags.Tag {
	t.Helper()
	for _, tag := range all {
		if tag.Name == name {
			return tag
		}
	}
	require.Failf(t, "tag not found", "tag %q not found", name)
	return tags.Tag{}
}

func TestIntegration_SynchronizeConfigBasedTags(t *testing.T) {
	t.Run("Inserts tags from config", func(t *testing.T) {
		ctx, _, _, _ := setupTagIntegrationTest(t)

		require.NoError(t, syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#ABC"),
		}))

		tag := findTagByName(t, getAllTags(ctx, t), "production")
		assert.Equal(t, "Defined in config", tag.Description)
		assert.Equal(t, "#aabbcc", tag.ColorHex)
		assert.Equal(t, tags.TagSourceConfig, tag.Source)
	})

	t.Run("Derives the same ID on every run and updates fields", func(t *testing.T) {
		ctx, _, _, _ := setupTagIntegrationTest(t)

		require.NoError(t, syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#aabbcc"),
		}))
		before := findTagByName(t, getAllTags(ctx, t), "production")

		require.NoError(t, syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#001122"),
		}))
		after := findTagByName(t, getAllTags(ctx, t), "production")

		assert.Equal(t, before.ID, after.ID)
		assert.Equal(t, "#001122", after.ColorHex)
	})

	t.Run("Deletes tags missing from config and keeps UI tags", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)
		uiTag, svcErr := tagService.CreateTag(ctx, tags.Tag{Name: "staging", ColorHex: "#aabbcc"})
		require.Nil(t, svcErr)

		require.NoError(t, syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#aabbcc"),
			"critical":   tagConfigEntry("#aabbcc"),
		}))
		require.NoError(t, syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#aabbcc"),
		}))

		all := getAllTags(ctx, t)
		require.Len(t, all, 2)
		findTagByName(t, all, "production")
		assert.Equal(t, uiTag.ID, findTagByName(t, all, "staging").ID)

		require.NoError(t, syncConfigTags(ctx, nil))
		all = getAllTags(ctx, t)
		require.Len(t, all, 1)
		assert.Equal(t, uiTag.ID, all[0].ID)
	})

	t.Run("Fails and writes nothing when name is used by a UI tag", func(t *testing.T) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)
		_, svcErr := tagService.CreateTag(ctx, tags.Tag{Name: "staging", ColorHex: "#aabbcc"})
		require.Nil(t, svcErr)

		err := syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#aabbcc"),
			"staging":    tagConfigEntry("#aabbcc"),
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "staging")
		assert.Len(t, getAllTags(ctx, t), 1)
	})

	t.Run("Rejects invalid entries without writing anything", func(t *testing.T) {
		withForbiddenKey := tagConfigEntry("#aabbcc")
		withForbiddenKey["name"] = "other"

		withUnknownField := tagConfigEntry("#aabbcc")
		withUnknownField["color"] = "#aabbcc"

		tests := []struct {
			name        string
			tagName     string
			entry       map[string]any
			errContains string
		}{
			{name: "Forbidden key", tagName: "production", entry: withForbiddenKey, errContains: "name"},
			{
				name:        "Unknown field",
				tagName:     "production",
				entry:       withUnknownField,
				errContains: `unknown field "color"`,
			},
			{name: "Invalid color", tagName: "production", entry: tagConfigEntry("red"), errContains: "colorHex"},
			{name: "Missing color", tagName: "production", entry: nil, errContains: "colorHex cannot be empty"},
			{
				name:        "Name with surrounding spaces",
				tagName:     " production ",
				entry:       tagConfigEntry("#aabbcc"),
				errContains: "invalid tag name",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ctx, _, _, _ := setupTagIntegrationTest(t)

				err := syncConfigTags(ctx, map[string]map[string]any{tt.tagName: tt.entry})
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errContains)
				assert.Empty(t, getAllTags(ctx, t))
			})
		}
	})
}

func TestIntegration_TagService_ConfigBasedTagIsReadOnly(t *testing.T) {
	setup := func(t *testing.T) (context.Context, *tags.TagService, tags.Tag) {
		ctx, tagService, _, _ := setupTagIntegrationTest(t)
		require.NoError(t, syncConfigTags(ctx, map[string]map[string]any{
			"production": tagConfigEntry("#aabbcc"),
		}))
		return ctx, tagService, findTagByName(t, getAllTags(ctx, t), "production")
	}

	t.Run("Update fails with 409", func(t *testing.T) {
		ctx, tagService, tag := setup(t)

		updated := tag
		updated.ColorHex = "#001122"
		_, svcErr := tagService.UpdateTag(ctx, updated)
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
		assert.Equal(t, "#aabbcc", findTagByName(t, getAllTags(ctx, t), "production").ColorHex)
	})

	t.Run("Delete fails with 409", func(t *testing.T) {
		ctx, tagService, tag := setup(t)

		svcErr := tagService.DeleteTag(ctx, tag.ID.String())
		require.NotNil(t, svcErr)
		assert.Equal(t, http.StatusConflict, svcErr.Code)
		assert.Len(t, getAllTags(ctx, t), 1)
	})
}
