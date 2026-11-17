package notifications

import (
	"io"

	"github.com/m-milek/leszmonitor/features/notifications/providers"
	"gopkg.in/yaml.v3"
)

type channelYAML struct {
	Kind    NotificationChannelKind `yaml:"kind"`
	Enabled bool                    `yaml:"enabled"`
	Config  providers.Config        `yaml:"config"`
}

func exportToYAML(channel NotificationChannel, writer io.Writer) error {
	cfg, err := parseKindConfig(channel.Kind, channel.Config)
	if err != nil {
		return err
	}
	cfg.PlaceholderSecrets(channel.Name)

	root := map[string]channelYAML{
		channel.Name: {
			Kind:    channel.Kind,
			Enabled: channel.Enabled,
			Config:  cfg,
		},
	}
	return yaml.NewEncoder(writer).Encode(root)
}
