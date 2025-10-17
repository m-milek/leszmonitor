package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
)

func CreateMonitorGroup(ctx context.Context, group *models.MonitorGroup) (*models.MonitorGroup, error) {
	dbRes, err := withTimeout(ctx, func() (*models.MonitorGroup, error) {
		rows := dbClient.conn.QueryRow(ctx,
			`INSERT INTO groups (display_id, team_id, name, description) VALUES ($1, $2, $3, $4) RETURNING id`,
			group.DisplayId, group.TeamId, group.Name, group.Description)

		err := rows.Scan(&group.Id)

		return group, err
	})

	logDbOperation("CreateMonitorGroup", dbRes, err)

	return dbRes.Result, err
}

func GetMonitorGroupById(ctx context.Context, displayId string) (*models.MonitorGroup, error) {
	dbRes, err := withTimeout(ctx, func() (*models.MonitorGroup, error) {
		rows := dbClient.conn.QueryRow(ctx,
			`SELECT id, display_id, team_id, name, description FROM groups WHERE display_id=$1`,
			displayId)

		group := &models.MonitorGroup{}
		err := rows.Scan(&group.Id, &group.DisplayId, &group.TeamId, &group.Name, &group.Description)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return group, err
	})

	logDbOperation("GetMonitorGroupById", dbRes, err)

	return dbRes.Result, err
}

func GetMonitorGroupsForTeam(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error) {
	dbRes, err := withTimeout(ctx, func() ([]models.MonitorGroup, error) {
		rows, err := dbClient.conn.Query(ctx,
			`SELECT id, display_id, team_id, name, description FROM groups WHERE team_id=$1`,
			team.Id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logging.Db.Info().Msgf("No monitor groups found for team %s", team.DisplayId)
				return []models.MonitorGroup{}, nil
			}
			return nil, err
		}

		var groups []models.MonitorGroup
		groups, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.MonitorGroup, error) {
			var group models.MonitorGroup
			err := row.Scan(&group.Id, &group.DisplayId, &group.TeamId, &group.Name, &group.Description)
			return group, err
		})
		if err != nil {
			return nil, err
		}
		return groups, nil
	})

	logDbOperation("GetTeamMonitorGroups", dbRes, err)

	return dbRes.Result, err
}

func UpdateMonitorGroup(ctx context.Context, team *models.Team, oldGroup, newGroup *models.MonitorGroup) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		result, err := dbClient.conn.Exec(ctx,
			`UPDATE groups SET name=$1, description=$2 WHERE id=$3 AND team_id=$4 RETURNING *`,
			newGroup.Name, newGroup.Description, oldGroup.Id, team.Id)

		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}

		return true, nil
	})

	logDbOperation("UpdateMonitorGroup", dbRes, err)

	return dbRes.Result, err
}

func DeleteMonitorGroup(ctx context.Context, team *models.Team, groupId string) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		result, err := dbClient.conn.Exec(ctx,
			`DELETE FROM groups WHERE display_id=$1 AND team_id=$2`,
			groupId, team.Id)
		if err != nil {
			return false, err
		}
		return result.RowsAffected() > 0, nil
	})

	logDbOperation("DeleteMonitorGroup", dbRes, err)

	return dbRes.Result, err
}
