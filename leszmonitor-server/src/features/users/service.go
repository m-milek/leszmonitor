package users

import (
	"context"
	"errors"
	"fmt"
	"os"

	jwt2 "github.com/golang-jwt/jwt/v5"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	config "github.com/m-milek/leszmonitor/platform/config"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"golang.org/x/crypto/bcrypt"
)

type IUserService interface {
	GetAllUsers(ctx context.Context) ([]User, *apperr.ServiceError)
	GetUserByUsername(ctx context.Context, username string) (*User, *apperr.ServiceError)
	RegisterUser(ctx context.Context, payload *UserRegisterPayload) *apperr.ServiceError
	Login(ctx context.Context, payload LoginPayload) (*LoginResponse, *apperr.ServiceError)
	SetUserRole(ctx context.Context, username string, payload SetUserRolePayload) (*User, *apperr.ServiceError)
}

type SetUserRolePayload struct {
	Role auth.Role `json:"role"`
}

type UserServiceDeps struct {
	DB db.DB
}

// UserService handles user-related operations such as registration, login, and retrieval.
type UserService struct {
	db db.DB
}

func NewUserService(deps UserServiceDeps) *UserService {
	return &UserService{
		db: deps.DB,
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
func (s *UserService) GetAllUsers(ctx context.Context) ([]User, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "GetAllUsers")
	logger.Trace().Msg("Retrieving all users")

	users, err := NewUserDAO(s.db.Querier()).GetAllUsers(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Error retrieving users")
		return nil, apperr.NewInternalError("error retrieving users: %w", err)
	}

	logger.Debug().Int("count", len(users)).Msg("Successfully retrieved all users")
	return users, nil
}

// GetUserByUsername retrieves a user by their username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*User, *apperr.ServiceError) {
	return s.internalGetUserByUsername(ctx, username)
}

// internalGetUserByUsername retrieves a user by their username without authorization checks.
func (s *UserService) internalGetUserByUsername(ctx context.Context, username string) (*User, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "internalGetUserByUsername")
	logger.Trace().Str("username", username).Msg("Retrieving user by username")

	user, err := NewUserDAO(s.db.Querier()).GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", username).Msg("User not found")
			return nil, apperr.NewNotFoundError("user %s not found", username)
		}
		logger.Error().Err(err).Str("username", username).Msg("Error retrieving user")
		return nil, apperr.NewInternalError("error retrieving user %s: %w", username, err)
	}

	logger.Debug().Str("username", username).Msg("Successfully retrieved user")
	return user, nil
}

// SetUserRole updates the role of the user identified by username.
func (s *UserService) SetUserRole(
	ctx context.Context,
	username string,
	payload SetUserRolePayload,
) (*User, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "SetUserRole")
	logger.Trace().Str("username", username).Str("role", string(payload.Role)).Msg("Setting user role")

	if err := payload.Role.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid role")
		return nil, apperr.NewBadRequestError("invalid role: %w", err)
	}

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewUnauthorizedError("user claims not found in context")
	}

	existingUser, getErr := s.internalGetUserByUsername(ctx, username)
	if getErr != nil {
		return nil, getErr
	}

	if existingUser.Role == payload.Role {
		logger.Debug().Str("username", username).Msg("User already has the requested role")
		return existingUser, nil
	}

	updatedUser, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*User, *audit.AuditLogParams, error) {
		u, err := NewUserDAO(q).UpdateUserRole(ctx, existingUser.ID, payload.Role)
		if err != nil {
			return nil, nil, err
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &u.ID,
			Action:     audit.ActionUpdateUser,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("User %s role changed from %s to %s", username, existingUser.Role, payload.Role),
			Before:     existingUser,
			After:      u,
		}
		return u, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrNotFound) {
			return nil, apperr.NewNotFoundError("user %s not found", username)
		}
		logger.Error().Err(txErr).Str("username", username).Msg("Failed to update user role")
		return nil, apperr.NewInternalError("failed to update user role: %w", txErr)
	}

	logger.Debug().Str("username", username).Str("role", string(payload.Role)).Msg("User role updated successfully")
	return updatedUser, nil
}

// RegisterUser registers a new user with the provided payload.
func (s *UserService) RegisterUser(ctx context.Context, payload *UserRegisterPayload) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "RegisterUser")
	logger.Trace().Str("username", payload.Username).Msg("Registering new user")

	hashedPassword, err := hashPassword(payload.Password)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Failed to hash password")
		return apperr.NewInternalError("failed to hash password: %w", err)
	}

	userModel, err := NewUser(payload.Username, hashedPassword)
	if err != nil {
		logger.Error().Err(err).Str("username", payload.Username).Msg("Invalid user data")
		return apperr.NewBadRequestError("invalid user data for %s: %w", payload.Username, err)
	}

	_, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*User, *audit.AuditLogParams, error) {
		u, err := NewUserDAO(q).InsertUser(ctx, userModel)
		if err != nil {
			return nil, nil, err
		}

		params := &audit.AuditLogParams{
			Username:   &payload.Username,
			ResourceID: &u.ID,
			Action:     audit.ActionCreateUser,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("User %s registered", u.Username),
			After:      u,
		}

		return u, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrAlreadyExists) {
			return apperr.NewUnauthorizedError("failed to register user")
		}
		logger.Error().Err(txErr).Str("username", payload.Username).Msg("Failed to create user in database")
		return apperr.NewInternalError("failed to register user %s: %w", payload.Username, txErr)
	}

	logger.Trace().Str("username", payload.Username).Msg("User registered successfully")

	return nil
}

