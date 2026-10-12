package monitors_test

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestIntegration_ExportToYAML(t *testing.T) {
	t.Run("Exports a monitor", func(t *testing.T) {
		ctx, client := setupMonitorTagsTest(t)
		tagID := insertTestTag(ctx, t, client, "prod")

		monitor := testMonitor([]uuid.UUID{tagID})
		probeConfig, err := json.Marshal(probe.HTTPProbe{
			Method:              "GET",
			URL:                 "https://example.com",
			ExpectedStatusCodes: []int{200},
		})
		require.NoError(t, err)
		monitor.ProbeConfig = string(probeConfig)

		created, err := monitors.NewMonitorDAO(client.Querier()).InsertMonitor(ctx, monitor)
		require.NoError(t, err)

		var buffer bytes.Buffer
		require.NoError(t, monitors.ExportToYAML(ctx, *created, tags.NewTagDAO(client.Querier()), &buffer))

		var exportedMonitors map[string]map[string]any
		require.NoError(t, yaml.Unmarshal(buffer.Bytes(), &exportedMonitors))
		exported := exportedMonitors[created.Slug]
		for _, key := range config.ForbiddenMonitorConfigKeys {
			assert.NotContains(t, exported, key)
		}
		assert.Equal(t, created.Name, exported["name"])
		assert.Equal(t, []any{"prod"}, exported["tags"])

		output := buffer.String()
		_ = output
	})
}
