package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
)

type IGroupRepository interface {
	CreateMonitorGroup(ctx context.Context, group *models.MonitorGroup) (*models.MonitorGroup, error)
	GetMonitorGroupByID(ctx context.Context, displayID string) (*models.MonitorGroup, error)
	GetMonitorGroupsForTeam(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error)
	UpdateMonitorGroup(ctx context.Context, team *models.Team, oldGroup, newGroup *models.MonitorGroup) (bool, error)
	DeleteMonitorGroup(ctx context.Context, team *models.Team, groupID string) (bool, error)
}

type groupRepository struct {
	baseRepository
}

func newGroupRepository(repository baseRepository) IGroupRepository {
	return &groupRepository{
		baseRepository: repository,
	}
}

func (r *groupRepository) CreateMonitorGroup(ctx context.Context, group *models.MonitorGroup) (*models.MonitorGroup, error) {
	return dbWrap(ctx, "CreateMonitorGroup", func() (*models.MonitorGroup, error) {
		rows := r.pool.QueryRow(ctx,
			`INSERT INTO groups (display_id, team_id, name, description) VALUES ($1, $2, $3, $4) RETURNING id`,
			group.DisplayID, group.TeamID, group.Name, group.Description)

		err := rows.Scan(&group.ID)

		return group, err
	})
}

func (r *groupRepository) GetMonitorGroupByID(ctx context.Context, displayID string) (*models.MonitorGroup, error) {
	return dbWrap(ctx, "GetMonitorGroupByID", func() (*models.MonitorGroup, error) {
		rows := r.pool.QueryRow(ctx,
			`SELECT id, display_id, team_id, name, description, created_at, updated_at FROM groups WHERE display_id=$1`,
			displayID)

		group := &models.MonitorGroup{}
		err := rows.Scan(&group.ID, &group.DisplayID, &group.TeamID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)

		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		return group, err
	})
}

func (r *groupRepository) GetMonitorGroupsForTeam(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error) {
	return dbWrap(ctx, "GetMonitorGroupsByTeamID", func() ([]models.MonitorGroup, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, display_id, team_id, name, description, created_at, updated_at FROM groups WHERE team_id=$1`,
			team.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logging.Db.Info().Msgf("No monitor groups found for team %s", team.DisplayID)
				return []models.MonitorGroup{}, nil
			}
			return nil, err
		}

		var groups []models.MonitorGroup
		groups, err = pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.MonitorGroup, error) {
			var group models.MonitorGroup
			err := row.Scan(&group.ID, &group.DisplayID, &group.TeamID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)
			return group, err
		})
		if err != nil {
			return nil, err
		}
		return groups, nil
	})
}

func (r *groupRepository) UpdateMonitorGroup(ctx context.Context, team *models.Team, oldGroup, newGroup *models.MonitorGroup) (bool, error) {
	return dbWrap(ctx, "UpdateMonitorGroup", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`UPDATE groups SET display_id=$1, name=$2, description=$3 WHERE id=$4 AND team_id=$5 RETURNING *`,
			newGroup.DisplayID, newGroup.Name, newGroup.Description, oldGroup.ID, team.ID)

		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}

		return true, nil
	})
}

func (r *groupRepository) DeleteMonitorGroup(ctx context.Context, team *models.Team, groupID string) (bool, error) {
	return dbWrap(ctx, "DeleteMonitorGroup", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`DELETE FROM groups WHERE display_id=$1 AND team_id=$2`,
			groupID, team.ID)
		if err != nil {
			return false, err
		}
		return result.RowsAffected() > 0, nil
	})
}
