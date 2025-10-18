package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser_Success(t *testing.T) {
	user, err := NewUser("testuser", "hashedpassword", "mail@example.com")

	assert.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "hashedpassword", user.PasswordHash)
		assert.Equal(t, "mail@example.com", user.Email)
	}
}

func TestNewUser_EmptyUsername_ReturnsError(t *testing.T) {
	user, err := NewUser("", "hashedpassword", "mail@example.com")

	assert.Nil(t, user)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "username cannot be empty")
	}
}

func TestNewUser_EmptyPasswordHash_ReturnsError(t *testing.T) {
	user, err := NewUser("testuser", "", "mail@example.com")

	assert.Nil(t, user)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "password hash cannot be empty")
	}
}

func TestNewUser_EmptyEmail_ReturnsError(t *testing.T) {
	user, err := NewUser("testuser", "hashedpassword", "")

	assert.Nil(t, user)
	if assert.Error(t, err) {
		assert.Contains(t, err.Error(), "email cannot be empty")
	}
}
