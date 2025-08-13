package models

import (
	"time"
)

type RawUser struct {
	Username     string `json:"username" bson:"username"`
	PasswordHash string `json:"passwordHash" bson:"passwordHash"`
	Email        string `json:"email" bson:"email"`
	CreatedAt    string `json:"created" bson:"created"`
	UpdatedAt    string `json:"updated" bson:"updated"`
}

// NewUser creates a new RawUser instance with the provided username, password, and email.
func NewUser(username, hashedPassword, email string) *RawUser {
	return &RawUser{
		Username:     username,
		PasswordHash: hashedPassword,
		Email:        email,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}
}

func (u *RawUser) IntoUser() *UserResponse {
	return &UserResponse{
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
	}
}

type UserResponse struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created"`
}
