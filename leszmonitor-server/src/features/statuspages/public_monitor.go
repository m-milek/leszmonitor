package statuspages

import (
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/stats"
	"github.com/m-milek/leszmonitor/features/tags"
)

// TODO: aggregate results into time buckets (e.g. daily) instead of returning the last N raw results.
const publicMonitorResultsLimit = 50

type PublicMonitor struct {
	ID      uuid.UUID             `json:"id"`
	Name    string                `json:"name"`
	Status  kind.MonitorStatus    `json:"status"`
	Stats   stats.MonitorStats    `json:"stats"`
	Tags    []tags.Tag            `json:"tags"`
	Results []PublicMonitorResult `json:"results"`
}

type PublicMonitorResult struct {
	Status    kind.MonitorStatus `json:"status"    db:"status"`
	CreatedAt time.Time          `json:"createdAt" db:"created_at"`
}

type publicMonitorRow struct {
	ID       uuid.UUID                `db:"id"`
	Name     string                   `db:"name"`
	RunState monitors.MonitorRunState `db:"run_state"`
}

func newPublicMonitor(row publicMonitorRow, monitorTags []tags.Tag, results []PublicMonitorResult) PublicMonitor {
	if monitorTags == nil {
		monitorTags = []tags.Tag{}
	}
	if results == nil {
		results = []PublicMonitorResult{}
	}

	status := kind.MonitorStatusUnknown
	if len(results) > 0 {
		status = results[0].Status
	}
	if row.RunState == monitors.MonitorRunStateStopped {
		status = kind.MonitorStatusPaused
	}

	return PublicMonitor{
		ID:      row.ID,
		Name:    row.Name,
		Status:  status,
		Tags:    monitorTags,
		Results: results,
	}
}
