package workers

import (
	"testing"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/statuschange"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_ProcessMonitorRunMessage_Thresholds(t *testing.T) {
	ctx, realDB, monitor := setupFullDB(t)
	monitor.FailureThreshold = 3
	monitor.RecoveryThreshold = 2

	base := time.Now().UTC().Add(-time.Minute)
	process := func(i int, status kind.MonitorStatus) kind.MonitorStatus {
		result := results.NewMonitorResult(monitor.ID, monitor.Type, status, false, nil, &results.HTTPResultDetails{})
		result.CreatedAt = base.Add(time.Duration(i) * time.Second)
		require.NoError(t, processMonitorRunMessage(ctx, realDB, monitors.MonitorRunMessage{Monitor: *monitor, Result: &result}))

		latest, err := statuschange.NewMonitorStatusChangeDAO(realDB.Querier()).
			GetLatestStatusChangeByMonitorID(ctx, monitor.ID, time.Now().UTC())
		require.NoError(t, err)
		return kind.MonitorStatus(latest.NextStatus)
	}

	steps := []struct {
		status   kind.MonitorStatus
		expected kind.MonitorStatus
	}{
		{kind.MonitorStatusUp, kind.MonitorStatusUp},
		{kind.MonitorStatusDown, kind.MonitorStatusUp},
		{kind.MonitorStatusDown, kind.MonitorStatusUp},
		{kind.MonitorStatusUp, kind.MonitorStatusUp},
		{kind.MonitorStatusDown, kind.MonitorStatusUp},
		{kind.MonitorStatusDown, kind.MonitorStatusUp},
		{kind.MonitorStatusDown, kind.MonitorStatusDown},
		{kind.MonitorStatusUp, kind.MonitorStatusDown},
		{kind.MonitorStatusUp, kind.MonitorStatusUp},
	}
	for i, step := range steps {
		assert.Equal(t, step.expected, process(i, step.status), "step %d", i)
	}
}
