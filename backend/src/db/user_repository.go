package db

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models"
)

type IUserRepository interface {
	InsertUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
}

type UserRepository struct {
	baseRepository
}

func newUserRepository(repository baseRepository) IUserRepository {
	return &UserRepository{
		baseRepository: repository,
	}
}

func (r *UserRepository) InsertUser(ctx context.Context, user *models.User) (*models.User, error) {
	return dbWrap(ctx, "CreateUser", func() (*models.User, error) {
		if user.ID == uuid.Nil {
			user.ID = uuid.New()
		}

		var createdUser models.User
		err := r.pool.QueryRowxContext(ctx,
			`INSERT INTO users (id, username, password_hash) VALUES ($1, $2, $3) RETURNING id, username, password_hash, created_at, updated_at`,
			user.ID, user.Username, user.PasswordHash,
		).StructScan(&createdUser)

		if err != nil {
			if isUniqueViolation(err) {
				return nil, ErrAlreadyExists
			}
			return nil, err
		}

		return &createdUser, nil
	})
}

func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	return dbWrap(ctx, "GetUserByUsername", func() (*models.User, error) {
		var user models.User
		err := r.pool.GetContext(ctx, &user,
			`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE username=$1`,
			username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})
}

func (r *UserRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return dbWrap(ctx, "GetUserByID", func() (*models.User, error) {
		var user models.User
		err := r.pool.GetContext(ctx, &user,
			`SELECT id, username, password_hash, created_at, updated_at FROM users WHERE id=$1`,
			id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})
}

func (r *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	return dbWrap(ctx, "GetAllUsers", func() ([]models.User, error) {
		var users []models.User
		err := r.pool.SelectContext(ctx, &users,
			`SELECT id, username, password_hash, created_at, updated_at FROM users`)
		if err != nil {
			return nil, err
		}
		if users == nil {
			users = []models.User{}
		}

		return users, nil
	})
}
