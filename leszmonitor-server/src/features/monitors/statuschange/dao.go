package statuschange

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/db"
)

type IMonitorStatusChangeDAO interface {
	InsertStatusChange(ctx context.Context, statusChange MonitorStatusChange) (any, error)
	GetStatusChangesByMonitorID(
		ctx context.Context,
		monitorID string,
		from time.Time,
		to time.Time,
	) ([]MonitorStatusChange, error)
	GetLatestStatusChangeByMonitorID(
		ctx context.Context,
		monitorID uuid.UUID,
		to time.Time,
	) (*MonitorStatusChange, error)
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
			INSERT INTO monitor_status_changes (id, monitor_id, caused_by_id, previous_status, next_status, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)`

		_, err := r.pool.ExecContext(ctx, query,
			statusChange.ID,
			statusChange.MonitorID,
			statusChange.CausedByID,
			statusChange.PreviousStatus,
			statusChange.NextStatus,
			statusChange.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		return nil, nil
	})
}

func (r *monitorStatusChangeDAO) GetLatestStatusChangeByMonitorID(
	ctx context.Context,
	monitorID uuid.UUID,
	to time.Time,
) (*MonitorStatusChange, error) {
	return db.Wrap(ctx, "GetLatestStatusChangeByMonitorID", func() (*MonitorStatusChange, error) {
		var statusChange MonitorStatusChange

		err := sqlx.GetContext(ctx, r.pool, &statusChange, `
			SELECT id, monitor_id, caused_by_id, previous_status, next_status, created_at
			FROM monitor_status_changes
			WHERE monitor_id = $1
			  AND created_at < $2
			ORDER BY created_at DESC
			LIMIT 1`,
			monitorID,
			to,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		return &statusChange, nil
	})
}

func (r *monitorStatusChangeDAO) GetStatusChangesByMonitorID(
	ctx context.Context,
	monitorID string,
	from time.Time,
	to time.Time,
) ([]MonitorStatusChange, error) {
	return db.Wrap(ctx, "GetStatusChangesByMonitorID", func() ([]MonitorStatusChange, error) {
		statusChanges := []MonitorStatusChange{}

		err := sqlx.SelectContext(ctx, r.pool, &statusChanges, `
			SELECT id, monitor_id, caused_by_id, previous_status, next_status, created_at
			FROM monitor_status_changes
			WHERE monitor_id = $1
			  AND created_at >= $2
			  AND created_at < $3
			ORDER BY created_at ASC`,
			monitorID,
			from,
			to,
		)

		if err != nil {
			return nil, err
		}

		return statusChanges, nil
	})
}
