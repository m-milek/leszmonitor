package db

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
)

type IGroupRepository interface {
	InsertGroup(ctx context.Context, group *models.MonitorGroup) error
	GetGroupByDisplayID(ctx context.Context, displayID string) (*models.MonitorGroup, error)
	GetGroupsByTeamID(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error)
	UpdateGroup(ctx context.Context, team *models.Team, oldGroup, newGroup *models.MonitorGroup) (bool, error)
	DeleteGroup(ctx context.Context, team *models.Team, groupID string) (bool, error)
}

type groupRepository struct {
	baseRepository
}

// groupFromCollectableRow maps a pgx.CollectableRow to a models.MonitorGroup struct.
func groupFromCollectableRow(row pgx.CollectableRow) (models.MonitorGroup, error) {
	group := models.MonitorGroup{}
	err := row.Scan(&group.ID, &group.TeamID, &group.DisplayID, &group.Name, &group.Description, &group.CreatedAt, &group.UpdatedAt)

	return group, err
}

func newGroupRepository(repository baseRepository) IGroupRepository {
	return &groupRepository{
		baseRepository: repository,
	}
}

func (r *groupRepository) InsertGroup(ctx context.Context, group *models.MonitorGroup) error {
	_, err := dbWrap(ctx, "InsertGroup", func() (*any, error) {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO groups (team_id, display_id, name, description) VALUES ($1, $2, $3, $4)`,
			group.TeamID, group.DisplayID, group.Name, group.Description)
		return nil, err
	})
	if pgErrIs(err, pgerrcode.UniqueViolation) {
		return ErrAlreadyExists
	}
	return err
}

func (r *groupRepository) GetGroupByDisplayID(ctx context.Context, displayID string) (*models.MonitorGroup, error) {
	return dbWrap(ctx, "GetGroupByDisplayID", func() (*models.MonitorGroup, error) {
		row, err := r.pool.Query(ctx,
			`SELECT id, team_id, display_id, name, description, created_at, updated_at FROM groups WHERE display_id=$1`,
			displayID)

		if err != nil {
			return nil, err
		}

		group, err := pgx.CollectExactlyOneRow(row, groupFromCollectableRow)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}

		return &group, err
	})
}

func (r *groupRepository) GetGroupsByTeamID(ctx context.Context, team *models.Team) ([]models.MonitorGroup, error) {
	return dbWrap(ctx, "GetGroupsByTeamID", func() ([]models.MonitorGroup, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, team_id, display_id, name, description, created_at, updated_at FROM groups WHERE team_id=$1`,
			team.ID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				logging.Db.Info().Msgf("No monitor groups found for team %s", team.DisplayID)
				return []models.MonitorGroup{}, nil
			}
			return nil, err
		}

		var groups []models.MonitorGroup
		groups, err = pgx.CollectRows(rows, groupFromCollectableRow)
		if err != nil {
			return nil, err
		}
		return groups, nil
	})
}

func (r *groupRepository) UpdateGroup(ctx context.Context, team *models.Team, oldGroup, newGroup *models.MonitorGroup) (bool, error) {
	return dbWrap(ctx, "UpdateGroup", func() (bool, error) {
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

func (r *groupRepository) DeleteGroup(ctx context.Context, team *models.Team, groupID string) (bool, error) {
	return dbWrap(ctx, "DeleteGroup", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`DELETE FROM groups WHERE display_id=$1 AND team_id=$2`,
			groupID, team.ID)
		if err != nil {
			return false, err
		}
		return result.RowsAffected() > 0, nil
	})
}
