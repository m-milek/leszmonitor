package workers

import (
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_MonitorExecutor_Execute(t *testing.T) {
	ctx := t.Context()

	_, realDB, monitor := setupFullDB(t)

	runChannel := monitors.MonitorRunChannel.Subscribe()
	defer monitors.MonitorRunChannel.Unsubscribe(runChannel)

	executor := NewMonitorExecutor(realDB)
	go executor.Run(ctx)

	time.Sleep(100 * time.Millisecond)

	monitors.MonitorExecuteChannel.Broadcast(monitors.MonitorExecuteMessage{Monitor: *monitor})

	select {
	case msg := <-runChannel:
		assert.Equal(t, monitor.ID, msg.Monitor.ID)
		assert.NotNil(t, msg.Result)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for monitor run result")
	}
}
