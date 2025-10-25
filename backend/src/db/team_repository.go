package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models"
)

type ITeamRepository interface {
	InsertTeam(ctx context.Context, team *models.Team) (*struct{}, error)
	GetTeamByID(ctx context.Context, displayID string) (*models.Team, error)
	GetAllTeams(ctx context.Context) ([]models.Team, error)
	UpdateTeam(ctx context.Context, team *models.Team) (bool, error)
	DeleteTeamByID(ctx context.Context, displayID string) (bool, error)
	AddMemberToTeam(ctx context.Context, teamDisplayID string, member *models.TeamMember) (bool, error)
	RemoveMemberFromTeam(ctx context.Context, teamDisplayID string, userID pgtype.UUID) (bool, error)
}

type teamRepository struct {
	baseRepository
}

func newTeamRepository(repository baseRepository) ITeamRepository {
	return &teamRepository{
		baseRepository: repository,
	}
}

// teamMemberFromCollectableRow maps a pgx.CollectableRow to a models.TeamMember struct.
func teamMemberFromCollectableRow(row pgx.CollectableRow) (models.TeamMember, error) {
	member := models.TeamMember{}
	err := row.Scan(&member.ID, &member.Role, &member.CreatedAt, &member.UpdatedAt)

	return member, err
}

// teamFromCollectableRow maps a pgx.CollectableRow to a models.Team struct.
func teamFromCollectableRow(row pgx.CollectableRow) (models.Team, error) {
	var team models.Team
	err := row.Scan(&team.ID, &team.DisplayID, &team.Name, &team.Description, &team.CreatedAt, &team.UpdatedAt)

	return team, err
}

func (r *teamRepository) InsertTeam(ctx context.Context, team *models.Team) (*struct{}, error) {
	return dbWrap(ctx, "InsertTeam", func() (*struct{}, error) {
		tx, err := r.pool.Begin(ctx)

		if err != nil {
			return nil, err
		}

		var teamID pgtype.UUID
		row := tx.QueryRow(ctx,
			`INSERT INTO teams (display_id, name, description) VALUES ($1, $2, $3) RETURNING id`,
			team.DisplayID, team.Name, team.Description)
		if err != nil {
			return nil, err
		}
		err = row.Scan(&teamID)
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO user_teams (team_id, user_id, role) VALUES ($1, $2, $3)`,
			teamID, team.Members[0].ID, team.Members[0].Role)
		if err != nil {
			return nil, err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

func (r *teamRepository) GetTeamByID(ctx context.Context, displayID string) (*models.Team, error) {
	return dbWrap(ctx, "GetTeamByID", func() (*models.Team, error) {
		var team models.Team
		row, err := r.pool.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM teams WHERE display_id=$1`,
			displayID)
		if err != nil {
			return nil, err
		}

		team, collectErr := pgx.CollectExactlyOneRow(row, teamFromCollectableRow)
		if collectErr != nil {
			if errors.Is(collectErr, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, collectErr
		}

		memberRows, membersErr := r.pool.Query(ctx,
			`SELECT user_id, role, created_at, updated_at FROM user_teams WHERE team_id=$1`,
			team.ID)
		if membersErr != nil {
			return nil, membersErr
		}

		members, err := pgx.CollectRows(memberRows, teamMemberFromCollectableRow)
		if err != nil {
			return nil, err
		}
		team.Members = members

		return &team, nil
	})
}

func (r *teamRepository) GetAllTeams(ctx context.Context) ([]models.Team, error) {
	return dbWrap(ctx, "GetAllTeam", func() ([]models.Team, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM teams`)
		if err != nil {
			return nil, err
		}

		teams, err := pgx.CollectRows(rows, teamFromCollectableRow)
		if err != nil {
			return nil, err
		}

		for i, team := range teams {
			memberRows, membersErr := r.pool.Query(ctx,
				`SELECT user_id, role, created_at, updated_at FROM user_teams WHERE team_id=$1`,
				team.ID)
			if membersErr != nil {
				return nil, membersErr
			}

			members, err := pgx.CollectRows(memberRows, teamMemberFromCollectableRow)
			if err != nil {
				return nil, err
			}
			teams[i].Members = members
		}

		return teams, nil
	})
}

func (r *teamRepository) UpdateTeam(ctx context.Context, team *models.Team) (bool, error) {
	return dbWrap(ctx, "UpdateTeam", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`UPDATE teams SET display_id=$1, name=$2, description=$3 WHERE id=$4`,
			team.DisplayID, team.Name, team.Description, team.ID)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}

func (r *teamRepository) DeleteTeamByID(ctx context.Context, displayID string) (bool, error) {
	return dbWrap(ctx, "DeleteTeam", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`DELETE FROM teams WHERE display_id=$1`,
			displayID)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}

func (r *teamRepository) AddMemberToTeam(ctx context.Context, teamDisplayID string, member *models.TeamMember) (bool, error) {
	return dbWrap(ctx, "AddMemberToTeam", func() (bool, error) {
		var teamID pgtype.UUID
		row := r.pool.QueryRow(ctx,
			`SELECT id FROM teams WHERE display_id=$1`,
			teamDisplayID)

		err := row.Scan(&teamID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.Exec(ctx,
			`INSERT INTO user_teams (team_id, user_id, role) VALUES ($1, $2, $3)`,
			teamID, member.ID, member.Role)
		if err != nil {
			return false, err
		}

		if result.RowsAffected() == 0 {
			return false, ErrAlreadyExists
		}
		return true, nil
	})
}

func (r *teamRepository) RemoveMemberFromTeam(ctx context.Context, teamDisplayID string, userID pgtype.UUID) (bool, error) {
	return dbWrap(ctx, "RemoveMemberFromTeam", func() (bool, error) {
		var teamID pgtype.UUID
		row := r.pool.QueryRow(ctx,
			`SELECT id FROM teams WHERE display_id=$1`,
			teamDisplayID)

		err := row.Scan(&teamID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.Exec(ctx,
			`DELETE FROM user_teams WHERE team_id=$1 AND user_id=$2`,
			teamID, userID)
		if err != nil {
			return false, err
		}

		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}
