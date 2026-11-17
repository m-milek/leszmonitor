package notifications

import (
	"github.com/google/uuid"
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
	Config          string                  `json:"config"  yaml:"config"  db:"config"`
	Source          util.ResourceSource     `json:"source"  yaml:"source"  db:"source"`
	util.Timestamps `yaml:",inline"`
}
