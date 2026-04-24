package monitors

import (
	"testing"

	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
)

func createTestBaseMonitor() BaseMonitor {
	name := "Test BaseMonitor"
	return BaseMonitor{
		Slug:        util.SlugFromString(name),
		Name:        name,
		Description: "Test Description",
		Interval:    60,
		Type:        httpType,
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

func TestBaseMonitorGenerateSlug(t *testing.T) {
	t.Run("Slug is empty", func(t *testing.T) {
		monitor := createTestBaseMonitor()
		monitor.Slug = ""
		monitor.GenerateSlug()
		assert.NotEmpty(t, monitor.Slug, "Generated slug should not be empty")
	})

	t.Run("Slug is already set", func(t *testing.T) {
		monitor := createTestBaseMonitor()
		originalSlug := monitor.Slug
		monitor.GenerateSlug()
		assert.Equal(t, originalSlug, monitor.Slug, "Slug should remain unchanged if already set")
	})
}
