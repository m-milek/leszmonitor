package monitors

import (
	"context"

	"github.com/m-milek/leszmonitor/platform/db"
)

type IMonitorStatusChangeDAO interface {
	InsertStatusChange(ctx context.Context, statusChange MonitorStatusChange) (any, error)
}

type monitorStatusChangeDAO struct {
	pool db.Querier
}

func NewMonitorStatusChangeDAO(pool db.Querier) IMonitorStatusChangeDAO {
	return &monitorStatusChangeDAO{
		pool: pool,
	}
}

func (r *monitorStatusChangeDAO) InsertStatusChange(ctx context.Context, statusChange MonitorStatusChange) (any, error) {
	return db.Wrap(ctx, "InsertStatusChange", func() (any, error) {
		query := `
			INSERT INTO monitor_status_changes (id, monitor_id, caused_by_id, previous_status, next_status)
			VALUES ($1, $2, $3, $4, $5)`

		_, err := r.pool.ExecContext(ctx, query,
			statusChange.ID,
			statusChange.MonitorID,
			statusChange.CausedByID,
			statusChange.PreviousStatus,
			statusChange.NextStatus,
		)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}
