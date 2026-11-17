package notifications

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/notifications/providers"
	"github.com/m-milek/leszmonitor/platform/util"
)

type NotificationChannelKind string

const (
	ChannelKindWebhook NotificationChannelKind = "webhook"
	ChannelKindNtfy    NotificationChannelKind = "ntfy"
	ChannelKindGotify  NotificationChannelKind = "gotify"
)

type NotificationChannel struct {
	ID              uuid.UUID               `json:"id"      yaml:"id"      db:"id"`
	Name            string                  `json:"name"    yaml:"name"    db:"name"`
	Kind            NotificationChannelKind `json:"kind"    yaml:"kind"    db:"kind"`
	Enabled         bool                    `json:"enabled" yaml:"enabled" db:"enabled"`
	Config          ChannelConfig           `json:"config"  yaml:"config"  db:"config"  swaggertype:"object"`
	Source          util.ResourceSource     `json:"source"  yaml:"source"  db:"source"`
	util.Timestamps `yaml:",inline"`
}

type ChannelConfig string

func (c ChannelConfig) MarshalJSON() ([]byte, error) {
	if c == "" {
		return []byte("null"), nil
	}
	return []byte(c), nil
}

func (c *ChannelConfig) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*c = ""
		return nil
	}
	*c = ChannelConfig(data)
	return nil
}

func parseKindConfig(kind NotificationChannelKind, raw ChannelConfig) (providers.Config, error) {
	var cfg providers.Config
	switch kind {
	case ChannelKindWebhook:
		cfg = &providers.WebhookConfig{}
	case ChannelKindNtfy:
		cfg = &providers.NtfyConfig{}
	case ChannelKindGotify:
		cfg = &providers.GotifyConfig{}
	default:
		return nil, fmt.Errorf("unknown channel kind %q", kind)
	}

	decoder := json.NewDecoder(strings.NewReader(string(raw)))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(cfg); err != nil {
		return nil, fmt.Errorf("invalid %s config: %w", kind, err)
	}
	return cfg, nil
}

// prepare normalizes and validates the channel. Secrets left empty are taken from previous, if it has the same kind.
func (c *NotificationChannel) prepare(previous *NotificationChannel) error {
	c.Name = strings.TrimSpace(c.Name)
	if c.Name == "" {
		return fmt.Errorf("channel name cannot be empty")
	}

	cfg, err := parseKindConfig(c.Kind, c.Config)
	if err != nil {
		return err
	}

	if previous != nil && previous.Kind == c.Kind {
		prevCfg, prevErr := parseKindConfig(previous.Kind, previous.Config)
		if prevErr != nil {
			return fmt.Errorf("stored config is invalid: %w", prevErr)
		}
		cfg.KeepSecretsFrom(prevCfg)
	}

	if err = cfg.Validate(); err != nil {
		return err
	}

	encoded, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	c.Config = ChannelConfig(encoded)
	return nil
}

func (c NotificationChannel) redacted() (NotificationChannel, error) {
	cfg, err := parseKindConfig(c.Kind, c.Config)
	if err != nil {
		return NotificationChannel{}, err
	}
	cfg.Redact()

	encoded, err := json.Marshal(cfg)
	if err != nil {
		return NotificationChannel{}, err
	}
	c.Config = ChannelConfig(encoded)
	return c, nil
}
