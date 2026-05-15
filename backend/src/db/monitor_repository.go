package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	monitors "github.com/m-milek/leszmonitor/models/monitors"
)

type IMonitorRepository interface {
	GetMonitorsByProjectID(ctx context.Context, projectID uuid.UUID) ([]monitors.Monitor, error)
	GetMonitorByID(ctx context.Context, id uuid.UUID) (*monitors.Monitor, error)
	GetMonitorBySlug(ctx context.Context, slug string, projectID uuid.UUID) (*monitors.Monitor, error)
	GetAllMonitors(ctx context.Context) ([]monitors.Monitor, error)
	DeleteMonitorBySlug(ctx context.Context, slug string) (*uuid.UUID, error)
	InsertMonitor(ctx context.Context, monitor monitors.Monitor) (*monitors.Monitor, error)
	UpdateMonitor(ctx context.Context, newMonitor monitors.Monitor) (interface{}, error)
	GetMonitorBySlugByProject(ctx context.Context, slug string, id uuid.UUID) (*monitors.Monitor, error)
}

type monitorRepository struct {
	baseRepository
}

func newMonitorRepository(repository baseRepository) IMonitorRepository {
	return &monitorRepository{
		baseRepository: repository,
	}
}

func (r *monitorRepository) GetMonitorsByProjectID(ctx context.Context, projectID uuid.UUID) ([]monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorsByProjectID", func() ([]monitors.Monitor, error) {
		var allMonitors []monitors.Monitor
		err := r.pool.SelectContext(ctx, &allMonitors,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.config, m.created_at, m.updated_at, m.project_id
			 FROM monitors m
			 WHERE m.project_id = $1`,
			projectID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if allMonitors == nil {
			allMonitors = []monitors.Monitor{}
		}
		return allMonitors, nil
	})
}

func (r *monitorRepository) GetMonitorBySlug(ctx context.Context, slug string, projectID uuid.UUID) (*monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorBySlug", func() (*monitors.Monitor, error) {
		var monitor monitors.Monitor
		err := r.pool.GetContext(ctx, &monitor,
			`SELECT m.id, m.slug, m.project_id, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.config, m.created_at, m.updated_at
			 FROM monitors m
			 WHERE m.slug = $1 
			   AND m.project_id = $2`,
			slug,
			projectID,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &monitor, nil
	})
}

func (r *monitorRepository) GetMonitorByID(ctx context.Context, id uuid.UUID) (*monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorByID", func() (*monitors.Monitor, error) {
		var monitor monitors.Monitor
		err := r.pool.GetContext(ctx, &monitor,
			`SELECT m.id, m.slug, m.project_id, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.config, m.created_at, m.updated_at
			 FROM monitors m
			 WHERE m.id = $1`, id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &monitor, nil
	})
}

func (r *monitorRepository) GetAllMonitors(ctx context.Context) ([]monitors.Monitor, error) {
	return dbWrap(ctx, "GetAllMonitors", func() ([]monitors.Monitor, error) {
		var allMonitors []monitors.Monitor
		err := r.pool.SelectContext(ctx, &allMonitors,
			`SELECT m.id, m.slug, m.project_id, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.config, m.created_at, m.updated_at
			 FROM monitors m`)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if allMonitors == nil {
			allMonitors = []monitors.Monitor{}
		}
		return allMonitors, nil
	})
}

func (r *monitorRepository) DeleteMonitorBySlug(ctx context.Context, slug string) (*uuid.UUID, error) {
	return dbWrap(ctx, "DeleteMonitor", func() (*uuid.UUID, error) {
		var id uuid.UUID
		err := r.pool.QueryRowxContext(ctx, `DELETE FROM monitors WHERE slug = $1 RETURNING id`, slug).Scan(&id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return &id, nil
	})
}

// InsertMonitor adds a new monitor to the database and returns the created monitor.
func (r *monitorRepository) InsertMonitor(ctx context.Context, monitor monitors.Monitor) (*monitors.Monitor, error) {
	return dbWrap(ctx, "InsertMonitor", func() (*monitors.Monitor, error) {
		id := monitor.ID
		if id == uuid.Nil {
			id = uuid.New()
		}

		_, err := r.pool.ExecContext(ctx,
			`INSERT INTO monitors (id, slug, project_id, name, description, interval, kind, result_retention_seconds,config)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			id,
			monitor.Slug,
			monitor.ProjectID,
			monitor.Name,
			monitor.Description,
			monitor.Interval,
			monitor.Type,
			monitor.ResultRetentionSeconds,
			monitor.ProbeConfig,
		)
		if err != nil {
			if isUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		return r.GetMonitorByID(ctx, id)
	})
}

func (r *monitorRepository) UpdateMonitor(ctx context.Context, newMonitor monitors.Monitor) (any, error) {
	return dbWrap(ctx, "UpdateMonitor", func() (any, error) {
		res, err := r.pool.ExecContext(ctx,
			`UPDATE monitors
			SET slug=$1, project_id=(SELECT p.id FROM projects p WHERE p.slug=$2), name=$3, description=$4, interval=$5, kind=$6, config=$7
			WHERE id=$8`,
			newMonitor.Slug,
			newMonitor.ProjectID,
			newMonitor.Name,
			newMonitor.Description,
			newMonitor.Interval,
			newMonitor.Type,
			newMonitor.ProbeConfig,
			newMonitor.ID,
		)
		if err != nil {
			if isUniqueViolation(err) {
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

		return nil, nil
	})
}

func (r *monitorRepository) GetMonitorBySlugByProject(ctx context.Context, slug string, id uuid.UUID) (*monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorBySlugByProject", func() (*monitors.Monitor, error) {
		var monitor monitors.Monitor
		err := r.pool.GetContext(ctx, &monitor,
			`SELECT m.id, m.slug, m.project_id, m.name, m.description, m.interval, m.kind, m.result_retention_seconds, m.config, m.created_at, m.updated_at
			 FROM monitors m
			 WHERE m.slug = $1 
			   AND m.project_id = $2`,
			slug,
			id,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &monitor, nil
	})
}
