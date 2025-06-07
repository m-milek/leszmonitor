package common

import (
	"time"
)

type User struct {
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"password_hash" bson:"password_hash"`
	Email        string `json:"email" bson:"email"`
	CreatedAt    string `json:"created" bson:"created"`
	UpdatedAt    string `json:"updated" bson:"updated"`
}

// NewUser creates a new User instance with the provided username, password, and email.
// It hashes the password using the api_util.HashPassword function and sets the CreatedAt and UpdatedAt fields to the current time.
//
// Fails if the password hashing fails due to length over 72 bytes.
func NewUser(username, hashedPassword, email string) *User {
	return &User{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        email,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}
}
