package db

import (
	"context"
	"errors"
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

func (r *userRepository) InsertUser(ctx context.Context, user *models.User) (*struct{}, error) {
	return dbWrap(ctx, "CreateUser", func() (*struct{}, error) {
		status, err := r.pool.Exec(ctx,
			`INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3) RETURNING *`,
			user.Username, user.Email, user.PasswordHash)

		if err != nil {
			return nil, err
		}
		if status.RowsAffected() == 0 {
			return nil, ErrAlreadyExists
		}
		return nil, err
	})
}

func (r *userRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return dbWrap(ctx, "GetUserByUsername", func() (*models.User, error) {
		var user models.User
		row := r.pool.QueryRow(ctx,
			`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username=$1`,
			username)
		err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
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
		var user models.User
		row := r.pool.QueryRow(ctx,
			`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE id=$1`,
			id)
		err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt)
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
			`SELECT id, username, email, password_hash, created_at, updated_at FROM users`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		usersList := make([]models.User, 0)

		for rows.Next() {
			var user models.User
			if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
				return nil, err
			}
			usersList = append(usersList, user)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return usersList, nil
	})
}
