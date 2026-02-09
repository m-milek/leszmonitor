package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/env"
	"github.com/m-milek/leszmonitor/models"
	"golang.org/x/crypto/bcrypt"
)

// UserServiceT handles user-related operations such as registration, login, and retrieval.
type UserServiceT struct {
	baseService
}

// NewUserService creates a new instance of UserServiceT.
func newUserService(base baseService) *UserServiceT {
	return &UserServiceT{
		baseService: base,
	}
}

var UserService = newUserService(newBaseService(newAuthorizationService(), "UserService"))

type UserRegisterPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginPayload struct {
	jwt2.MapClaims
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Jwt string `json:"jwt"`
}

// GetAllUsers retrieves all users from the database.
// Requires no permissions.
func (s *UserServiceT) GetAllUsers(ctx context.Context) ([]models.User, *ServiceError) {
	logger := s.getMethodLogger("GetAllUsers")
	logger.Trace().Msg("Retrieving all users")

	users, err := s.getDB().Users().GetAllUsers(ctx)

	if err != nil {
		logger.Error().Err(err).Msg("Error retrieving users")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("error retrieving users: %w", err),
		}
	}

	return users, nil
}

// GetUserByUsername retrieves a user by their username.
// Requires no permissions.
func (s *UserServiceT) GetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError) {
	return s.internalGetUserByUsername(ctx, username)
}

// internalGetUserByUsername retrieves a user by their username without authorization checks.
// This is used internally by other services to avoid circular dependencies.
func (s *UserServiceT) internalGetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError) {
	logger := s.getMethodLogger("internalGetUserByUsername")
	logger.Trace().Str("username", username).Msg("Retrieving user by username")

	user, err := s.getDB().Users().GetUserByUsername(ctx, username)

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

// RegisterUser registers a new user with the provided payload.
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

	userModel, err := models.NewUser(payload.Username, hashedPassword)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Invalid user data")
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid user data for %s: %w", payload.Username, err),
		}
	}

	user, err := s.getDB().Users().InsertUser(ctx, userModel)

	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to create user in database")
		if errors.Is(err, db.ErrAlreadyExists) {
			return &ServiceError{
				Code: http.StatusConflict,
				Err:  fmt.Errorf("user %s already exists", payload.Username),
			}
		}
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to register user %s: %w", payload.Username, err),
		}
	}

	teamName := fmt.Sprintf("%s's Space", payload.Username)
	teamDescription := fmt.Sprintf("Default space for %s", payload.Username)

	team, teamErr := models.NewTeam(teamName, teamDescription, user.ID)
	if teamErr != nil {
		logger.Error().Err(teamErr).Str("username", payload.Username).Msg("Failed to create default team model for user")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create default team for user %s: %w", payload.Username, teamErr),
		}
	}

	team, teamServiceErr := TeamService.internalCreateTeam(ctx, team)
	if teamServiceErr != nil {
		logger.Error().Err(teamServiceErr.Err).Str("username", payload.Username).Msg("Failed to create default team for user")
		return &ServiceError{
			Code: teamServiceErr.Code,
			Err:  fmt.Errorf("failed to create default team for user %s: %w", payload.Username, teamServiceErr.Err),
		}
	}

	logger.Trace().Str("username", payload.Username).Str("team", team.Name).Msg("Successfully registered user and created default team")

	return nil
}

// Login authenticates a user and returns a JWT token if successful.
// On success, returns a LoginResponse containing the JWT token.
func (s *UserServiceT) Login(ctx context.Context, payload LoginPayload) (*LoginResponse, *ServiceError) {
	logger := s.getMethodLogger("Login")
	logger.Trace().Str("username", payload.Username).Msg("User login attempt")

	user, err := s.getDB().Users().GetUserByUsername(ctx, payload.Username)

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

	err = checkPasswordHash(payload.Password, user.PasswordHash)
	if err != nil {
		logger.Warn().Err(err).Str("username", payload.Username).Msg("Invalid password during login")
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

// hashPassword hashes a plaintext password using bcrypt.
// Returns the hashed password or an error if hashing fails.
func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// checkPasswordHash compares a plaintext password with a hashed password.
// Returns true if they match, false otherwise.
func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
