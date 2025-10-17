package monitors

import (
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createTestBaseMonitor() BaseMonitor {
	name := "Test BaseMonitor"
	return BaseMonitor{
		DisplayId:   util.IdFromString(name),
		Name:        name,
		Description: "Test Description",
		Interval:    60,
		Type:        Http,
	}
}

func TestBaseMonitorValidateSuccess(t *testing.T) {
	monitor := createTestBaseMonitor()
	err := monitor.validateBase()
	assert.NoError(t, err)
}

func TestBaseMonitorValidateEmptyName(t *testing.T) {
	monitor := createTestBaseMonitor()
	monitor.Name = ""
	err := monitor.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "name cannot be empty")
}

func TestBaseMonitorValidateZeroInterval(t *testing.T) {
	monitor := createTestBaseMonitor()
	monitor.Interval = 0
	err := monitor.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interval must be greater than zero")
}

func TestBaseMonitorValidateNegativeInterval(t *testing.T) {
	monitor := createTestBaseMonitor()
	monitor.Interval = -10
	err := monitor.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "interval must be greater than zero")
}

func TestBaseMonitorValidateEmptyType(t *testing.T) {
	monitor := createTestBaseMonitor()
	monitor.Type = ""
	err := monitor.Validate()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "type cannot be empty")
}

func TestBaseMonitorGenerateId(t *testing.T) {
	t.Run("DisplayId is empty", func(t *testing.T) {
		monitor := createTestBaseMonitor()
		monitor.DisplayId = ""
		monitor.GenerateDisplayId()
		assert.NotEmpty(t, monitor.DisplayId, "Generated DisplayId should not be empty")
	})

	t.Run("DisplayId is already set", func(t *testing.T) {
		monitor := createTestBaseMonitor()
		originalId := monitor.DisplayId
		monitor.GenerateDisplayId()
		assert.Equal(t, originalId, monitor.DisplayId, "DisplayId should remain unchanged if already set")
	})
}
