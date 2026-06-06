package services

import (
	"context"
	"net/http"

	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/auth"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
)

type IAuditLogService interface {
	GetEntries(ctx context.Context, userClaims *auth.UserClaims, filter security.AuditLogFilter, pagination util.Pagination) ([]security.AuditLogEntry, *ServiceError)
	Record(ctx context.Context, entry security.AuditLogEntry) error
}

// authorizationServiceT handles authorization-related operations.
// It provides methods to authorize actions based on project membership and permissions.
type auditLogServiceT struct {
	baseService
}

// newAuthorizationService creates a new instance of authorizationServiceT.
func newAuditLogService() *auditLogServiceT {
	return &auditLogServiceT{
		baseService: baseService{
			serviceLogger: log.NewServiceLogger("AuditLogService"),
			authService:   newAuthorizationService(),
		},
	}
}

var AuditLogService = newAuditLogService()

func (s *auditLogServiceT) GetEntries(ctx context.Context, userClaims *auth.UserClaims, filter security.AuditLogFilter, pagination util.Pagination) ([]security.AuditLogEntry, *ServiceError) {
	logger := log.FromContext(ctx).With().Str("method", "GetEntries").Logger()

	authErr := s.authorizeReadAccess(ctx, userClaims, filter)
	if authErr != nil {
		return nil, authErr
	}

	entries, dbErr := s.getDB().AuditLog().GetAuditLogEntries(ctx, filter, pagination)
	if dbErr != nil {
		logger.Error().Err(dbErr).Msg("Failed to retrieve audit log entries")
		return nil, newServiceError(http.StatusInternalServerError, dbErr)
	}

	logger.Trace().Int("entryCount", len(entries)).Msg("Successfully retrieved audit log entries")
	return entries, nil
}

func (s *auditLogServiceT) authorizeReadAccess(ctx context.Context, claims *auth.UserClaims, filter security.AuditLogFilter) *ServiceError {
	logger := s.getMethodLogger("authorizeReadAccess")

	isInstanceAdmin, err := s.authService.isInstanceAdmin(ctx, claims.Username)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to check instance admin status for authorizeReadAccess")
		return newServiceError(http.StatusInternalServerError, err)
	}

	if isInstanceAdmin {
		logger.Debug().Str("username", claims.Username).Msg("User is instance admin, skipping project authorization for authorizeReadAccess")
		return nil
	}

	if err := filter.ValidateForNonInstanceAdmin(); err != nil {
		logger.Warn().Err(err).Msg("Invalid filter for non-instance admin user in authorizeReadAccess")
		return newServiceError(http.StatusBadRequest, err)
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

func (s *auditLogServiceT) Record(ctx context.Context, entry security.AuditLogEntry) error {
	entry.BeforeCreate()

	_, err := s.getDB().AuditLog().InsertAuditLogEntry(ctx, entry)
	if err != nil {
		logger := s.getMethodLogger("SaveAction")
		logger.Error().Err(err).Msg("Failed to save audit log entry")
		return err
	}
	return nil
}
