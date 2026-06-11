package services

import (
	"context"

	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
)

type IAuditLogger interface {
	GetEntries(ctx context.Context, userClaims *auth.UserClaims, filter security.AuditLogFilter, pagination util.Pagination) ([]security.AuditLogEntry, *ServiceError)
	Record(ctx context.Context, entry security.AuditLogEntry) error
}

// AuditLogService provides methods to manage audit log entries, including retrieval and recording of actions for auditing purposes.
type AuditLogService struct {
	db          db.DB
	authService IAuthorizer
}

type AuditLogServiceDeps struct {
	DB          db.DB
	AuthService IAuthorizer
}

// NewAuditLogService creates a new instance of AuditLogService with the provided dependencies.
func NewAuditLogService(deps AuditLogServiceDeps) AuditLogService {
	return AuditLogService{
		db:          deps.DB,
		authService: deps.AuthService,
	}
}

func (s *AuditLogService) GetEntries(ctx context.Context, userClaims *auth.UserClaims, filter security.AuditLogFilter, pagination util.Pagination) ([]security.AuditLogEntry, *ServiceError) {
	logger := log.FromContext(ctx).With().Str("method", "GetEntries").Logger()

	authErr := s.authorizeReadAccess(ctx, userClaims, filter)
	if authErr != nil {
		return nil, authErr
	}

	entries, dbErr := s.db.AuditLog().GetAuditLogEntries(ctx, filter, pagination)
	if dbErr != nil {
		logger.Error().Err(dbErr).Msg("Failed to retrieve audit log entries")
		return nil, NewInternalError("failed to retrieve audit log entries: %w", dbErr)
	}

	logger.Trace().Int("entryCount", len(entries)).Msg("Successfully retrieved audit log entries")
	return entries, nil
}

func (s *AuditLogService) authorizeReadAccess(ctx context.Context, claims *auth.UserClaims, filter security.AuditLogFilter) *ServiceError {
	logger := MethodLoggerFromContext(ctx, "AuditLogService", "authorizeReadAccess")

	isInstanceAdmin, err := s.authService.isInstanceAdmin(ctx, claims.Username)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to check instance admin status for authorizeReadAccess")
		return NewInternalError("failed to check admin status: %w", err)
	}

	if isInstanceAdmin {
		logger.Debug().Str("username", claims.Username).Msg("User is instance admin, skipping project authorization for authorizeReadAccess")
		return nil
	}

	if err := filter.ValidateForNonInstanceAdmin(); err != nil {
		logger.Warn().Err(err).Msg("Invalid filter for non-instance admin user in authorizeReadAccess")
		return NewBadRequestError("invalid filter: %w", err)
	}

	_, authErr := s.authService.authorizeProjectAction(ctx, &authorization.ProjectAuthorization{
		Username:  claims.Username,
		ProjectID: *filter.ProjectID,
	}, models.PermissionProjectAdmin)
	if authErr != nil {
		logger.Warn().Err(authErr).Msg("Unauthorized access to authorizeReadAccess")
		return authErr
	}

	logger.Trace().Str("username", claims.Username).Str("projectID", filter.ProjectID.String()).Msg("User authorized for read access in authorizeReadAccess")
	return nil
}

func (s *AuditLogService) Record(ctx context.Context, entry security.AuditLogEntry) error {
	logger := MethodLoggerFromContext(ctx, "AuditLogService", "Record")
	entry.BeforeCreate()

	_, err := s.db.AuditLog().InsertAuditLogEntry(ctx, entry)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to save audit log entry")
		return err
	}
	return nil
}
