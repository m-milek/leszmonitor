package monitors

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHttpMonitorFromReader(t *testing.T) {
	jsonInput := `{
		"name": "Test Monitor",
		"type": "http",
		"config": {
			"url": "http://example.com",
			"method": "GET",
			"expectedStatusCodes": [200]
		}
	}`

	reader := strings.NewReader(jsonInput)
	monitor, err := FromReader(reader)

	assert.NoError(t, err)
	assert.NotNil(t, monitor)
	assert.Equal(t, "Test Monitor", monitor.GetName())

	httpMonitor, ok := monitor.(*HttpMonitor)
	assert.True(t, ok)
	assert.Equal(t, "http://example.com", httpMonitor.Config.URL)
}

func TestHttpMonitorFromReaderInvalidJSON(t *testing.T) {
	jsonInput := `invalid json`

	reader := strings.NewReader(jsonInput)
	monitor, err := FromReader(reader)

	assert.Error(t, err)
	assert.Nil(t, monitor)
}

func TestHttpMonitorFromReaderMissingType(t *testing.T) {
	jsonInput := `{
		"name": "Test Monitor",
		"config": {
			"url": "http://example.com"
		}
	}`

	reader := strings.NewReader(jsonInput)
	monitor, err := FromReader(reader)

	assert.Error(t, err)
	assert.Nil(t, monitor)
}
