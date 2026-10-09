package tags

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/pkg/errors"
)

var forbiddenKeys = []string{
	"id", "name", "source", "createdAt", "updatedAt",
}

var uuidNamespace = uuid.MustParse("abadcafe-deaf-dead-beef-cafebabefeed")

func SynchronizeConfigBasedTags(ctx context.Context, database db.DB, rawEntries map[string]map[string]any) error {
	configBasedTags := make([]Tag, 0)

	for name, entry := range rawEntries {
		for k := range entry {
			if slices.Contains(forbiddenKeys, k) {
				return fmt.Errorf("forbidden key in tag config (%s): %s", name, k)
			}
		}

		jsonTagConfig, err := json.Marshal(entry)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse tag config (%s)", name))
		}

		var tag Tag
		decoder := json.NewDecoder(bytes.NewReader(jsonTagConfig))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&tag); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode tag config (%s)", name))
		}

		tag.Name = name
		tag.Normalize()
		if tag.Name != name {
			return fmt.Errorf("invalid tag name: %q", name)
		}

		tag.ID = uuid.NewSHA1(uuidNamespace, []byte(name))
		tag.Source = TagSourceConfig

		if err = tag.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to validate tag config (%s)", name))
		}

		configBasedTags = append(configBasedTags, tag)
	}

	deletedCount := 0
	err := database.WithTx(ctx, func(q db.Querier) error {
		dao := NewTagDAO(q)
		allTags, err := dao.GetAllTags(ctx)
		if err != nil {
			return err
		}

		for _, tag := range allTags {
			isNotConfigBased := tag.Source != TagSourceConfig
			isInCurrentFile := slices.ContainsFunc(configBasedTags, func(t Tag) bool {
				return t.ID == tag.ID
			})

			if isNotConfigBased || isInCurrentFile {
				continue
			}

			if _, err = dao.DeleteTagByID(ctx, tag.ID); err != nil {
				return err
			}
			deletedCount++
		}

		for i := range configBasedTags {
			if err = dao.UpsertConfigBasedTag(ctx, &configBasedTags[i]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to upsert tag (%s)", configBasedTags[i].Name))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to synchronize tag config in DB")
	}

	log.FromContext(ctx).Info().
		Int("synced", len(configBasedTags)).
		Int("deleted", deletedCount).
		Msg("Config based tags synchronized")

	return nil
}
