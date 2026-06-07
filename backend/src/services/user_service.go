package services

import (
	"context"
	"errors"
	"fmt"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]models.User, *ServiceError)
	GetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError)
	RegisterUser(ctx context.Context, payload *UserRegisterPayload) *ServiceError
	Login(ctx context.Context, payload LoginPayload) (*LoginResponse, *ServiceError)
}

type UserServiceDeps struct {
	DB             db.DB
	Auth           IAuthorizer
	ProjectService IProjectService
}

// UserService handles user-related operations such as registration, login, and retrieval.
type UserService struct {
	db             db.DB
	auth           IAuthorizer
	projectService IProjectService
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		db:             deps.DB,
		auth:           deps.Auth,
		projectService: deps.ProjectService,
	}
}

type UserRegisterPayload struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"passwordConfirm"`
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
func (s *UserService) GetAllUsers(ctx context.Context) ([]models.User, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "UserService", "GetAllUsers")
	logger.Trace().Msg("Retrieving all users")

	users, err := s.db.Users().GetAllUsers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error retrieving users")
		return nil, NewInternalError("error retrieving users: %w", err)
	}

	return users, nil
}

// GetUserByUsername retrieves a user by their username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError) {
	return s.internalGetUserByUsername(ctx, username)
}

// internalGetUserByUsername retrieves a user by their username without authorization checks.
func (s *UserService) internalGetUserByUsername(ctx context.Context, username string) (*models.User, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "UserService", "internalGetUserByUsername")
	logger.Trace().Str("username", username).Msg("Retrieving user by username")

	user, err := s.db.Users().GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("username", username).Msg("User not found")
			return nil, NewNotFoundError("user %s not found", username)
		}
		logger.Error().Err(err).Str("username", username).Msg("Error retrieving user")
		return nil, NewInternalError("error retrieving user %s: %w", username, err)
	}

	return user, nil
}

// RegisterUser registers a new user with the provided payload.
func (s *UserService) RegisterUser(ctx context.Context, payload *UserRegisterPayload) *ServiceError {
	logger := MethodLoggerFromContext(ctx, "UserService", "RegisterUser")
	logger.Trace().Str("username", payload.Username).Msg("Registering new user")

	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to hash password")
		return NewInternalError("failed to hash password: %w", err)
	}

	userModel, err := models.NewUser(payload.Username, hashedPassword)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Invalid user data")
		return NewBadRequestError("invalid user data for %s: %w", payload.Username, err)
	}

	_, err = s.db.Users().InsertUser(ctx, userModel)
	if err != nil {
		if errors.Is(err, db.ErrAlreadyExists) {
			return NewUnauthorizedError("failed to register user")
		}
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to create user in database")
		return NewInternalError("failed to register user %s: %w", payload.Username, err)
	}

	logger.Trace().Str("username", payload.Username).Msg("User registered successfully")

	_, projectErr := s.projectService.CreateProject(ctx, payload.Username, CreateProjectPayload{
		Name:        fmt.Sprintf("%s's Sandbox", payload.Username),
		Description: "Your default sandbox project",
	})
	if projectErr != nil {
		logger.Error().Err(projectErr.Err).Msg("Failed to auto-create sandbox project for new user")
		// We don't fail the registration if the sandbox fails
	}

	return nil
}

// Login authenticates a user and returns a JWT token if successful.
func (s *UserService) Login(ctx context.Context, payload LoginPayload) (*LoginResponse, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "UserService", "Login")
	logger.Info().Str("username", payload.Username).Msg("User login attempt")

	user, err := s.db.Users().GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, NewUnauthorizedError("invalid credentials")
		}
		return nil, NewInternalError("error retrieving user %s: %w", payload.Username, err)
	}

	if err = checkPasswordHash(payload.Password, user.PasswordHash); err != nil {
		return nil, NewUnauthorizedError("invalid credentials")
	}

	jwtToken, err := auth.NewJwt(payload.Username)
	if jwtToken == nil {
		logger.Error().Str("username", payload.Username).Err(err).Msg("Failed to generate JWT token")
		return nil, NewInternalError("failed to generate JWT token")
	}

	return &LoginResponse{Jwt: *jwtToken}, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
