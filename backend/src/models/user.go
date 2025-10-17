package models

import "github.com/jackc/pgx/v5/pgtype"

type User struct {
	Id           pgtype.UUID `json:"id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"-"`
	Email        string      `json:"email"`
	Timestamps
}

// NewUser creates a new User instance with the provided username, password, and email.
func NewUser(username, hashedPassword, email string) *User {
	return &User{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        email,
	}
}
