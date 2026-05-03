package monitors

import (
	"testing"

	shared "github.com/m-milek/leszmonitor/models/consts"
	"github.com/stretchr/testify/assert"
)

func TestHttpMonitorType(t *testing.T) {
	monitor := &HttpMonitor{
		BaseMonitor: BaseMonitor{
			Type: shared.HttpConfigType,
		},
	}
	assert.Equal(t, shared.HttpConfigType, monitor.GetType())
}

func TestHttpMonitorGetConfig(t *testing.T) {
	config := HttpConfig{
		Method: "POST",
		URL:    "http://test.com",
	}
	monitor := &HttpMonitor{
		Config: config,
	}
	assert.Equal(t, &config, monitor.GetConfig())
}
