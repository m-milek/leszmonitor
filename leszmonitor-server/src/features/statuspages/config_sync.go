package statuspages

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/pkg/errors"
)

var forbiddenKeys = []string{
	"id", "slug", "tagIds", "source", "monitorCount", "createdAt", "updatedAt",
}

var requiredKeys = []string{
	"enabled", "timeWindowSeconds",
}

var uuidNamespace = uuid.MustParse("decafbad-c0de-face-feed-deadbeefcafe")

func SynchronizeConfigBasedStatusPages(ctx context.Context, database db.DB, rawEntries map[string]map[string]any) error {
	configBasedPages := make([]StatusPage, 0)
	tagNamesBySlug := make(map[string][]string)

	for slug, entry := range rawEntries {
		isSlugValid := slug == util.SlugFromString(slug) && len(slug) >= 2 && len(slug) <= 50
		if !isSlugValid {
			return fmt.Errorf("invalid status page slug: %s", slug)
		}

		for k := range entry {
			if slices.Contains(forbiddenKeys, k) {
				return fmt.Errorf("forbidden key in status page config (%s): %s", slug, k)
			}
		}

		for _, k := range requiredKeys {
			if _, ok := entry[k]; !ok {
				return fmt.Errorf("missing key in status page config (%s): %s", slug, k)
			}
		}

		tagNames, err := monitors.ParseTagNames(entry["tags"])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid tags in status page config (%s)", slug))
		}
		if len(tagNames) == 0 {
			return fmt.Errorf("status page config (%s) must have at least one tag", slug)
		}
		tagNamesBySlug[slug] = tagNames

		pageEntry := maps.Clone(entry)
		delete(pageEntry, "tags")

		jsonPageConfig, err := json.Marshal(pageEntry)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse status page config (%s)", slug))
		}

		var page StatusPage
		decoder := json.NewDecoder(bytes.NewReader(jsonPageConfig))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&page); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode status page config (%s)", slug))
		}

		page.ID = uuid.NewSHA1(uuidNamespace, []byte(slug))
		page.Slug = slug
		page.Source = util.ResourceSourceConfig

		if err = page.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to validate status page config (%s)", slug))
		}

		configBasedPages = append(configBasedPages, page)
	}

	deletedCount := 0
	err := database.WithTx(ctx, func(q db.Querier) error {
		allTags, err := tags.NewTagDAO(q).GetAllTags(ctx)
		if err != nil {
			return err
		}

		tagIDsByName := make(map[string]uuid.UUID, len(allTags))
		for _, tag := range allTags {
			tagIDsByName[tag.Name] = tag.ID
		}

		for i := range configBasedPages {
			page := &configBasedPages[i]
			page.TagIDs = make([]uuid.UUID, 0)
			for _, tagName := range tagNamesBySlug[page.Slug] {
				tagID, ok := tagIDsByName[tagName]
				if !ok {
					return fmt.Errorf("unknown tag in status page config (%s): %s", page.Slug, tagName)
				}
				page.TagIDs = append(page.TagIDs, tagID)
			}
		}

		dao := NewStatusPageDAO(q)
		allPages, err := dao.GetAllStatusPages(ctx)
		if err != nil {
			return err
		}

		for _, page := range allPages {
			isNotConfigBased := page.Source != util.ResourceSourceConfig
			isInCurrentFile := slices.ContainsFunc(configBasedPages, func(p StatusPage) bool {
				return p.ID == page.ID
			})

			if isNotConfigBased || isInCurrentFile {
				continue
			}

			if _, err = dao.DeleteStatusPageByID(ctx, page.ID); err != nil {
				return err
			}
			deletedCount++
		}

		for i := range configBasedPages {
			if err = dao.UpsertConfigBasedStatusPage(ctx, &configBasedPages[i]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to upsert status page (%s)", configBasedPages[i].Slug))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to synchronize status page config in DB")
	}

	log.FromContext(ctx).Info().
		Int("synced", len(configBasedPages)).
		Int("deleted", deletedCount).
		Msg("Config based status pages synchronized")

	return nil
}