// Login authenticates a user and returns a JWT token if successful.
func (s *UserService) Login(ctx context.Context, payload LoginPayload) (*LoginResponse, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "Login")
	logger.Info().Str("username", payload.Username).Msg("User login attempt")

	user, err := NewUserDAO(s.db.Querier()).GetUserByUsername(ctx, payload.Username)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("username", payload.Username).Msg("User not found for login")

			_ = audit.NewAuditLogDAO(s.db.Querier()).Record(ctx, audit.AuditLogParams{
				Username:  &payload.Username,
				Action:    audit.ActionFailedLogin,
				IsSuccess: false,
				Summary:   fmt.Sprintf("Failed login attempt for unknown user: %s", payload.Username),
			})

			return nil, apperr.NewUnauthorizedError("invalid credentials")
		}
		logger.Error().Err(err).Str("username", payload.Username).Msg("Error retrieving user for login")
		return nil, apperr.NewInternalError("error retrieving user %s: %w", payload.Username, err)
	}

	if err = checkPasswordHash(payload.Password, user.PasswordHash); err != nil {
		logger.Error().Str("username", payload.Username).Msg("Invalid password for login")

		_ = audit.NewAuditLogDAO(s.db.Querier()).Record(ctx, audit.AuditLogParams{
			Username:   &payload.Username,
			ResourceID: &user.ID,
			Action:     audit.ActionFailedLogin,
			IsSuccess:  false,
			Summary:    fmt.Sprintf("Failed login attempt for user: %s", payload.Username),
		})

		return nil, apperr.NewUnauthorizedError("invalid credentials")
	}

	jwtToken, err := auth.NewJwt(payload.Username, GetIsInstanceAdmin(*user))
	if jwtToken == nil {
		logger.Error().Str("username", payload.Username).Err(err).Msg("Failed to generate JWT token")
		return nil, apperr.NewInternalError("failed to generate JWT token")
	}

	_ = audit.NewAuditLogDAO(s.db.Querier()).Record(ctx, audit.AuditLogParams{
		Username:   &payload.Username,
		ResourceID: &user.ID,
		Action:     audit.ActionLogin,
		IsSuccess:  true,
		Summary:    fmt.Sprintf("User %s logged in", payload.Username),
	})

	logger.Debug().Str("username", payload.Username).Msg("Login successful")
	return &LoginResponse{Jwt: *jwtToken}, nil
}

func (s *UserService) EnsureAdminUserExists(ctx context.Context) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "EnsureAdminUserExists")
	logger.Info().Msg("Ensuring admin user exists")

	adminUser, err := NewUserDAO(s.db.Querier()).GetUserByUsername(ctx, os.Getenv(config.InstanceAdminUsername))
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Info().Msg("Admin user not found, creating...")
			return s.createAdminUser(ctx)
		}
		return apperr.NewInternalError("error occurred while checking for admin user: %w", err)
	}

	logger.Info().Str("username", adminUser.Username).Msg("Admin user already exists")

	return nil
}

func (s *UserService) createAdminUser(ctx context.Context) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameUser, "CreateAdminUser")

	hashedPassword, err := hashPassword(os.Getenv(config.InstanceAdminPassword))
	if err != nil {
		return apperr.NewInternalError("failed to hash admin password: %w", err)
	}

	adminUserModel, err := NewUser(os.Getenv(config.InstanceAdminUsername), hashedPassword)
	if err != nil {
		return apperr.NewBadRequestError("invalid admin user data: %w", err)
	}
	adminUserModel.IsInstanceAdmin = true
	adminUserModel.Role = auth.RoleOwner

	_, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*User, *audit.AuditLogParams, error) {
		u, err := NewUserDAO(q).InsertUser(ctx, adminUserModel)
		if err != nil {
			return nil, nil, err
		}

		params := &audit.AuditLogParams{
			Username:   &u.Username,
			ResourceID: &u.ID,
			Action:     audit.ActionCreateUser,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Admin user %s created", u.Username),
			After:      u,
		}

		return u, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrAlreadyExists) {
			logger.Info().Msg("Admin user already exists")
			return nil
		}
		return apperr.NewInternalError("failed to create admin user: %w", txErr)
	}

	logger.Info().Msg("Admin user created successfully")

	return nil
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
