package models

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/models/util"
)

// User represents a user in the system.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	PasswordHash string    `json:"-" db:"password_hash"`
	util.Timestamps
}

// NewUser creates a new User instance with the provided username, password, and email.
func NewUser(username, hashedPassword string) (*User, error) {
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
