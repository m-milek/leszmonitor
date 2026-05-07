package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/log"
	monitors "github.com/m-milek/leszmonitor/models/monitors"
)

type IMonitorRepository interface {
	GetMonitorsByProjectID(ctx context.Context, projectID uuid.UUID) ([]monitors.Monitor, error)
	GetMonitorByID(ctx context.Context, id uuid.UUID) (*monitors.Monitor, error)
	GetMonitorBySlug(ctx context.Context, slug string) (*monitors.Monitor, error)
	GetAllMonitors(ctx context.Context) ([]monitors.Monitor, error)
	DeleteMonitorBySlug(ctx context.Context, slug string) (*uuid.UUID, error)
	InsertMonitor(ctx context.Context, monitor monitors.Monitor) (*monitors.Monitor, error)
	UpdateMonitor(ctx context.Context, newMonitor monitors.Monitor) (interface{}, error)
}

type monitorRepository struct {
	baseRepository
}

func newMonitorRepository(repository baseRepository) IMonitorRepository {
	return &monitorRepository{
		baseRepository: repository,
	}
}

func mapRowsToMonitors(rows *sqlx.Rows) ([]monitors.Monitor, error) {
	defer rows.Close()
	var allMonitors []monitors.Monitor
	for rows.Next() {
		var b monitors.Monitor

		err := rows.Scan(&b.ID, &b.Slug, &b.Name, &b.Description, &b.Interval, &b.Type, &b.ProbeConfig, &b.CreatedAt, &b.UpdatedAt, &b.ProjectSlug)
		if err != nil {
			return nil, err
		}

		allMonitors = append(allMonitors, b)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if allMonitors == nil {
		allMonitors = []monitors.Monitor{}
	}
	return allMonitors, nil
}

func mapRowToMonitor(row *sqlx.Row) (*monitors.Monitor, error) {
	var b monitors.Monitor

	err := row.Scan(&b.ID, &b.Slug, &b.Name, &b.Description, &b.Interval, &b.Type, &b.ProbeConfig, &b.CreatedAt, &b.UpdatedAt, &b.ProjectSlug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &b, nil
}

func (r *monitorRepository) GetMonitorsByProjectID(ctx context.Context, projectID uuid.UUID) ([]monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorsByProjectID", func() ([]monitors.Monitor, error) {
		rows, err := r.pool.QueryxContext(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.project_id = $1`,
			projectID)
		if err != nil {
			return nil, err
		}
		mappedMonitors, err := mapRowsToMonitors(rows)
		if err != nil {
			log.Db.Error().Err(err).Msg("Error mapRowsToMonitors")
			return nil, err
		}
		return mappedMonitors, nil
	})
}

func (r *monitorRepository) GetMonitorBySlug(ctx context.Context, slug string) (*monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorBySlug", func() (*monitors.Monitor, error) {
		row := r.pool.QueryRowxContext(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.slug = $1`,
			slug)
		return mapRowToMonitor(row)
	})
}

func (r *monitorRepository) GetMonitorByID(ctx context.Context, id uuid.UUID) (*monitors.Monitor, error) {
	return dbWrap(ctx, "GetMonitorByID", func() (*monitors.Monitor, error) {
		row := r.pool.QueryRowxContext(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.id = $1`,
			id)
		return mapRowToMonitor(row)
	})
}

func (r *monitorRepository) GetAllMonitors(ctx context.Context) ([]monitors.Monitor, error) {
	return dbWrap(ctx, "GetAllMonitors", func() ([]monitors.Monitor, error) {
		rows, err := r.pool.QueryxContext(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id`)
		if err != nil {
			return nil, err
		}
		return mapRowsToMonitors(rows)
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
			id = uuid.New() // manually inject UUID
		}

		var projectID uuid.UUID
		if err := r.pool.QueryRowxContext(ctx,
			`SELECT id FROM projects WHERE slug = $1`,
			monitor.ProjectSlug,
		).Scan(&projectID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		_, err := r.pool.ExecContext(ctx,
			`INSERT INTO monitors (id, slug, project_id, name, description, interval, kind, config)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id,
			monitor.Slug,
			projectID,
			monitor.Name,
			monitor.Description,
			monitor.Interval,
			monitor.Type,
			monitor.ProbeConfig,
		)
		if err != nil {
			if isUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		row := r.pool.QueryRowxContext(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.id = $1`,
			id,
		)
		created, err := mapRowToMonitor(row)
		if err != nil {
			return nil, err
		}
		return created, nil
	})
}

func (r *monitorRepository) UpdateMonitor(ctx context.Context, newMonitor monitors.Monitor) (any, error) {
	return dbWrap(ctx, "UpdateMonitor", func() (any, error) {
		res, err := r.pool.ExecContext(ctx,
			`UPDATE monitors
			SET slug=$1, project_id=(SELECT p.id FROM projects p WHERE p.slug=$2), name=$3, description=$4, interval=$5, kind=$6, config=$7
			WHERE id=$8`,
			newMonitor.Slug,
			newMonitor.ProjectSlug,
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
