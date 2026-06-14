package probes

import (
	"context"
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ProbeRunner_CheckExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, realDB, monitor := setupFullDB(t)

	runChannel := events.MonitorRunChannel.Subscribe()
	defer events.MonitorRunChannel.Unsubscribe(runChannel)

	runner := &probeRunner{
		monitor:    *monitor,
		db:         realDB,
		cancel:     cancel,
		updates:    make(chan monitors.Monitor, 1),
		onExit:     func() {},
		baseLogger: log.New(),
	}

	go runner.run(ctx)

	// We should receive a run result since the interval is 1s
	select {
	case msg := <-runChannel:
		assert.Equal(t, monitor.ID, msg.Monitor.ID)
		assert.NotNil(t, msg.Result)
	case <-time.After(3 * time.Second):
		t.Fatal("Timeout waiting for monitor run result")
	}

	// Verify that the result was inserted in the DB
	results, err := realDB.MonitorResults().GetMonitorResultsByMonitorID(ctx, monitor.ID.String(), &util.Pagination{Page: 1, PerPage: 10})
	require.NoError(t, err)
	assert.GreaterOrEqual(t, len(results), 1)
}

func TestIntegration_ProbeRunner_SelfTermination(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, realDB, monitor := setupFullDB(t)

	// Make monitor invalid to trigger self-termination during initialization
	monitor.Interval = 0

	exitCalled := make(chan bool, 1)

	runner := &probeRunner{
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
