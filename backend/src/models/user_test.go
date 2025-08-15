package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("Creates user with valid data", func(t *testing.T) {
		user := NewUser("testuser", "hashedpassword", "mail@example.com")

		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "hashedpassword", user.PasswordHash)
		assert.Equal(t, "mail@example.com", user.Email)
		assert.NotEmpty(t, user.CreatedAt)
		assert.NotEmpty(t, user.UpdatedAt)
	})

	t.Run("Creates user with empty fields", func(t *testing.T) {
		user := NewUser("", "", "")

		assert.Equal(t, "", user.Username)
		assert.Equal(t, "", user.PasswordHash)
		assert.Equal(t, "", user.Email)
		assert.NotEmpty(t, user.CreatedAt)
		assert.NotEmpty(t, user.UpdatedAt)
	})

	t.Run("IntoUser converts RawUser to UserResponse", func(t *testing.T) {
		rawUser := NewUser("testuser", "hashedpassword", "mail@example.com")

		userResponse := rawUser.IntoUser()

		assert.Equal(t, "testuser", userResponse.Username)
		assert.Equal(t, "mail@example.com", userResponse.Email)
		assert.NotEmpty(t, userResponse.CreatedAt)
		assert.NotEmpty(t, userResponse.CreatedAt)
	})
}
