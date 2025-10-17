package db

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/m-milek/leszmonitor/models"
)

func CreateUser(ctx context.Context, user *models.User) error {
	dbRes, err := withTimeout(ctx, func() (*struct{}, error) {
		status, err := dbClient.conn.Exec(ctx,
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

	logDbOperation("CreateUser", dbRes, err)

	if err != nil {
		return err
	}
	return nil
}

func GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	dbRes, err := withTimeout(ctx, func() (*models.User, error) {
		var user models.User
		var timestamps models.RawTimestamps
		row := dbClient.conn.QueryRow(ctx,
			`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username=$1`,
			username)
		err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &timestamps.CreatedAt, &timestamps.UpdatedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		user.SetTimestamps(timestamps)
		return &user, nil
	})

	logDbOperation("GetUserByUsername", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetRawUserByUsername(ctx context.Context, username string) (*models.User, error) {
	dbRes, err := withTimeout(ctx, func() (*models.User, error) {
		var user models.User
		var timestamps models.RawTimestamps
		row := dbClient.conn.QueryRow(ctx,
			`SELECT id, username, email, password_hash, created_at, updated_at FROM users WHERE username=$1`,
			username)
		err := row.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &timestamps.CreatedAt, &timestamps.UpdatedAt)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}
		user.SetTimestamps(timestamps)
		return &user, nil
	})

	logDbOperation("GetRawUserByUsername", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}

func GetAllUsers(ctx context.Context) ([]models.User, error) {
	dbRes, err := withTimeout(ctx, func() ([]models.User, error) {
		rows, err := dbClient.conn.Query(ctx,
			`SELECT id, username, email, password_hash, created_at, updated_at FROM users`)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		usersList := make([]models.User, 0)

		for rows.Next() {
			var user models.User
			var timestamps models.RawTimestamps
			if err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &timestamps.CreatedAt, &timestamps.UpdatedAt); err != nil {
				return nil, err
			}
			user.SetTimestamps(timestamps)
			usersList = append(usersList, user)
		}

		if err := rows.Err(); err != nil {
			return nil, err
		}

		return usersList, nil
	})

	logDbOperation("GetAllUsers", dbRes, err)

	if err != nil {
		return nil, err
	}
	return dbRes.Result, nil
}
