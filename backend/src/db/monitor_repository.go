package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/log"
	monitors "github.com/m-milek/leszmonitor/models/monitors"
)

type IMonitorRepository interface {
	GetMonitorsByProjectID(ctx context.Context, projectID uuid.UUID) ([]monitors.IConcreteMonitor, error)
	GetMonitorByID(ctx context.Context, id string) (monitors.IConcreteMonitor, error)
	GetAllMonitors(ctx context.Context) ([]monitors.IConcreteMonitor, error)
	DeleteMonitorBySlug(ctx context.Context, slug string) (*uuid.UUID, error)
	InsertMonitor(ctx context.Context, monitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error)
	UpdateMonitor(ctx context.Context, newMonitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error)
}

type monitorRepository struct {
	baseRepository
}

func newMonitorRepository(repository baseRepository) IMonitorRepository {
	return &monitorRepository{
		baseRepository: repository,
	}
}

func mapRowsToMonitors(rows *sqlx.Rows) ([]monitors.IConcreteMonitor, error) {
	defer rows.Close()
	var allMonitors []monitors.IConcreteMonitor
	for rows.Next() {
		var config []byte
		var b monitors.BaseMonitor

		err := rows.Scan(&b.ID, &b.Slug, &b.Name, &b.Description, &b.Interval, &b.Type, &config, &b.CreatedAt, &b.UpdatedAt, &b.ProjectSlug)
		if err != nil {
			return nil, err
		}

		parsedConfig, err := monitors.UnmarshalConfigFromBytes(b.Type, config)
		if err != nil {
			return nil, fmt.Errorf("failed to create monitor from config: %w", err)
		}

		monitor, err := monitors.NewConcreteMonitor(b, parsedConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to create monitor: %w", err)
		}
		allMonitors = append(allMonitors, monitor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if allMonitors == nil {
		allMonitors = []monitors.IConcreteMonitor{}
	}
	return allMonitors, nil
}

func mapRowToMonitor(row *sqlx.Row) (monitors.IConcreteMonitor, error) {
	var config []byte
	var b monitors.BaseMonitor

	err := row.Scan(&b.ID, &b.Slug, &b.Name, &b.Description, &b.Interval, &b.Type, &config, &b.CreatedAt, &b.UpdatedAt, &b.ProjectSlug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	parsedConfig, err := monitors.UnmarshalConfigFromBytes(b.Type, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitor from config: %w", err)
	}

	monitor, err := monitors.NewConcreteMonitor(b, parsedConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitor: %w", err)
	}
	return monitor, nil
}

func (r *monitorRepository) GetMonitorsByProjectID(ctx context.Context, projectID uuid.UUID) ([]monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetMonitorsByProjectID", func() ([]monitors.IConcreteMonitor, error) {
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

func (r *monitorRepository) GetMonitorByID(ctx context.Context, id string) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetMonitorByID", func() (monitors.IConcreteMonitor, error) {
		row := r.pool.QueryRowxContext(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.slug = $1`,
			id)
		return mapRowToMonitor(row)
	})
}

func (r *monitorRepository) GetAllMonitors(ctx context.Context) ([]monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetAllMonitors", func() ([]monitors.IConcreteMonitor, error) {
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
func (r *monitorRepository) InsertMonitor(ctx context.Context, monitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "InsertMonitor", func() (monitors.IConcreteMonitor, error) {
		id := monitor.GetID()
		if id == uuid.Nil {
			id = uuid.New() // manually inject UUID
		}

		configBytes, err := json.Marshal(monitor.GetConfig())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal monitor config: %w", err)
		}

		var projectID uuid.UUID
		if err := r.pool.QueryRowxContext(ctx,
			`SELECT id FROM projects WHERE slug = $1`,
			monitor.GetProjectSlug(),
		).Scan(&projectID); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		_, err = r.pool.ExecContext(ctx,
			`INSERT INTO monitors (id, slug, project_id, name, description, interval, kind, config)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			id,
			monitor.GetSlug(),
			projectID,
			monitor.GetName(),
			monitor.GetDescription(),
			int(monitor.GetInterval().Seconds()),
			string(monitor.GetType()),
			configBytes,
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

func (r *monitorRepository) UpdateMonitor(ctx context.Context, newMonitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "UpdateMonitor", func() (monitors.IConcreteMonitor, error) {
		configBytes, err := json.Marshal(newMonitor.GetConfig())
		if err != nil {
			return nil, fmt.Errorf("failed to marshal monitor config: %w", err)
		}

		row := r.pool.QueryRowxContext(ctx,
			`UPDATE monitors m
			SET slug=$1, project_id=(SELECT p.id FROM projects p WHERE p.slug=$2), name=$3, description=$4, interval=$5, kind=$6, config=$7
			WHERE id=$8
			RETURNING
				m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at,
				$2 AS project_slug`,
			newMonitor.GetSlug(),
			newMonitor.GetProjectSlug(),
			newMonitor.GetName(),
			newMonitor.GetDescription(),
			int(newMonitor.GetInterval().Seconds()),
			string(newMonitor.GetType()),
			configBytes,
			newMonitor.GetID(),
		)

		updatedMonitor, err := mapRowToMonitor(row)
		if err != nil {
			return nil, err
		}
		return updatedMonitor, nil
	})
}
