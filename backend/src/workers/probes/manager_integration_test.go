package probes

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ProbesManager_Lifecycle(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, realDB, monitor := setupFullDB(t)

	manager := NewProbesManager(realDB)
	go manager.Run(ctx)

	// Wait for DB loading to finish
	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, 1, manager.ActiveCount())

	// Test Edited
	monitor.Interval = 10
	_, err := realDB.Monitors().UpdateMonitor(ctx, *monitor)
	require.NoError(t, err)

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		Status: monitors.Edited,
		ID:     monitor.ID,
	})
	time.Sleep(200 * time.Millisecond)

	runner := manager.get(monitor.ID)
	require.NotNil(t, runner)
	assert.Equal(t, 10, runner.monitor.Interval)

	// Test Deleted
	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
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

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		Status:  monitors.Created,
		Monitor: &newMonitor,
	})
	time.Sleep(200 * time.Millisecond)

	assert.Equal(t, 1, manager.ActiveCount())
}
