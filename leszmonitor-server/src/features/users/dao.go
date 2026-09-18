package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/db"
)

type IUserDAO interface {
	InsertUser(ctx context.Context, user *User) (*User, error)
	UpdateUserRole(ctx context.Context, userID uuid.UUID, role auth.Role) (*User, error)
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	GetAllUsers(ctx context.Context) ([]User, error)
}

type UserDAO struct {
	pool db.Querier
}

func NewUserDAO(pool db.Querier) IUserDAO {
	return &UserDAO{
		pool: pool,
	}
}

func (r *UserDAO) InsertUser(ctx context.Context, user *User) (*User, error) {
	return db.Wrap(ctx, "CreateUser", func() (*User, error) {
		if user.ID == uuid.Nil {
			user.ID = uuid.New()
		}

		var createdUser User
		err := r.pool.QueryRowxContext(
			ctx,
			`INSERT INTO users (id, username, role, password_hash) VALUES ($1, $2, $3, $4) RETURNING id, username, role, password_hash, created_at, updated_at`,
			user.ID,
			user.Username,
			user.Role,
			user.PasswordHash,
		).StructScan(&createdUser)

		if err != nil {
			if db.IsUniqueViolation(err) {
				return nil, db.ErrAlreadyExists
			}
			return nil, err
		}

		return &createdUser, nil
	})
}

func (r *UserDAO) UpdateUserRole(
	ctx context.Context,
	userID uuid.UUID,
	role auth.Role,
) (*User, error) {
	return db.Wrap(ctx, "UpdateUserRole", func() (*User, error) {
		var updatedUser User
		err := r.pool.QueryRowxContext(
			ctx,
			`UPDATE users SET role = $1 WHERE id = $2 RETURNING id, username, role, password_hash, created_at, updated_at`,
			role,
			userID,
		).StructScan(&updatedUser)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}
		return &updatedUser, nil
	})
}

func (r *UserDAO) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return db.Wrap(ctx, "GetUserByUsername", func() (*User, error) {
		var user User
		err := sqlx.GetContext(ctx, r.pool, &user,
			`SELECT id, username, role, password_hash, created_at, updated_at FROM users WHERE username=$1`,
			username)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})
}

func (r *UserDAO) GetUserByID(ctx context.Context, id uuid.UUID) (*User, error) {
	return db.Wrap(ctx, "GetUserByID", func() (*User, error) {
		var user User
		err := sqlx.GetContext(ctx, r.pool, &user,
			`SELECT id, username, role, password_hash, created_at, updated_at FROM users WHERE id=$1`,
			id)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}
		return &user, nil
	})
}

func (r *UserDAO) GetAllUsers(ctx context.Context) ([]User, error) {
	return db.Wrap(ctx, "GetAllUsers", func() ([]User, error) {
		var users []User
		err := sqlx.SelectContext(ctx, r.pool, &users,
			`SELECT id, username, role, password_hash, created_at, updated_at FROM users`)
		if err != nil {
			return nil, err
		}
		if users == nil {
			users = []User{}
		}

		return users, nil
	})
}
