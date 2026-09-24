package workers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Manager_Lifecycle(t *testing.T) {
	ctx := t.Context()

	_, realDB, monitor := setupFullDB(t)

	manager := NewMonitorScheduler(realDB)
	go manager.Run(ctx)

	// Wait for DB loading to finish
	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, 1, manager.ActiveCount())

	// Test Edited
	monitor.Interval = 10
	_, err := monitors.NewMonitorDAO(realDB.Querier()).UpdateMonitor(ctx, *monitor)
	require.NoError(t, err)

	monitors.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		Status: monitors.Edited,
		ID:     monitor.ID,
	})
	time.Sleep(200 * time.Millisecond)

	runner := manager.get(monitor.ID)
	require.NotNil(t, runner)
	assert.Equal(t, 10, runner.monitor.Interval)

	// Test Deleted
	monitors.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		Status: monitors.Deleted,
		ID:     monitor.ID,
	})
	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, 0, manager.ActiveCount())

	// Test Created
	newMonitorID := uuid.New()
	newMonitor := *monitor
	newMonitor.ID = newMonitorID
	newMonitor.Name = "Brand New Monitor"
	newMonitor.GenerateSlug()

	monitors.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		Status:  monitors.Created,
		Monitor: &newMonitor,
	})
	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, 1, manager.ActiveCount())
}
