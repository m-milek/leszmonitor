package monitors_test

import (
	"bytes"
	"context"
	"encoding/json"
	"testing"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/features/tags"
	"github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func insertExportTestMonitor(ctx context.Context, t *testing.T) *monitors.Monitor {
	t.Helper()
	tag, err := tags.NewTagDAO(db.Get().Querier()).
		InsertTag(ctx, tags.Tag{Name: "prod", ColorHex: "#aabbcc", Source: tags.TagSourceUI})
	require.NoError(t, err)

	probeConfig, err := json.Marshal(probe.HTTPProbe{
		Method:              "GET",
		URL:                 "https://example.com",
		ExpectedStatusCodes: []int{200},
	})
	require.NoError(t, err)

	monitor := testMonitor([]uuid.UUID{tag.ID})
	monitor.ProbeConfig = monitors.ProbeConfig(string(probeConfig))
	created, err := monitors.NewMonitorDAO(db.Get().Querier()).InsertMonitor(ctx, monitor)
	require.NoError(t, err)
	return created
}

func TestIntegration_ExportToYAML(t *testing.T) {
	t.Run("Exports a monitor", func(t *testing.T) {
		ctx, monitorService, _, _ := setupMonitorIntegrationTest(t)
		created := insertExportTestMonitor(ctx, t)

		var buffer bytes.Buffer
		require.Nil(t, monitorService.ExportToYAML(ctx, created.ID, &buffer))

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

	t.Run("Exported monitor round trips through config sync", func(t *testing.T) {
		ctx, monitorService, _, owner := setupMonitorIntegrationTest(t)
		created := insertExportTestMonitor(ctx, t)

		var firstExport bytes.Buffer
		require.Nil(t, monitorService.ExportToYAML(ctx, created.ID, &firstExport))
		_, err := monitors.NewMonitorDAO(db.Get().Querier()).DeleteMonitorByID(ctx, created.ID)
		require.NoError(t, err)

		var entries map[string]map[string]any
		require.NoError(t, yaml.Unmarshal(firstExport.Bytes(), &entries))
		require.NoError(t, syncConfigMonitors(ctx, owner, entries))

		var secondExport bytes.Buffer
		synced := getMonitorBySlug(ctx, t, created.Slug)
		require.Nil(t, monitorService.ExportToYAML(ctx, synced.ID, &secondExport))
		assert.Equal(t, firstExport.String(), secondExport.String())
	})
}
