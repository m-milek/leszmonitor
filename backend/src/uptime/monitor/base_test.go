package monitors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func createTestBaseMonitor() BaseMonitor {
	return BaseMonitor{
		Id:          generateId(),
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

func TestGenerateMonitorId(t *testing.T) {
	t.Run("Valid ID Generation", func(t *testing.T) {
		id := generateId()
		assert.NotEmpty(t, id, "Generated ID should not be empty")
	})

	t.Run("ID Uniqueness", func(t *testing.T) {
		id1 := generateId()
		id2 := generateId()
		assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
	})

	t.Run("Id generated multiple times is different", func(t *testing.T) {
		ids := make(map[string]bool)
		for i := 0; i < 100; i++ {
			id := generateId()
			assert.NotEmpty(t, id, "Generated ID should not be empty")
			_, exists := ids[id]
			assert.False(t, exists, "Generated ID should be unique")
			ids[id] = true
		}
	})
}
