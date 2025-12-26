package db

import (
	"context"
	"errors"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models"
)

type IUserRepository interface {
	InsertUser(ctx context.Context, user *models.User) (*struct{}, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id pgtype.UUID) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type userRepository struct {
	baseRepository
}

func newUserRepository(repository baseRepository) IUserRepository {
	return &userRepository{
		baseRepository: repository,
	}
}

func userFromCollectableRow(row pgx.CollectableRow) (models.User, error) {
	var user models.User
	err := row.Scan(&user.ID, &user.Username, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)

	return user, err
}

func (r *userRepository) InsertUser(ctx context.Context, user *models.User) (*struct{}, error) {
	return dbWrap(ctx, "CreateUser", func() (*struct{}, error) {
		_, err := r.pool.Exec(ctx,
			`INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING *`,
			user.Username, user.PasswordHash)

		if err != nil {
			if pgErrIs(err, pgerrcode.UniqueViolation) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}
		return nil, err
	})
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return dbWrap(ctx, "GetUserByUsername", func() (*models.User, error) {
		row, err := r.pool.Query(ctx,
			`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username=$1`,
			username)
		if err != nil {
			return nil, err
		}

		user, err := pgx.CollectExactlyOneRow(row, userFromCollectableRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})
}

func (r *userRepository) GetUserByID(ctx context.Context, id pgtype.UUID) (*models.User, error) {
	return dbWrap(ctx, "GetUserByUsername", func() (*models.User, error) {
		row, err := r.pool.Query(ctx,
			`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id=$1`,
			id)
		if err != nil {
			return nil, err
		}

		user, err := pgx.CollectExactlyOneRow(row, userFromCollectableRow)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return dbWrap(ctx, "GetAllUsers", func() ([]models.User, error) {
		rows, err := r.pool.Query(ctx,
			`SELECT id, username, password_hash, created_at, updated_at FROM users`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		users, err := pgx.CollectRows(rows, userFromCollectableRow)

		return users, err
	})
}
