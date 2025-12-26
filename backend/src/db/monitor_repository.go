package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/logging"
	monitors "github.com/m-milek/leszmonitor/uptime/monitor"
)

type IMonitorRepository interface {
	GetMonitorsByTeamID(ctx context.Context, teamID pgtype.UUID) ([]monitors.IConcreteMonitor, error)
	GetMonitorByID(ctx context.Context, id string) (monitors.IConcreteMonitor, error)
	GetAllMonitors(ctx context.Context) ([]monitors.IConcreteMonitor, error)
	DeleteMonitorByDisplayID(ctx context.Context, displayID string) (*pgtype.UUID, error)
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

	err := row.Scan(&b.ID, &b.DisplayID, &b.TeamID, &b.GroupID, &b.Name, &b.Description, &b.Interval, &b.Type, &config, &b.CreatedAt, &b.UpdatedAt)
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

func (r *monitorRepository) GetMonitorsByTeamID(ctx context.Context, teamID pgtype.UUID) ([]monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetMonitorsByTeamID", func() ([]monitors.IConcreteMonitor, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT * FROM monitors WHERE team_id=$1`,
			teamID)

		if err != nil {
			return nil, err
		}
		var allMonitors []monitors.IConcreteMonitor
		allMonitors, err = pgx.CollectRows(rows, monitorFromCollectableRow)

		if err != nil {
			return nil, err
		}
		if err := rows.Err(); err != nil {
			logging.Db.Error().Err(err).Msg("Error occurred while iterating over monitor rows")
			return nil, err
		}

		return allMonitors, nil
	})
}

func (r *monitorRepository) GetMonitorByID(ctx context.Context, id string) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "GetMonitorByID", func() (monitors.IConcreteMonitor, error) {
		row, err := r.pool.Query(ctx,
			`SELECT * FROM monitors WHERE id=$1`,
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
			`SELECT * FROM monitors`)

		if err != nil {
			return nil, err
		}

		return pgx.CollectRows(rows, monitorFromCollectableRow)
	})
}

func (r *monitorRepository) DeleteMonitorByDisplayID(ctx context.Context, displayID string) (*pgtype.UUID, error) {
	return dbWrap(ctx, "DeleteMonitor", func() (*pgtype.UUID, error) {
		result := r.pool.QueryRow(ctx, `DELETE FROM monitors WHERE display_id=$1 RETURNING id`, displayID)

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

// InsertMonitor adds a new monitor to the database and returns its DisplayID (short DisplayID).
func (r *monitorRepository) InsertMonitor(ctx context.Context, monitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "InsertMonitor", func() (monitors.IConcreteMonitor, error) {
		rows, err := r.pool.Query(ctx,
			`INSERT INTO monitors (display_id, team_id, group_id, name, description, interval, kind, config)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
			monitor.GetDisplayID(),
			monitor.GetTeamID(),
			monitor.GetGroupID(),
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

		monitor, err := pgx.CollectOneRow(rows, monitorFromCollectableRow)
		if err != nil {
			return nil, err
		}

		return monitor, nil
	})
}

func (r *monitorRepository) UpdateMonitor(ctx context.Context, newMonitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error) {
	return dbWrap(ctx, "UpdateMonitor", func() (monitors.IConcreteMonitor, error) {
		result, err := r.pool.Query(ctx,
			`UPDATE monitors SET display_id=$1, team_id=$2, group_id=$3, name=$4, description=$5, interval=$6, kind=$7, config=$8 WHERE id=$9 RETURNING *`,
			newMonitor.GetDisplayID(),
			newMonitor.GetTeamID(),
			newMonitor.GetGroupID(),
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
