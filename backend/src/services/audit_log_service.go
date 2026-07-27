package services

import (
	"context"

	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
)

type IAuditLogger interface {
	GetEntries(
		ctx context.Context,
		filter security.AuditLogFilter,
		pagination util.Pagination,
	) ([]security.AuditLogEntry, *ServiceError)
	Record(ctx context.Context, entry security.AuditLogEntry) error
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
	filter security.AuditLogFilter,
	pagination util.Pagination,
) ([]security.AuditLogEntry, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameAuditLog, "GetEntries")
	logger.Trace().Interface("filter", filter).Interface("pagination", pagination).Msg("Retrieving audit log entries")

	entries, dbErr := s.db.AuditLog().GetAuditLogEntries(ctx, filter, pagination)
	if dbErr != nil {
		logger.Error().Err(dbErr).Msg("Failed to retrieve audit log entries")
		return nil, NewInternalError("failed to retrieve audit log entries: %w", dbErr)
	}

	logger.Debug().Int("entryCount", len(entries)).Msg("Successfully retrieved audit log entries")
	return entries, nil
}

func (s *AuditLogService) Record(ctx context.Context, entry security.AuditLogEntry) error {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameAuditLog, "Record")
	logger.Trace().Interface("entry", entry).Msg("Recording audit log entry")

	entry.BeforeCreate()

	_, err := s.db.AuditLog().InsertAuditLogEntry(ctx, entry)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to save audit log entry")
		return err
	}

	logger.Debug().Msg("Audit log entry recorded successfully")
	return nil
}
