package db

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models"
)

type IOrgRepository interface {
	InsertOrg(ctx context.Context, org *models.Org) (*struct{}, error)
	GetOrgByDisplayID(ctx context.Context, displayID string) (*models.Org, error)
	GetAllOrgs(ctx context.Context) ([]models.Org, error)
	UpdateOrg(ctx context.Context, org *models.Org) (bool, error)
	DeleteOrgByID(ctx context.Context, displayID string) (bool, error)
	AddMemberToOrg(ctx context.Context, orgDisplayID string, member *models.OrgMember) (bool, error)
	RemoveMemberFromOrg(ctx context.Context, orgDisplayID string, userID pgtype.UUID) (bool, error)
}

type orgRepository struct {
	baseRepository
}

func newOrgRepository(repository baseRepository) IOrgRepository {
	return &orgRepository{
		baseRepository: repository,
	}
}

// orgMemberFromCollectableRow maps a pgx.CollectableRow to a models.OrgMember struct.
func orgMemberFromCollectableRow(row pgx.CollectableRow) (models.OrgMember, error) {
	member := models.OrgMember{}
	err := row.Scan(&member.ID, &member.Username, &member.Role, &member.CreatedAt, &member.UpdatedAt)

	return member, err
}

// orgFromCollectableRow maps a pgx.CollectableRow to a models.Org struct.
func orgFromCollectableRow(row pgx.CollectableRow) (models.Org, error) {
	var org models.Org
	err := row.Scan(&org.ID, &org.DisplayID, &org.Name, &org.Description, &org.CreatedAt, &org.UpdatedAt)

	return org, err
}

func (r *orgRepository) InsertOrg(ctx context.Context, org *models.Org) (*struct{}, error) {
	return dbWrap(ctx, "InsertOrg", func() (*struct{}, error) {
		tx, err := r.pool.Begin(ctx)
		if err != nil {
			return nil, err
		}
		defer func() {
			_ = tx.Rollback(ctx)
		}()

		// Create the org and get its ID created by the DB
		var orgID pgtype.UUID
		row := tx.QueryRow(ctx,
			`INSERT INTO orgs (display_id, name, description) VALUES ($1, $2, $3) RETURNING id`,
			org.DisplayID, org.Name, org.Description)
		err = row.Scan(&orgID)
		if err != nil {
			if pgErrIs(err, pgerrcode.UniqueViolation) {
				return nil, ErrAlreadyExists
			}
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		// Create the owner record
		_, err = tx.Exec(ctx,
			`INSERT INTO user_orgs (org_id, user_id, role) VALUES ($1, $2, $3)`,
			orgID, org.Members[0].ID, org.Members[0].Role)
		if err != nil {
			if pgErrIs(err, pgerrcode.UniqueViolation) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		if err = tx.Commit(ctx); err != nil {
			return nil, err
		}

		return nil, nil
	})
}

func (r *orgRepository) GetOrgByDisplayID(ctx context.Context, displayID string) (*models.Org, error) {
	return dbWrap(ctx, "GetOrgByDisplayID", func() (*models.Org, error) {
		var org models.Org
		row, err := r.pool.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM orgs WHERE display_id=$1`,
			displayID)
		if err != nil {
			return nil, err
		}

		org, collectErr := pgx.CollectExactlyOneRow(row, orgFromCollectableRow)
		if collectErr != nil {
			if errors.Is(collectErr, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, collectErr
		}

		memberRows, membersErr := r.pool.Query(ctx,
			`SELECT ut.user_id, u.username, ut.role, ut.created_at, ut.updated_at 
             FROM user_orgs ut JOIN users u ON u.id = ut.user_id 
             WHERE org_id=$1`,
			org.ID)
		if membersErr != nil {
			return nil, membersErr
		}

		members, err := pgx.CollectRows(memberRows, orgMemberFromCollectableRow)
		if err != nil {
			return nil, err
		}
		org.Members = members

		return &org, nil
	})
}

func (r *orgRepository) GetAllOrgs(ctx context.Context) ([]models.Org, error) {
	return dbWrap(ctx, "GetAllOrgs", func() ([]models.Org, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, display_id, name, description, created_at, updated_at FROM orgs`)
		if err != nil {
			return nil, err
		}

		orgs, err := pgx.CollectRows(rows, orgFromCollectableRow)
		if err != nil {
			return nil, err
		}

		for i, org := range orgs {
			memberRows, membersErr := r.pool.Query(ctx,
				`SELECT user_id, u.username, role, ut.created_at, ut.updated_at FROM user_orgs ut JOIN users u ON u.id = ut.user_id WHERE org_id=$1`,
				org.ID)
			if membersErr != nil {
				return nil, membersErr
			}

			members, err := pgx.CollectRows(memberRows, orgMemberFromCollectableRow)
			if err != nil {
				return nil, err
			}
			orgs[i].Members = members
		}

		return orgs, nil
	})
}

func (r *orgRepository) UpdateOrg(ctx context.Context, org *models.Org) (bool, error) {
	return dbWrap(ctx, "UpdateOrg", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`UPDATE orgs SET display_id=$1, name=$2, description=$3 WHERE id=$4`,
			org.DisplayID, org.Name, org.Description, org.ID)
		if err != nil {
			return false, err
		}
		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}

func (r *orgRepository) DeleteOrgByID(ctx context.Context, displayID string) (bool, error) {
	return dbWrap(ctx, "DeleteOrg", func() (bool, error) {
		result, err := r.pool.Exec(ctx,
			`DELETE FROM orgs WHERE display_id=$1`,
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

func (r *orgRepository) AddMemberToOrg(ctx context.Context, orgDisplayID string, member *models.OrgMember) (bool, error) {
	return dbWrap(ctx, "AddMemberToOrg", func() (bool, error) {
		var orgID pgtype.UUID
		row := r.pool.QueryRow(ctx,
			`SELECT id FROM orgs WHERE display_id=$1`,
			orgDisplayID)

		err := row.Scan(&orgID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.Exec(ctx,
			`INSERT INTO user_orgs (org_id, user_id, role) VALUES ($1, $2, $3)`,
			orgID, member.ID, member.Role)
		if err != nil {
			return false, err
		}

		if result.RowsAffected() == 0 {
			return false, ErrAlreadyExists
		}
		return true, nil
	})
}

func (r *orgRepository) RemoveMemberFromOrg(ctx context.Context, orgDisplayID string, userID pgtype.UUID) (bool, error) {
	return dbWrap(ctx, "RemoveMemberFromOrg", func() (bool, error) {
		var orgID pgtype.UUID
		row := r.pool.QueryRow(ctx,
			`SELECT id FROM orgs WHERE display_id=$1`,
			orgDisplayID)

		err := row.Scan(&orgID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return false, ErrNotFound
			}
			return false, err
		}

		result, err := r.pool.Exec(ctx,
			`DELETE FROM user_orgs WHERE org_id=$1 AND user_id=$2`,
			orgID, userID)
		if err != nil {
			return false, err
		}

		if result.RowsAffected() == 0 {
			return false, ErrNotFound
		}
		return true, nil
	})
}
