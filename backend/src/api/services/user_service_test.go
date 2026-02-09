package services

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// setupTestUserService creates a new UserServiceT with a mock database for testing.
func setupTestUserService() (context.Context, *UserServiceT, *db.MockDB) {
	authService := newAuthorizationService()
	ctx := context.Background()
	mockDB := &db.MockDB{
		UsersRepo: new(db.MockUserRepository),
	}
	db.Set(mockDB)
	base := newBaseService(authService, "UserServiceTest")
	return ctx, newUserService(base), mockDB
}

func TestUserServiceT_RegisterUser(t *testing.T) {
	t.Run("Registers a new user successfully", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("InsertUser", ctx, mock.AnythingOfType("*models.User")).Return(&struct{}{}, nil)

		payload := &UserRegisterPayload{
			Username: "testuser",
			Password: "password123",
		}

		err := userService.RegisterUser(ctx, payload)

		assert.Nil(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Fails to register an invalid user", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("InsertUser", ctx, mock.AnythingOfType("*models.User")).Return(&struct{}{}, nil)

		payload := &UserRegisterPayload{
			Username: "",
			Password: "",
		}

		err := userService.RegisterUser(ctx, payload)

		assert.Equal(t, http.StatusBadRequest, err.Code)
		mockUserRepo.AssertNotCalled(t, "InsertUser")
	})

	t.Run("Fails to register user with password too long", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("InsertUser", ctx, mock.AnythingOfType("*models.User")).Return(&struct{}{}, nil)

		payload := &UserRegisterPayload{
			Username: "username",
			Password: "vu_dLvxG=d8?fd4bpMZHML$nMX:J_fTRjw{d1SUS=(EG*VL*Ffy]n*-.t=zUUfz1q3G]TxSH123",
		}

		err := userService.RegisterUser(ctx, payload)

		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockUserRepo.AssertNotCalled(t, "InsertUser")
	})

	t.Run("Fails to register user when database returns error", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("InsertUser", ctx, mock.AnythingOfType("*models.User")).Return(&struct{}{}, errors.New("database error"))

		payload := &UserRegisterPayload{
			Username: "username",
			Password: "password",
		}

		err := userService.RegisterUser(ctx, payload)

		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Fails to register user when user is duplicated", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("InsertUser", ctx, mock.AnythingOfType("*models.User")).Return(&struct{}{}, db.ErrAlreadyExists)

		payload := &UserRegisterPayload{
			Username: "username",
			Password: "password",
		}

		err := userService.RegisterUser(ctx, payload)

		assert.Equal(t, http.StatusConflict, err.Code)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserServiceT_GetAllUsers(t *testing.T) {
	t.Run("Returns all users successfully", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetAllUsers", ctx).Return([]models.User{
			{Username: "user1"},
			{Username: "user2"},
		}, nil)

		users, err := userService.GetAllUsers(ctx)

		assert.Contains(t, users, models.User{Username: "user1"})
		assert.Contains(t, users, models.User{Username: "user2"})
		assert.Nil(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Fails to get all users when database returns error", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetAllUsers", ctx).Return([]models.User{}, errors.New(""))

		users, err := userService.GetAllUsers(ctx)

		assert.Nil(t, users)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
	})
}

func TestUserServiceT_GetUserByUsername(t *testing.T) {
	t.Run("Returns user successfully", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "testuser").Return(&models.User{
			Username: "testuser",
		}, nil)

		user, err := userService.GetUserByUsername(ctx, "testuser")

		assert.NotNil(t, user)
		assert.Equal(t, "testuser", user.Username)
		assert.Nil(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Fails to get user when user not found", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "nonexistent").Return((*models.User)(nil), db.ErrNotFound)

		user, err := userService.GetUserByUsername(ctx, "nonexistent")

		assert.Nil(t, user)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Fails to get user when database returns error", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "testuser").Return((*models.User)(nil), errors.New("database error"))

		user, err := userService.GetUserByUsername(ctx, "testuser")

		assert.Nil(t, user)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUserServiceT_Login(t *testing.T) {
	t.Run("Login succeeds with correct credentials", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		t.Setenv("JWT_SECRET", "test-secret-key")
		t.Setenv("JWT_EXPIRY_HOURS", "24")

		hashedPassword, _ := hashPassword("password123")
		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "testuser").Return(&models.User{
			Username:     "testuser",
			PasswordHash: hashedPassword,
		}, nil)

		payload := LoginPayload{
			Username: "testuser",
			Password: "password123",
		}

		response, err := userService.Login(ctx, payload)

		assert.NotNil(t, response)
		assert.NotEmpty(t, response.Jwt)
		assert.Nil(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Login fails when user not found", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "nonexistent").Return((*models.User)(nil), db.ErrNotFound)

		payload := LoginPayload{
			Username: "nonexistent",
			Password: "password123",
		}

		response, err := userService.Login(ctx, payload)

		assert.Nil(t, response)
		assert.Equal(t, http.StatusNotFound, err.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Login fails with incorrect password", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		hashedPassword, _ := hashPassword("correctpassword")
		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "testuser").Return(&models.User{
			Username:     "testuser",
			PasswordHash: hashedPassword,
		}, nil)

		payload := LoginPayload{
			Username: "testuser",
			Password: "wrongpassword",
		}

		response, err := userService.Login(ctx, payload)

		assert.Nil(t, response)
		assert.Equal(t, http.StatusUnauthorized, err.Code)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Login fails when database returns error", func(t *testing.T) {
		ctx, userService, mockDB := setupTestUserService()

		mockUserRepo := mockDB.UsersRepo.(*db.MockUserRepository)
		mockUserRepo.On("GetUserByUsername", ctx, "testuser").Return((*models.User)(nil), errors.New("database error"))

		payload := LoginPayload{
			Username: "testuser",
			Password: "password123",
		}

		response, err := userService.Login(ctx, payload)

		assert.Nil(t, response)
		assert.Equal(t, http.StatusInternalServerError, err.Code)
		mockUserRepo.AssertExpectations(t)
	})
}
