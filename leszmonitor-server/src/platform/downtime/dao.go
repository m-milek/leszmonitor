package downtime

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/db"
)

type IAppDowntimeDAO interface {
	GetHeartbeat(ctx context.Context) (*HeartbeatRecord, error)
	InsertHeartbeat(ctx context.Context, beat *HeartbeatRecord) (*HeartbeatRecord, error)
	InsertDowntime(ctx context.Context, downtime *AppDowntime) (*AppDowntime, error)
	GetAllDowntimes(ctx context.Context, from, to time.Time) ([]AppDowntime, error)
}

type AppDowntimeDAO struct {
	pool db.Querier
}

func NewAppDowntimeDAO(pool db.Querier) IAppDowntimeDAO {
	return &AppDowntimeDAO{
		pool: pool,
	}
}

func (r *AppDowntimeDAO) GetHeartbeat(ctx context.Context) (*HeartbeatRecord, error) {
	return db.Wrap(ctx, "GetHeartbeat", func() (*HeartbeatRecord, error) {
		var beat HeartbeatRecord

		err := sqlx.GetContext(ctx, r.pool, &beat, `
			SELECT beat_at
			FROM heartbeat
			WHERE id = $1`,
			1,
		)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		return &beat, nil
	})
}

func (r *AppDowntimeDAO) InsertHeartbeat(ctx context.Context, beat *HeartbeatRecord) (*HeartbeatRecord, error) {
	return db.Wrap(ctx, "InsertHeartbeat", func() (*HeartbeatRecord, error) {
		_, err := r.pool.ExecContext(ctx, `
			INSERT INTO heartbeat (id, beat_at)
			VALUES ($1, $2)
			ON CONFLICT (id) DO UPDATE SET beat_at = excluded.beat_at`,
			1,
			beat.BeatAt,
		)

		if err != nil {
			return nil, err
		}

		return beat, nil
	})
}

func (r *AppDowntimeDAO) InsertDowntime(ctx context.Context, downtime *AppDowntime) (*AppDowntime, error) {
	return db.Wrap(ctx, "InsertAppDowntime", func() (*AppDowntime, error) {
		_, err := r.pool.ExecContext(ctx, `
			INSERT INTO app_downtime_windows (id, started_at, ended_at)
			VALUES ($1, $2, $3)`,
			downtime.ID,
			downtime.StartedAt,
			downtime.EndedAt,
		)

		if err != nil {
			return nil, err
		}

		return downtime, nil
	})
}

func (r *AppDowntimeDAO) GetAllDowntimes(ctx context.Context, from, to time.Time) ([]AppDowntime, error) {
	return db.Wrap(ctx, "GetAllAppDowntime", func() ([]AppDowntime, error) {
		downtimes := []AppDowntime{}

		err := sqlx.SelectContext(ctx, r.pool, &downtimes, `
			SELECT id, started_at, ended_at
			FROM app_downtime_windows
			WHERE ended_at > $1
			  AND started_at < $2
			ORDER BY started_at ASC`,
			from,
			to,
		)

		if err != nil {
			return nil, err
		}

		return downtimes, nil
	})
}
