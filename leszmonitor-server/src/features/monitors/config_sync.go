package monitors

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"
	"slices"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/pkg/errors"
)

var forbiddenKeys = []string{
	"id", "slug", "runState", "ownerId", "tagIds", "source", "createdAt", "updatedAt",
}

var uuidNamespace = uuid.MustParse("feedface-dead-beef-cafe-baddecafbeef")

func SynchronizeConfigBasedMonitors(
	ctx context.Context,
	database db.DB,
	rawEntries map[string]map[string]any,
	ownerID uuid.UUID,
) error {
	configBasedMonitors := make([]Monitor, 0)
	tagNamesBySlug := make(map[string][]string)

	for slug, entry := range rawEntries {
		isSlugValid := slug == util.SlugFromString(slug) && len(slug) >= 2 && len(slug) <= 50
		if !isSlugValid {
			return fmt.Errorf("invalid monitor slug: %s", slug)
		}

		for k := range entry {
			if slices.Contains(forbiddenKeys, k) {
				return fmt.Errorf("forbidden key in monitor config (%s): %s", slug, k)
			}
		}

		tagNames, err := parseTagNames(entry["tags"])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("invalid tags in monitor config (%s)", slug))
		}
		tagNamesBySlug[slug] = tagNames

		monitorEntry := maps.Clone(entry)
		delete(monitorEntry, "tags")

		jsonMonitorConfig, err := json.Marshal(monitorEntry)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse monitor config (%s)", slug))
		}

		monitor, err := decodeMonitorJSON(jsonMonitorConfig)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode monitor config (%s)", slug))
		}

		uniqueKey := string(monitor.Type) + "/" + slug
		monitor.ID = uuid.NewSHA1(uuidNamespace, []byte(uniqueKey))

		monitor.Slug = slug
		monitor.Source = MonitorSourceConfig
		monitor.RunState = MonitorRunStateActive
		monitor.OwnerID = ownerID
		if monitor.ResultRetentionSeconds == 0 {
			monitor.ResultRetentionSeconds = defaultResultRetentionSeconds
		}

		if err = monitor.Validate(); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to validate monitor config (%s)", slug))
		}

		if _, err = probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse monitor's probe config (%s)", slug))
		}

		configBasedMonitors = append(configBasedMonitors, monitor)
	}

	deletedCount := 0
	err := database.WithTx(ctx, func(database db.Querier) error {
		allTags, err := tags.NewTagDAO(database).GetAllTags(ctx)
		if err != nil {
			return err
		}

		tagIDsByName := make(map[string]uuid.UUID, len(allTags))
		for _, tag := range allTags {
			tagIDsByName[tag.Name] = tag.ID
		}

		for i := range configBasedMonitors {
			monitor := &configBasedMonitors[i]
			monitor.TagIDs = make([]uuid.UUID, 0)
			for _, tagName := range tagNamesBySlug[monitor.Slug] {
				tagID, ok := tagIDsByName[tagName]
				if !ok {
					return fmt.Errorf("unknown tag in monitor config (%s): %s", monitor.Slug, tagName)
				}
				monitor.TagIDs = append(monitor.TagIDs, tagID)
			}
		}

		dao := NewMonitorDAO(database)
		allMonitors, err := dao.GetAllMonitors(ctx)
		if err != nil {
			return err
		}

		for _, monitor := range allMonitors {
			isNotConfigBased := monitor.Source != MonitorSourceConfig
			isInCurrentFile := slices.ContainsFunc(configBasedMonitors, func(m Monitor) bool {
				return m.ID == monitor.ID
			})

			if isNotConfigBased || isInCurrentFile {
				continue
			}

			if _, err = dao.DeleteMonitorByID(ctx, monitor.ID); err != nil {
				return err
			}
			deletedCount++
		}

		for i := range configBasedMonitors {
			if err := dao.UpsertConfigBasedMonitor(ctx, &configBasedMonitors[i]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to upsert monitor (%s)", configBasedMonitors[i].Slug))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to synchronize monitor config in DB")
	}

	log.FromContext(ctx).Info().
		Int("synced", len(configBasedMonitors)).
		Int("deleted", deletedCount).
		Msg("Config based monitors synchronized")

	return nil
}

func parseTagNames(rawTags any) ([]string, error) {
	if rawTags == nil {
		return nil, nil
	}

	list, ok := rawTags.([]any)
	if !ok {
		return nil, fmt.Errorf("tags must be a list of tag names")
	}

	tagNames := make([]string, 0, len(list))
	for _, rawTag := range list {
		tagName, isString := rawTag.(string)
		if !isString {
			return nil, fmt.Errorf("tags must be a list of tag names")
		}
		tagNames = append(tagNames, tagName)
	}

	return tagNames, nil
}
