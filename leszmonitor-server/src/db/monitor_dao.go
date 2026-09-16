package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/models/monitors"
	platformdb "github.com/m-milek/leszmonitor/platform/db"
)

type IMonitorDAO interface {
	GetMonitorByID(ctx context.Context, id uuid.UUID) (*monitors.Monitor, error)
	GetMonitorBySlug(ctx context.Context, slug string) (*monitors.Monitor, error)
	GetAllMonitors(ctx context.Context) ([]monitors.Monitor, error)
	DeleteMonitorByID(ctx context.Context, monitorID uuid.UUID) (*uuid.UUID, error)
	InsertMonitor(ctx context.Context, monitor monitors.Monitor) (*monitors.Monitor, error)
	UpdateMonitor(ctx context.Context, newMonitor monitors.Monitor) (any, error)
}

type monitorDAO struct {
	baseDAO
}

func newMonitorDAO(dao baseDAO) IMonitorDAO {
	return &monitorDAO{
		baseDAO: dao,
	}
}

func (r *monitorDAO) GetMonitorBySlug(
	ctx context.Context,
	slug string,
) (*monitors.Monitor, error) {
	return platformdb.Wrap(ctx, "GetMonitorBySlug", func() (*monitors.Monitor, error) {
		var monitor monitors.Monitor
		err := sqlx.GetContext(
			ctx,
			r.pool,
			&monitor,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.run_state, m.config, m.owner_id, m.created_at, m.updated_at
			 FROM monitors m
			 WHERE m.slug = $1`,
			slug,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		if err = r.attachTags(ctx, &monitor); err != nil {
			return nil, err
		}

		return &monitor, nil
	})
}

func (r *monitorDAO) GetMonitorByID(ctx context.Context, id uuid.UUID) (*monitors.Monitor, error) {
	return platformdb.Wrap(ctx, "GetMonitorByID", func() (*monitors.Monitor, error) {
		var monitor monitors.Monitor
		err := sqlx.GetContext(
			ctx,
			r.pool,
			&monitor,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.run_state, m.config, m.owner_id, m.created_at, m.updated_at
			 FROM monitors m
			 WHERE m.id = $1`,
			id,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		if err = r.attachTags(ctx, &monitor); err != nil {
			return nil, err
		}

		return &monitor, nil
	})
}

func (r *monitorDAO) GetAllMonitors(ctx context.Context) ([]monitors.Monitor, error) {
	return platformdb.Wrap(ctx, "GetAllMonitors", func() ([]monitors.Monitor, error) {
		var allMonitors []monitors.Monitor
		err := sqlx.SelectContext(
			ctx,
			r.pool,
			&allMonitors,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.run_state, m.config, m.owner_id, m.created_at, m.updated_at
			 FROM monitors m`,
		)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if allMonitors == nil {
			allMonitors = []monitors.Monitor{}
		}

		if err = r.attachTagsToAll(ctx, allMonitors); err != nil {
			return nil, err
		}

		return allMonitors, nil
	})
}

func (r *monitorDAO) DeleteMonitorByID(ctx context.Context, monitorID uuid.UUID) (*uuid.UUID, error) {
	return platformdb.Wrap(ctx, "DeleteMonitor", func() (*uuid.UUID, error) {
		var id uuid.UUID
		err := r.pool.QueryRowxContext(ctx, `DELETE FROM monitors WHERE id = $1 RETURNING id`, monitorID.String()).
			Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return &id, nil
	})
}

func (r *monitorDAO) attachTags(ctx context.Context, monitor *monitors.Monitor) error {
	tagIDs := []uuid.UUID{}
	err := sqlx.SelectContext(
		ctx,
		r.pool,
		&tagIDs,
		`SELECT tag_id FROM monitor_tags WHERE monitor_id = $1 ORDER BY tag_id`,
		monitor.ID,
	)
	if err != nil {
		return err
	}

	monitor.TagIDs = tagIDs
	return nil
}

func (r *monitorDAO) attachTagsToAll(ctx context.Context, allMonitors []monitors.Monitor) error {
	var rows []struct {
		MonitorID uuid.UUID `db:"monitor_id"`
		TagID     uuid.UUID `db:"tag_id"`
	}
	err := sqlx.SelectContext(
		ctx,
		r.pool,
		&rows,
		`SELECT monitor_id, tag_id FROM monitor_tags ORDER BY tag_id`,
	)
	if err != nil {
		return err
	}

	byMonitorID := make(map[uuid.UUID][]uuid.UUID, len(allMonitors))
	for _, row := range rows {
		byMonitorID[row.MonitorID] = append(byMonitorID[row.MonitorID], row.TagID)
	}

	for i := range allMonitors {
		tagIDs := byMonitorID[allMonitors[i].ID]
		if tagIDs == nil {
			tagIDs = []uuid.UUID{}
		}
		allMonitors[i].TagIDs = tagIDs
	}

	return nil
}

func (r *monitorDAO) replaceMonitorTags(
	ctx context.Context,
	monitorID uuid.UUID,
	tagIDs []uuid.UUID,
) error {
	_, err := r.pool.ExecContext(ctx, `DELETE FROM monitor_tags WHERE monitor_id = $1`, monitorID)
	if err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		_, err = r.pool.ExecContext(
			ctx,
			`INSERT OR IGNORE INTO monitor_tags (monitor_id, tag_id) VALUES ($1, $2)`,
			monitorID,
			tagID,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// InsertMonitor adds a new monitor to the database and returns the created monitor.
func (r *monitorDAO) InsertMonitor(ctx context.Context, monitor monitors.Monitor) (*monitors.Monitor, error) {
	return platformdb.Wrap(ctx, "InsertMonitor", func() (*monitors.Monitor, error) {
		id := monitor.ID
		if id == uuid.Nil {
			id = uuid.New()
		}

		_, err := r.pool.ExecContext(
			ctx,
			`INSERT INTO monitors (id, slug, name, description, interval, kind, result_retention_seconds, run_state, config, owner_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			id,
			monitor.Slug,
			monitor.Name,
			monitor.Description,
			monitor.Interval,
			monitor.Type,
			monitor.ResultRetentionSeconds,
			monitor.RunState,
			monitor.ProbeConfig,
			monitor.OwnerID,
		)
		if err != nil {
			if platformdb.IsUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		if err = r.replaceMonitorTags(ctx, id, monitor.TagIDs); err != nil {
			return nil, err
		}

		return r.GetMonitorByID(ctx, id)
	})
}

func (r *monitorDAO) UpdateMonitor(ctx context.Context, newMonitor monitors.Monitor) (any, error) {
	return platformdb.Wrap(ctx, "UpdateMonitor", func() (any, error) {
		res, err := r.pool.ExecContext(ctx,
			`UPDATE monitors
			SET slug=$1, name=$2, description=$3, interval=$4, kind=$5, run_state=$6, config=$7
			WHERE id=$8`,
			newMonitor.Slug,
			newMonitor.Name,
			newMonitor.Description,
			newMonitor.Interval,
			newMonitor.Type,
			newMonitor.RunState,
			newMonitor.ProbeConfig,
			newMonitor.ID,
		)
		if err != nil {
			if platformdb.IsUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if rowsAffected == 0 {
			return nil, ErrNotFound
		}

		if err = r.replaceMonitorTags(ctx, newMonitor.ID, newMonitor.TagIDs); err != nil {
			return nil, err
		}

		return nil, nil
	})
}
