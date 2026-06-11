package services

import (
	"net/http"
	"testing"

	"github.com/m-milek/leszmonitor/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegration_UserService_RegisterUser(t *testing.T) {
	t.Run("Successfully registers a new user and creates sandbox project", func(t *testing.T) {
		ctx, _, userService, _ := setupIntegrationTest(t)

		payload := &UserRegisterPayload{
			Username:        "new_user",
			Password:        "Password123!",
			PasswordConfirm: "Password123!",
		}

		err := userService.RegisterUser(ctx, payload)
		require.Nil(t, err)

		// Verify user was created in DB
		user, getErr := userService.GetUserByUsername(ctx, "new_user")
		require.Nil(t, getErr)
		assert.Equal(t, "new_user", user.Username)

		// Verify sandbox project was auto-created
		projects, projErr := db.Get().Projects().GetProjectsByQuery(ctx, db.GetProjectsQuery{RequestingUserID: user.ID})
		require.NoError(t, projErr)
		require.Len(t, projects, 1)
		assert.Equal(t, "new_user's Sandbox", projects[0].Name)
	})

	t.Run("Fails to register a duplicate user", func(t *testing.T) {
		ctx, _, userService, owner := setupIntegrationTest(t)

		payload := &UserRegisterPayload{
			Username:        owner.Username, // Already registered by setupIntegrationTest
			Password:        "Password123!",
			PasswordConfirm: "Password123!",
		}

		err := userService.RegisterUser(ctx, payload)
		require.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.Code) // Service returns 401 for duplicate registration
	})
}

func TestIntegration_UserService_Login(t *testing.T) {
	t.Run("Successfully logs in with valid credentials", func(t *testing.T) {
		ctx, _, userService, _ := setupIntegrationTest(t)

		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{
			Username:        "login_user",
			Password:        "MySecretPassword!",
			PasswordConfirm: "MySecretPassword!",
		}))

		payload := LoginPayload{
			Username: "login_user",
			Password: "MySecretPassword!",
		}

		resp, err := userService.Login(ctx, payload)
		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.NotEmpty(t, resp.Jwt)
	})

	t.Run("Fails to log in with invalid password", func(t *testing.T) {
		ctx, _, userService, _ := setupIntegrationTest(t)

		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{
			Username:        "login_user2",
			Password:        "MySecretPassword!",
			PasswordConfirm: "MySecretPassword!",
		}))

		payload := LoginPayload{
			Username: "login_user2",
			Password: "WrongPassword!",
		}

		resp, err := userService.Login(ctx, payload)
		require.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.Code)
		assert.Nil(t, resp)
	})

	t.Run("Fails to log in with nonexistent user", func(t *testing.T) {
		ctx, _, userService, _ := setupIntegrationTest(t)

		payload := LoginPayload{
			Username: "nonexistent",
			Password: "Password!",
		}

		resp, err := userService.Login(ctx, payload)
		require.NotNil(t, err)
		assert.Equal(t, http.StatusUnauthorized, err.Code)
		assert.Nil(t, resp)
	})
}

func TestIntegration_UserService_GetUserByUsername(t *testing.T) {
	t.Run("Successfully retrieves an existing user", func(t *testing.T) {
		ctx, _, userService, owner := setupIntegrationTest(t)

		user, err := userService.GetUserByUsername(ctx, owner.Username)
		require.Nil(t, err)
		assert.Equal(t, owner.Username, user.Username)
		assert.Equal(t, owner.ID, user.ID)
	})

	t.Run("Fails to retrieve a nonexistent user", func(t *testing.T) {
		ctx, _, userService, _ := setupIntegrationTest(t)

		user, err := userService.GetUserByUsername(ctx, "nobody")
		require.NotNil(t, err)
		assert.Equal(t, http.StatusNotFound, err.Code)
		assert.Nil(t, user)
	})
}

func TestIntegration_UserService_GetAllUsers(t *testing.T) {
	t.Run("Successfully retrieves all users", func(t *testing.T) {
		ctx, _, userService, owner := setupIntegrationTest(t)

		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{
			Username:        "user1",
			Password:        "Pass1!",
			PasswordConfirm: "Pass1!",
		}))
		require.Nil(t, userService.RegisterUser(ctx, &UserRegisterPayload{
			Username:        "user2",
			Password:        "Pass2!",
			PasswordConfirm: "Pass2!",
		}))

		users, err := userService.GetAllUsers(ctx)
		require.Nil(t, err)
		// owner + user1 + user2 = at least 3 users
		assert.GreaterOrEqual(t, len(users), 3)

		// Verify all these users are in the list
		foundOwner, foundUser1, foundUser2 := false, false, false
		for _, u := range users {
			if u.Username == owner.Username {
				foundOwner = true
			}
			if u.Username == "user1" {
				foundUser1 = true
			}
			if u.Username == "user2" {
				foundUser2 = true
			}
		}

		assert.True(t, foundOwner)
		assert.True(t, foundUser1)
		assert.True(t, foundUser2)
	})
}
