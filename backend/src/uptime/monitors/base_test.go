package monitors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func createTestBaseMonitor() BaseMonitor {
	return BaseMonitor{
		Id:          generateMonitorId(),
		Name:        "Test BaseMonitor",
		Description: "Test Description",
		Interval:    60,
		Type:        Http,
		OwnerId:     "test-owner-id",
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
