package monitors

import (
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func createTestBaseMonitor() BaseMonitor {
	name := "Test BaseMonitor"
	return BaseMonitor{
		Id:          util.IdFromString(name),
		Name:        name,
		Description: "Test Description",
		Interval:    60,
		Type:        Http,
		GroupId:     "test-group-id",
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
	t.Run("Id is empty", func(t *testing.T) {
		monitor := createTestBaseMonitor()
		monitor.Id = ""
		monitor.GenerateId()
		assert.NotEmpty(t, monitor.Id, "Generated ID should not be empty")
	})

	t.Run("Id is already set", func(t *testing.T) {
		monitor := createTestBaseMonitor()
		originalId := monitor.Id
		monitor.GenerateId()
		assert.Equal(t, originalId, monitor.Id, "ID should remain unchanged if already set")
	})
}
