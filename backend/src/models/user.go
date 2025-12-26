package models

import (
	"fmt"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/m-milek/leszmonitor/models/util"
)

// User represents a user in the system.
type User struct {
	ID           pgtype.UUID `json:"id"`
	Username     string      `json:"username"`
	PasswordHash string      `json:"-"`
	util.Timestamps
}

// NewUser creates a new User instance with the provided username, password, and email.
func NewUser(username, hashedPassword, email string) (*User, error) {
	user := &User{
		Username:     username,
		PasswordHash: hashedPassword,
	}
	err := user.Validate()

	if err != nil {
		return nil, err
	}
	return user, nil
}

// Validate checks if the User has valid fields.
func (u *User) Validate() error {
	if u.Username == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if u.PasswordHash == "" {
		return fmt.Errorf("password hash cannot be empty")
	}
	return nil
}
