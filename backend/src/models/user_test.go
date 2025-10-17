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
	})

	t.Run("Creates user with empty fields", func(t *testing.T) {
		user := NewUser("", "", "")

		assert.Equal(t, "", user.Username)
		assert.Equal(t, "", user.PasswordHash)
		assert.Equal(t, "", user.Email)
	})
}
