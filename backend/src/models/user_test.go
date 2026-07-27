package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewUser_Success(t *testing.T) {
	t.Parallel()
	user, err := NewUser("testuser", "hashedpassword")

	require.NoError(t, err)
	if assert.NotNil(t, user) {
		assert.Equal(t, "testuser", user.Username)
		assert.Equal(t, "hashedpassword", user.PasswordHash)
	}
}

func TestNewUser_EmptyUsername_ReturnsError(t *testing.T) {
	t.Parallel()
	user, err := NewUser("", "hashedpassword")

	assert.Nil(t, user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "username cannot be empty")
}

func TestNewUser_EmptyPasswordHash_ReturnsError(t *testing.T) {
	t.Parallel()
	user, err := NewUser("testuser", "")

	assert.Nil(t, user)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "password hash cannot be empty")
}
