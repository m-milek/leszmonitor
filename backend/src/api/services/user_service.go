package services

import (
	"context"
	"errors"
	"fmt"
	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"os"
	"strconv"
	"time"
)

type UserServiceT struct {
	BaseService
}

// NewUserService creates a new instance of UserServiceT.
func newUserService() *UserServiceT {
	return &UserServiceT{
		BaseService{
			serviceLogger: logging.NewServiceLogger("UserService"),
		},
	}
}

var UserService = newUserService()

type UserRegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginPayload struct {
	jwt2.MapClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

func (s *UserServiceT) GetAllUsers(ctx context.Context) (result []models.UserResponse, error *ServiceError) {
	logger := s.getMethodLogger("GetAllUsers")
	logger.Trace().Msg("Retrieving all users")

	users, err := db.GetAllUsers(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Error retrieving users")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving users: %w", err),
		}
	}

	return users, nil
}

func (s *UserServiceT) GetUserByUsername(ctx context.Context, username string) (*models.UserResponse, *ServiceError) {
	logger := s.getMethodLogger("GetUserByUsername")
	logger.Trace().Str("username", username).Msg("Retrieving user by username")

	user, err := db.GetUserByUsername(ctx, username)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("username", username).Msg("User not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("user %s not found", username),
			}
		}
		logger.Error().Err(err).Str("username", username).Msg("Error retrieving user")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving user %s: %w", username, err),
		}
	}

	return user, nil
}

func (s *UserServiceT) RegisterUser(ctx context.Context, payload *UserRegisterPayload) *ServiceError {
	logger := s.getMethodLogger("RegisterUser")
	logger.Trace().Str("username", payload.Username).Msg("Registering new user")

	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to hash password")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to hash password: %w", err),
		}
	}

	user := models.NewUser(payload.Username, hashedPassword, payload.Email)

	_, err = db.CreateUser(ctx, user)

	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to create user in database")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to register user %s: %w", payload.Username, err),
		}
	}

	return nil
}

func (s *UserServiceT) Login(ctx context.Context, payload LoginPayload) (*LoginResponse, *ServiceError) {
	logger := s.getMethodLogger("Login")
	logger.Trace().Str("username", payload.Username).Msg("User login attempt")

	user, err := db.GetRawUserByUsername(ctx, payload.Username)

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("username", payload.Username).Msg("User not found during login")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("user %s not found", payload.Username),
			}
		}
		logger.Error().Err(err).Str("username", payload.Username).Msg("Error retrieving user during login")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving user %s: %w", payload.Username, err),
		}
	}

	matches := checkPasswordHash(payload.Password, user.PasswordHash)
	if !matches {
		logger.Warn().Str("username", payload.Username).Msg("Invalid password during login")
		return nil, &ServiceError{
			Code: http.StatusUnauthorized,
			Err:  fmt.Errorf("invalid password for user %s", payload.Username),
		}
	}

	expiryHours, err := strconv.Atoi(os.Getenv(env.JwtExpiryHours))
	if err != nil {
		logger.Error().Err(err).Msg("Invalid JwtExpiryHours environment variable")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("invalid JwtExpiryHours value: %w", err),
		}
	}
	validFor := time.Duration(expiryHours) * time.Hour
	expiryDate := time.Now().Add(validFor)

	jwt := jwt2.NewWithClaims(
		jwt2.SigningMethodHS256,
		jwt2.MapClaims{
			"username": payload.Username,
			"exp":      jwt2.NewNumericDate(expiryDate),
			"iat":      jwt2.NewNumericDate(time.Now()),
		},
	)
	token, err := jwt.SignedString([]byte(os.Getenv(env.JwtSecret)))

	if err != nil {
		logger.Error().Err(err).Msg("Failed to create JWT token")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create JWT token: %w", err),
		}
	}

	return &LoginResponse{Jwt: token}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
