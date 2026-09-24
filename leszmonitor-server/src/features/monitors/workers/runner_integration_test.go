package workers

import (
	"context"
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_ProbeRunner_CheckExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, realDB, monitor := setupFullDB(t)

	runChannel := monitors.MonitorExecuteChannel.Subscribe()
	defer monitors.MonitorExecuteChannel.Unsubscribe(runChannel)

	runner := &monitorRunner{
		monitor:    *monitor,
		db:         realDB,
		cancel:     cancel,
		updates:    make(chan monitors.Monitor, 1),
		onExit:     func() {},
		baseLogger: log.New(),
	}

	go runner.run(ctx)

	// We should receive an execute message since the interval is 1s
	select {
	case msg := <-runChannel:
		assert.Equal(t, monitor.ID, msg.Monitor.ID)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for monitor execute message")
	}
}

func TestIntegration_ProbeRunner_SelfTermination(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, realDB, monitor := setupFullDB(t)

	// Make monitor invalid to trigger self-termination during initialization
	monitor.Interval = 0

	exitCalled := make(chan bool, 1)

	runner := &monitorRunner{
		monitor:    *monitor,
		db:         realDB,
		cancel:     cancel,
		updates:    make(chan monitors.Monitor, 1),
		onExit:     func() { exitCalled <- true },
		baseLogger: log.New(),
	}

	go runner.run(ctx)

	select {
	case <-exitCalled:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for self-termination")
	}
}
