package db

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/security"
	"github.com/m-milek/leszmonitor/util"
)

type IAuditLogRepository interface {
	InsertAuditLogEntry(ctx context.Context, entry security.AuditLogEntry) (any, error)
	GetAuditLogEntries(ctx context.Context, filter security.AuditLogFilter, pagination util.Pagination) ([]security.AuditLogEntry, error)
}

type auditLogRepository struct {
	baseRepository
}

func newAuditLogRepository(base baseRepository) IAuditLogRepository {
	return &auditLogRepository{
		baseRepository: base,
	}
}

func (a auditLogRepository) InsertAuditLogEntry(ctx context.Context, entry security.AuditLogEntry) (any, error) {
	return dbWrap(ctx, "InsertAuditLogEntry", func() (any, error) {
		_, err := a.pool.ExecContext(ctx,
			`INSERT INTO audit_logs (id, username, project_id, resource_id, action, is_success, summary, before, after, trace_id, created_at)
			 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
			entry.ID, entry.Username, entry.ProjectID, entry.ResourceID, entry.Action, entry.IsSuccess, entry.Summary, entry.Before, entry.After, entry.TraceID, entry.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
}

func (a auditLogRepository) GetAuditLogEntries(ctx context.Context, filter security.AuditLogFilter, pagination util.Pagination) ([]security.AuditLogEntry, error) {
	return dbWrap(ctx, "GetAuditLogEntries", func() ([]security.AuditLogEntry, error) {
		var (
			entries    []security.AuditLogEntry
			conditions []string
			args       []interface{}
		)
		if filter.UserID != nil {
			conditions = append(conditions, "user_id = ?")
			args = append(args, *filter.UserID)
		}
		if filter.ProjectID != nil {
			conditions = append(conditions, "project_id = ?")
			args = append(args, *filter.ProjectID)
		}
		if filter.ResourceID != nil {
			conditions = append(conditions, "resource_id = ?")
			args = append(args, *filter.ResourceID)
		}
		if filter.Action != nil {
			conditions = append(conditions, "action = ?")
			args = append(args, *filter.Action)
		}
		if filter.IsSuccess != nil {
			conditions = append(conditions, "is_success = ?")
			args = append(args, *filter.IsSuccess)
		}
		if filter.TraceID != nil {
			conditions = append(conditions, "trace_id = ?")
			args = append(args, *filter.TraceID)
		}
		if filter.StartDate != nil {
			conditions = append(conditions, "created_at >= ?")
			args = append(args, *filter.StartDate)
		}
		if filter.EndDate != nil {
			conditions = append(conditions, "created_at <= ?")
			args = append(args, *filter.EndDate)
		}

		query := `SELECT id, username, project_id, resource_id, action, is_success, summary, before, after, trace_id, created_at
	          FROM audit_logs`
		if len(conditions) > 0 {
			query += " WHERE " + strings.Join(conditions, " AND ")
		}
		query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
		args = append(args, pagination.PerPage, pagination.Offset())

		err := sqlx.SelectContext(context.Background(), a.pool, &entries, query, args...)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if entries == nil {
			entries = []security.AuditLogEntry{}
		}
		return entries, nil
	})
}
