package monitors

import (
	"testing"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeMonitorJSON(t *testing.T) {
	t.Run("Decodes fields and keeps probeConfig as raw JSON", func(t *testing.T) {
		monitor, err := decodeMonitorJSON([]byte(`{
			"name": "API health",
			"type": "http",
			"interval": 60,
			"probeConfig": {"method": "GET", "url": "https://example.com"}
		}`))
		require.NoError(t, err)

		assert.Equal(t, "API health", monitor.Name)
		assert.Equal(t, kind.HTTPConfigType, monitor.Type)
		assert.Equal(t, 60, monitor.Interval)
		assert.JSONEq(t, `{"method": "GET", "url": "https://example.com"}`, monitor.ProbeConfig)
	})

	t.Run("Leaves probeConfig empty when missing", func(t *testing.T) {
		monitor, err := decodeMonitorJSON([]byte(`{"name": "API health"}`))
		require.NoError(t, err)
		assert.Empty(t, monitor.ProbeConfig)
	})

	t.Run("Rejects unknown fields", func(t *testing.T) {
		_, err := decodeMonitorJSON([]byte(`{"name": "API health", "intervall": 60}`))
		require.Error(t, err)
		assert.Contains(t, err.Error(), `unknown field "intervall"`)
	})

	t.Run("Rejects invalid JSON", func(t *testing.T) {
		_, err := decodeMonitorJSON([]byte(`{"name": `))
		require.Error(t, err)
	})
}
