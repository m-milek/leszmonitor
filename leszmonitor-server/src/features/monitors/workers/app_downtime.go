package workers

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/statuschange"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/downtime"
	"github.com/pkg/errors"
)

// RecordDowntimeStatusChanges marks every monitor as unknown from the start of an application downtime.
func RecordDowntimeStatusChanges(ctx context.Context, q db.Querier, appDowntime downtime.AppDowntime) error {
	allMonitors, err := monitors.NewMonitorDAO(q).GetAllMonitors(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to retrieve monitors")
	}

	statusChangeDAO := statuschange.NewMonitorStatusChangeDAO(q)
	for _, monitor := range allMonitors {
		latest, err := statusChangeDAO.GetLatestStatusChangeByMonitorID(ctx, monitor.ID, appDowntime.StartedAt)
		if errors.Is(err, db.ErrNotFound) {
			continue
		}
		if err != nil {
			return errors.Wrap(err, "failed to retrieve latest status change")
		}
		if latest.NextStatus == string(kind.MonitorStatusUnknown) {
			continue
		}

		_, err = statusChangeDAO.InsertStatusChange(ctx, statuschange.MonitorStatusChange{
			ID:             uuid.New(),
			MonitorID:      monitor.ID,
			PreviousStatus: latest.NextStatus,
			NextStatus:     string(kind.MonitorStatusUnknown),
			CreatedAt:      appDowntime.StartedAt,
		})
		if err != nil {
			return errors.Wrap(err, "failed to insert status change")
		}
	}

	return nil
}
