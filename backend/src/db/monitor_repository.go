package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/log"
	monitors "github.com/m-milek/leszmonitor/models/monitors"
)

type IMonitorRepository interface {
	GetMonitorsByProjectID(ctx context.Context, projectID pgtype.UUID) ([]monitors.IConcreteMonitor, error)
	GetMonitorByID(ctx context.Context, id string) (monitors.IConcreteMonitor, error)
	GetAllMonitors(ctx context.Context) ([]monitors.IConcreteMonitor, error)
	DeleteMonitorBySlug(ctx context.Context, slug string) (*pgtype.UUID, error)
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

// monitorFromCollectableRow maps a pgx.CollectableRow to a monitors.IConcreteMonitor.
func monitorFromCollectableRow(row pgx.CollectableRow) (monitors.IConcreteMonitor, error) {
	var config []byte
	var b monitors.BaseMonitor

	err := row.Scan(&b.ID, &b.Slug, &b.Name, &b.Description, &b.Interval, &b.Type, &config, &b.CreatedAt, &b.UpdatedAt, &b.ProjectSlug)
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

	return monitor, nil
}

func (r *monitorRepository) GetMonitorsByProjectID(ctx context.Context, projectID pgtype.UUID) ([]monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetMonitorsByProjectID", func() ([]monitors.IConcreteMonitor, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.project_id = $1`,
			projectID)
		if err != nil {
			return nil, err
		}

		allMonitors, err := pgx.CollectRows(rows, monitorFromCollectableRow)
		if err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			log.Db.Error().Err(err).Msg("Error occurred while iterating over monitor rows")
			return nil, err
		}

		return allMonitors, nil
	})
}

func (r *monitorRepository) GetMonitorByID(ctx context.Context, id string) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetMonitorByID", func() (monitors.IConcreteMonitor, error) {
		row, err := r.pool.Query(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id
			 WHERE m.slug = $1`,
			id)
		if err != nil {
			return nil, err
		}

		monitor, err := pgx.CollectOneRow(row, monitorFromCollectableRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return monitor, nil
	})
}

func (r *monitorRepository) GetAllMonitors(ctx context.Context) ([]monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetAllMonitors", func() ([]monitors.IConcreteMonitor, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at, p.slug AS project_slug
			 FROM monitors m
			 JOIN projects p ON p.id = m.project_id`)
		if err != nil {
			return nil, err
		}
		return pgx.CollectRows(rows, monitorFromCollectableRow)
	})
}

func (r *monitorRepository) DeleteMonitorBySlug(ctx context.Context, slug string) (*pgtype.UUID, error) {
	return dbWrap(ctx, "DeleteMonitor", func() (*pgtype.UUID, error) {
		result := r.pool.QueryRow(ctx, `DELETE FROM monitors WHERE slug = $1 RETURNING id`, slug)

		var id pgtype.UUID
		err := result.Scan(&id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
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
		rows, err := r.pool.Query(ctx,
			`INSERT INTO monitors (slug, project_id, name, description, interval, kind, config)
			SELECT $1, p.id, $3, $4, $5, $6, $7
			FROM projects p WHERE p.slug = $2
			RETURNING
				id, slug, name, description, interval, kind, config, created_at, updated_at,
				$2::TEXT AS project_slug`,
			monitor.GetSlug(),
			monitor.GetProjectSlug(),
			monitor.GetName(),
			monitor.GetDescription(),
			int(monitor.GetInterval().Seconds()),
			string(monitor.GetType()),
			monitor.GetConfig(),
		)
		if err != nil {
			if pgErrIs(err, pgerrcode.UniqueViolation) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		created, err := pgx.CollectOneRow(rows, monitorFromCollectableRow)
		if err != nil {
			return nil, err
		}

		return created, nil
	})
}

func (r *monitorRepository) UpdateMonitor(ctx context.Context, newMonitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "UpdateMonitor", func() (monitors.IConcreteMonitor, error) {
		result, err := r.pool.Query(ctx,
			`UPDATE monitors m
			SET slug=$1, project_id=(SELECT p.id FROM projects p WHERE p.slug=$2), name=$3, description=$4, interval=$5, kind=$6, config=$7
			WHERE id=$8
			RETURNING
				m.id, m.slug, m.name, m.description, m.interval, m.kind, m.config, m.created_at, m.updated_at,
				$2::TEXT AS project_slug`,
			newMonitor.GetSlug(),
			newMonitor.GetProjectSlug(),
			newMonitor.GetName(),
			newMonitor.GetDescription(),
			int(newMonitor.GetInterval().Seconds()),
			string(newMonitor.GetType()),
			newMonitor.GetConfig(),
			newMonitor.GetID(),
		)
		if err != nil {
			return nil, err
		}
		updatedMonitor, err := pgx.CollectOneRow(result, monitorFromCollectableRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return updatedMonitor, nil
	})
}
