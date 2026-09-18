package auditlog

import (
	"context"

	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
)

type IAuditLogger interface {
	GetEntries(
		ctx context.Context,
		filter audit.AuditLogFilter,
		pagination util.Pagination,
	) ([]audit.AuditLogEntry, *apperr.ServiceError)
	Record(ctx context.Context, params audit.AuditLogParams) error
}

// AuditLogService provides methods to manage audit log entries, including retrieval and recording of actions for auditing purposes.
type AuditLogService struct {
	db db.DB
}

type AuditLogServiceDeps struct {
	DB db.DB
}

// NewAuditLogService creates a new instance of AuditLogService with the provided dependencies.
func NewAuditLogService(deps AuditLogServiceDeps) AuditLogService {
	return AuditLogService{
		db: deps.DB,
	}
}

func (s *AuditLogService) GetEntries(
	ctx context.Context,
	filter audit.AuditLogFilter,
	pagination util.Pagination,
) ([]audit.AuditLogEntry, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameAuditLog, "GetEntries")
	logger.Trace().Interface("filter", filter).Interface("pagination", pagination).Msg("Retrieving audit log entries")

	entries, dbErr := audit.NewAuditLogDAO(s.db.Querier()).GetAuditLogEntries(ctx, filter, pagination)
	if dbErr != nil {
		logger.Error().Err(dbErr).Msg("Failed to retrieve audit log entries")
		return nil, apperr.NewInternalError("failed to retrieve audit log entries: %w", dbErr)
	}

	logger.Debug().Int("entryCount", len(entries)).Msg("Successfully retrieved audit log entries")
	return entries, nil
}

func (s *AuditLogService) Record(ctx context.Context, params audit.AuditLogParams) error {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameAuditLog, "Record")
	logger.Trace().Interface("params", params).Msg("Recording audit log entry")

	err := audit.NewAuditLogDAO(s.db.Querier()).Record(ctx, params)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to save audit log entry")
		return err
	}

	logger.Debug().Msg("Audit log entry recorded successfully")
	return nil
}
