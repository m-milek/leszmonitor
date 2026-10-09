package workers

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_PushWatcher_CheckMonitor(t *testing.T) {
	base := time.Now().UTC().Add(-time.Hour)
	monitor := monitors.Monitor{
		Timestamps:  util.Timestamps{CreatedAt: base},
		ID:          uuid.New(),
		Interval:    60,
		Type:        kind.PushConfigType,
		ProbeConfig: `{"gracePeriodSeconds":30}`,
		RunState:    monitors.MonitorRunStateActive,
	}

	t.Run("Broadcasts a missed heartbeat after interval and grace period", func(t *testing.T) {
		ctx, realDB := setupDB(t)
		watcher := &PushWatcher{db: realDB, startedAt: base}

		runChannel := monitors.MonitorRunChannel.Subscribe()
		defer monitors.MonitorRunChannel.Unsubscribe(runChannel)

		watcher.checkMonitor(ctx, monitor, base.Add(91*time.Second))

		select {
		case msg := <-runChannel:
			assert.Equal(t, monitor.ID, msg.Monitor.ID)
			assert.Equal(t, kind.MonitorStatusDown, msg.Result.GetStatus())
			require.Len(t, msg.Result.GetFailures(), 1)
			assert.Equal(t, results.FailureReasonPushMissedHeartbeat, msg.Result.GetFailures()[0].Reason)
		default:
			t.Fatal("Expected a monitor run message")
		}
	})

	t.Run("Does nothing within interval and grace period", func(t *testing.T) {
		ctx, realDB := setupDB(t)
		watcher := &PushWatcher{db: realDB, startedAt: base}

		runChannel := monitors.MonitorRunChannel.Subscribe()
		defer monitors.MonitorRunChannel.Unsubscribe(runChannel)

		watcher.checkMonitor(ctx, monitor, base.Add(89*time.Second))

		select {
		case <-runChannel:
			t.Fatal("Expected no monitor run message")
		default:
		}
	})
}

func TestPushWatcher_Baseline(t *testing.T) {
	base := time.Now().UTC()
	at := func(minutes int) time.Time { return base.Add(time.Duration(minutes) * time.Minute) }
	resultAt := func(createdAt time.Time) results.IMonitorResult {
		result := results.NewMonitorResult(uuid.New(), kind.PushConfigType, kind.MonitorStatusUp, false, nil, nil)
		result.CreatedAt = createdAt
		return &result
	}

	tests := []struct {
		name         string
		startedAt    time.Time
		createdAt    time.Time
		updatedAt    *time.Time
		latestResult results.IMonitorResult
		expected     time.Time
	}{
		{name: "Watcher start is latest", startedAt: at(3), createdAt: at(1), updatedAt: new(at(2)), latestResult: resultAt(at(0)), expected: at(3)},
		{name: "Monitor creation is latest", startedAt: at(0), createdAt: at(3), expected: at(3)},
		{name: "Monitor update is latest", startedAt: at(0), createdAt: at(1), updatedAt: new(at(3)), latestResult: resultAt(at(2)), expected: at(3)},
		{name: "Latest result is latest", startedAt: at(0), createdAt: at(1), updatedAt: new(at(2)), latestResult: resultAt(at(3)), expected: at(3)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watcher := &PushWatcher{startedAt: tt.startedAt}
			monitor := monitors.Monitor{Timestamps: util.Timestamps{CreatedAt: tt.createdAt, UpdatedAt: tt.updatedAt}}

			assert.Equal(t, tt.expected, watcher.baseline(monitor, tt.latestResult))
		})
	}
}
