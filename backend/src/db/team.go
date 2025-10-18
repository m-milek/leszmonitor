package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models"
)

// teamMemberFromCollectableRow maps a pgx.CollectableRow to a models.TeamMember struct.
func teamMemberFromCollectableRow(row pgx.CollectableRow) (models.TeamMember, error) {
	member := models.TeamMember{}
	err := row.Scan(&member.ID, &member.Role, &member.CreatedAt, &member.UpdatedAt)

	return member, err
}

// teamFromCollectableRow maps a pgx.CollectableRow to a models.Team struct.
func teamFromCollectableRow(row pgx.CollectableRow) (models.Team, error) {
	var team models.Team
	//var timestamps models.Timestamps
	err := row.Scan(&team.ID, &team.DisplayID, &team.Name, &team.Description, &team.CreatedAt, &team.UpdatedAt)

	return team, err
}

func CreateTeam(ctx context.Context, team *models.Team) (*struct{}, error) {
	dbRes, err := withTimeout(ctx, func() (*struct{}, error) {
		tx, err := dbClient.conn.Begin(ctx)

		if err != nil {
			return nil, err
		}

		var teamId pgtype.UUID
		row := tx.QueryRow(ctx,
			`INSERT INTO teams (display_id, name, description) VALUES ($1, $2, $3) RETURNING id`,
			team.DisplayID, team.Name, team.Description)
		if err != nil {
			return nil, err
		}
		err = row.Scan(&teamId)
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(ctx,
			`INSERT INTO user_teams (team_id, user_id, role) VALUES ($1, $2, $3)`,
			teamId, team.Members[0].ID, team.Members[0].Role)
		if err != nil {
			return nil, err
		}

		err = tx.Commit(ctx)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})

	logDbOperation("CreateTeam", dbRes, err)

	return dbRes.Result, err
}

func GetTeamById(ctx context.Context, displayId string) (*models.Team, error) {
	dbRes, err := withTimeout(ctx, func() (*models.Team, error) {
		var team models.Team
		row, err := dbClient.conn.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM teams WHERE display_id=$1`,
			displayId)
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

		memberRows, membersErr := dbClient.conn.Query(ctx,
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

	logDbOperation("GetTeamById", dbRes, err)

	return dbRes.Result, err
}

func GetAllTeams(ctx context.Context) ([]models.Team, error) {
	dbRes, err := withTimeout(ctx, func() ([]models.Team, error) {
		rows, err := dbClient.conn.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM teams`)
		if err != nil {
			return nil, err
		}

		teams, err := pgx.CollectRows(rows, teamFromCollectableRow)
		if err != nil {
			return nil, err
		}

		for i, team := range teams {
			memberRows, membersErr := dbClient.conn.Query(ctx,
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

	logDbOperation("GetAllTeams", dbRes, err)

	return dbRes.Result, err
}

func UpdateTeam(ctx context.Context, team *models.Team) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		result, err := dbClient.conn.Exec(ctx,
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

	logDbOperation("UpdateTeam", dbRes, err)

	return dbRes.Result, err
}

func DeleteTeam(ctx context.Context, displayId string) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		result, err := dbClient.conn.Exec(ctx,
			`DELETE FROM teams WHERE display_id=$1`,
			displayId)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})

	logDbOperation("DeleteTeam", dbRes, err)

	return dbRes.Result, err
}

func AddMemberToTeam(ctx context.Context, teamDisplayId string, member *models.TeamMember) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		var teamId pgtype.UUID
		row := dbClient.conn.QueryRow(ctx,
			`SELECT id FROM teams WHERE display_id=$1`,
			teamDisplayId)

		err := row.Scan(&teamId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := dbClient.conn.Exec(ctx,
			`INSERT INTO user_teams (team_id, user_id, role) VALUES ($1, $2, $3)`,
			teamId, member.ID, member.Role)
		if err != nil {
			return false, err
		}

		if result.RowsAffected() == 0 {
			return false, ErrAlreadyExists
		}
		return true, nil
	})

	logDbOperation("AddMemberToTeam", dbRes, err)

	return dbRes.Result, err
}

func RemoveMemberFromTeam(ctx context.Context, teamDisplayId string, userId pgtype.UUID) (bool, error) {
	dbRes, err := withTimeout(ctx, func() (bool, error) {
		var teamId pgtype.UUID
		row := dbClient.conn.QueryRow(ctx,
			`SELECT id FROM teams WHERE display_id=$1`,
			teamDisplayId)

		err := row.Scan(&teamId)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := dbClient.conn.Exec(ctx,
			`DELETE FROM user_teams WHERE team_id=$1 AND user_id=$2`,
			teamId, userId)
		if err != nil {
			return false, err
		}

		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})

	logDbOperation("RemoveMemberFromTeam", dbRes, err)

	return dbRes.Result, err
}
