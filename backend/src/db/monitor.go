package db

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/logging"
	monitors "github.com/m-milek/leszmonitor/uptime/monitor"
)

// monitorFromCollectableRow maps a pgx.CollectableRow to a monitors.IConcreteMonitor.
func monitorFromCollectableRow(row pgx.CollectableRow) (monitors.IConcreteMonitor, error) {
	var config []byte
	var b monitors.BaseMonitor

	err := row.Scan(&b.Id, &b.DisplayId, &b.TeamId, &b.GroupId, &b.Name, &b.Description, &b.Interval, &b.Type, &config, &b.CreatedAt, &b.UpdatedAt)
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

// CreateMonitor adds a new monitor to the database and returns its DisplayID (short DisplayID).
func CreateMonitor(ctx context.Context, monitor monitors.IConcreteMonitor) (monitors.IConcreteMonitor, error) {
	dbRes, err := withTimeout(ctx, func() (monitors.IConcreteMonitor, error) {
		rows, err := dbClient.conn.Query(ctx,
			`INSERT INTO monitors (display_id, team_id, group_id, name, description, interval, kind, config)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING *`,
			monitor.GetDisplayId(),
			monitor.GetTeamId(),
			monitor.GetGroupId(),
			monitor.GetName(),
			monitor.GetDescription(),
			int(monitor.GetInterval().Seconds()),
			string(monitor.GetType()),
			monitor.GetConfig(),
		)
		if err != nil {
			return nil, err
		}

		monitor, err := pgx.CollectOneRow(rows, monitorFromCollectableRow)
		if err != nil {
			return nil, err
		}

		return monitor, nil
	})

	logDbOperation("InsertMonitor", dbRes, err)

	if err != nil {
		return nil, err
	}

	// Broadcast that a monitor has been added
	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		Id:      dbRes.Result.GetId(),
		Status:  monitors.Created,
		Monitor: &dbRes.Result,
	})

	return dbRes.Result, nil
}

func GetAllMonitors(ctx context.Context) ([]monitors.IConcreteMonitor, error) {
	dbRes, err := withTimeout(ctx, func() ([]monitors.IConcreteMonitor, error) {
		rows, err := dbClient.conn.Query(ctx,
			`SELECT id, display_id, team_id, group_id, name, description, interval, kind, config FROM monitors`)

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

	logDbOperation("GetAllMonitors", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func DeleteMonitor(ctx context.Context, displayId string) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (*pgtype.UUID, error) {
		result := dbClient.conn.QueryRow(ctx, `DELETE FROM monitors WHERE display_id=$1 RETURNING id`, displayId)

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

	logDbOperation("DeleteMonitor", dbRes, err)

	if err != nil {
		return false, err
	}

	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		Id:      *dbRes.Result,
		Status:  monitors.Deleted,
		Monitor: nil,
	})

	if dbRes.Result == nil {
		return false, nil
	}
	return true, nil
}

func GetMonitorById(ctx context.Context, id string) (monitors.IMonitor, error) {
	dbRes, err := withTimeout(ctx, func() (monitors.IMonitor, error) {
		row, err := dbClient.conn.Query(ctx,
			`SELECT id, display_id, team_id, group_id, name, description, interval, kind, config FROM monitors WHERE id=$1`,
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

	logDbOperation("GetMonitorById", dbRes, err)

	return dbRes.Result, err
}

func UpdateMonitor(ctx context.Context, newMonitor monitors.IConcreteMonitor) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		result, err := dbClient.conn.Exec(ctx,
			`UPDATE monitors SET display_id=$1, team_id=$2, group_id=$3, name=$4, description=$5, interval=$6, kind=$7, config=$8 WHERE id=$9`,
			newMonitor.GetDisplayId(),
			newMonitor.GetTeamId(),
			newMonitor.GetGroupId(),
			newMonitor.GetName(),
			newMonitor.GetDescription(),
			int(newMonitor.GetInterval().Seconds()),
			string(newMonitor.GetType()),
			newMonitor.GetConfig(),
			newMonitor.GetId(),
		)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return result.RowsAffected() > 0, nil
	})

	logDbOperation("UpdateMonitor", dbRes, err)

	if err != nil {
		return false, err
	}

	// Broadcast that a monitor has been updated
	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		Id:      newMonitor.GetId(),
		Status:  monitors.Edited,
		Monitor: &newMonitor,
	})

	return dbRes.Result, nil
}
