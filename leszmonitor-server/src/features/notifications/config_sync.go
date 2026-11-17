package notifications

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/pkg/errors"
)

var forbiddenKeys = []string{
	"id", "name", "source", "createdAt", "updatedAt",
}

var uuidNamespace = uuid.MustParse("c0ffee00-beef-face-cafe-feedbadcab1e")

type channelConfigEntry struct {
	Kind    NotificationChannelKind `json:"kind"`
	Enabled *bool                   `json:"enabled"`
	Config  json.RawMessage         `json:"config"`
}

func SynchronizeConfigBasedChannels(ctx context.Context, database db.DB, rawEntries map[string]map[string]any) error {
	configBasedChannels := make([]NotificationChannel, 0)

	for name, entry := range rawEntries {
		for k := range entry {
			if slices.Contains(forbiddenKeys, k) {
				return fmt.Errorf("forbidden key in notification channel config (%s): %s", name, k)
			}
		}

		expandedConfig, err := expandEnv(entry["config"])
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to expand notification channel config (%s)", name))
		}
		expandedEntry := map[string]any{}
		for k, v := range entry {
			expandedEntry[k] = v
		}
		expandedEntry["config"] = expandedConfig

		jsonEntry, err := json.Marshal(expandedEntry)
		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to parse notification channel config (%s)", name))
		}

		var parsed channelConfigEntry
		decoder := json.NewDecoder(bytes.NewReader(jsonEntry))
		decoder.DisallowUnknownFields()
		if err = decoder.Decode(&parsed); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to decode notification channel config (%s)", name))
		}

		channel := NotificationChannel{
			ID:      uuid.NewSHA1(uuidNamespace, []byte(name)),
			Name:    name,
			Kind:    parsed.Kind,
			Enabled: parsed.Enabled == nil || *parsed.Enabled,
			Config:  ChannelConfig(parsed.Config),
			Source:  util.ResourceSourceConfig,
		}
		if err = channel.prepare(nil); err != nil {
			return errors.Wrap(err, fmt.Sprintf("failed to validate notification channel config (%s)", name))
		}
		if channel.Name != name {
			return fmt.Errorf("invalid notification channel name: %q", name)
		}

		configBasedChannels = append(configBasedChannels, channel)
	}

	deletedCount := 0
	err := database.WithTx(ctx, func(q db.Querier) error {
		dao := NewNotificationChannelDAO(q)
		allChannels, err := dao.GetAllChannels(ctx)
		if err != nil {
			return err
		}

		for _, channel := range allChannels {
			isNotConfigBased := channel.Source != util.ResourceSourceConfig
			isInCurrentFile := slices.ContainsFunc(configBasedChannels, func(c NotificationChannel) bool {
				return c.ID == channel.ID
			})

			if isNotConfigBased || isInCurrentFile {
				continue
			}

			if _, err = dao.DeleteChannelByID(ctx, channel.ID); err != nil {
				return err
			}
			deletedCount++
		}

		for i := range configBasedChannels {
			if err = dao.UpsertConfigBasedChannel(ctx, &configBasedChannels[i]); err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to upsert notification channel (%s)", configBasedChannels[i].Name))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "failed to synchronize notification channel config in DB")
	}

	log.FromContext(ctx).Info().
		Int("synced", len(configBasedChannels)).
		Int("deleted", deletedCount).
		Msg("Config based notification channels synchronized")

	return nil
}

func expandEnv(value any) (any, error) {
	switch v := value.(type) {
	case string:
		var missing []string
		expanded := os.Expand(v, func(name string) string {
			envValue, ok := os.LookupEnv(name)
			if !ok {
				missing = append(missing, name)
			}
			return envValue
		})
		if len(missing) > 0 {
			return nil, fmt.Errorf("environment variables not set: %v", missing)
		}
		return expanded, nil
	case map[string]any:
		out := make(map[string]any, len(v))
		for key, item := range v {
			expanded, err := expandEnv(item)
			if err != nil {
				return nil, err
			}
			out[key] = expanded
		}
		return out, nil
	case []any:
		out := make([]any, len(v))
		for i, item := range v {
			expanded, err := expandEnv(item)
			if err != nil {
				return nil, err
			}
			out[i] = expanded
		}
		return out, nil
	default:
		return value, nil
	}
}
