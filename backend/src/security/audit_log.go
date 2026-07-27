package security

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type AuditLogAction string

const (
	ActionCreateProject  AuditLogAction = "project.create"
	ActionUpdateProject  AuditLogAction = "project.update"
	ActionDeleteProject  AuditLogAction = "project.delete"
	ActionCreateMonitor  AuditLogAction = "monitor.create"
	ActionUpdateMonitor  AuditLogAction = "monitor.update"
	ActionDeleteMonitor  AuditLogAction = "monitor.delete"
	ActionCreateUser     AuditLogAction = "user.create"
	ActionUpdateUser     AuditLogAction = "user.update"
	ActionDeleteUser     AuditLogAction = "user.delete"
	ActionLogin          AuditLogAction = "auth.login"
	ActionLogout         AuditLogAction = "auth.logout"
	ActionFailedLogin    AuditLogAction = "auth.failed_login"
	ActionPasswordChange AuditLogAction = "auth.password_change"
)

type AuditLogEntry struct {
	ID         uuid.UUID      `json:"id"                   db:"id"`
	Username   *string        `json:"username,omitempty"   db:"username"` /// "system" if the action was performed by the system (e.g. scheduled task)
	ProjectID  *uuid.UUID     `json:"projectId,omitempty"  db:"project_id"`
	ResourceID *uuid.UUID     `json:"resourceId,omitempty" db:"resource_id"` // ID of the resource that was acted upon, e.g. monitor ID, project ID, etc. Can be empty if not applicable.
	Action     AuditLogAction `json:"action"               db:"action"`
	IsSuccess  bool           `json:"isSuccess"            db:"is_success"`
	Summary    string         `json:"summary,omitempty"    db:"summary"`
	Before     *string        `json:"before,omitempty"     db:"before"`   // JSON string representing the state of the resource before the action. Can be empty if not applicable.
	After      *string        `json:"after,omitempty"      db:"after"`    // JSON string representing the state of the resource after the action. Can be empty if not applicable.
	TraceID    *string        `json:"traceId,omitempty"    db:"trace_id"` // Trace ID for correlating with other logs/traces. Can be empty if not applicable.
	CreatedAt  time.Time      `json:"createdAt"            db:"created_at"`
}

// AuditLogParams encapsulates the parameters needed to create a new audit log entry.
type AuditLogParams struct {
	Username   *string
	ProjectID  *uuid.UUID
	ResourceID *uuid.UUID
	Action     AuditLogAction
	IsSuccess  bool
	Summary    string
	Before     any
	After      any
}

// NewAuditLogEntry creates a new AuditLogEntry, handling JSON marshaling of before/after states,
// capturing the TraceID from context, and calling BeforeCreate.
func NewAuditLogEntry(
	ctx context.Context,
	params AuditLogParams,
) (AuditLogEntry, error) {
	var beforeStr *string
	if params.Before != nil {
		b, err := json.Marshal(params.Before)
		if err != nil {
			return AuditLogEntry{}, fmt.Errorf("failed to marshal before state: %w", err)
		}
		str := string(b)
		beforeStr = &str
	}

	var afterStr *string
	if params.After != nil {
		b, err := json.Marshal(params.After)
		if err != nil {
			return AuditLogEntry{}, fmt.Errorf("failed to marshal after state: %w", err)
		}
		str := string(b)
		afterStr = &str
	}

	traceID, _ := ctx.Value("X-Trace-Id").(string)

	entry := AuditLogEntry{
		Username:   params.Username,
		ProjectID:  params.ProjectID,
		ResourceID: params.ResourceID,
		Action:     params.Action,
		IsSuccess:  params.IsSuccess,
		Summary:    params.Summary,
		Before:     beforeStr,
		After:      afterStr,
		TraceID:    &traceID,
	}
	entry.BeforeCreate()

	return entry, nil
}

func (a *AuditLogEntry) BeforeCreate() {
	a.ID = uuid.New()
	a.CreatedAt = time.Now().UTC()
}

type AuditLogFilter struct {
	UserID     *string
	ProjectID  *uuid.UUID
	ResourceID *uuid.UUID
	Action     *AuditLogAction
	IsSuccess  *bool
	TraceID    *string
	StartDate  *time.Time
	EndDate    *time.Time
}

func (f *AuditLogFilter) ValidateForNonInstanceAdmin() error {
	// Allow filtering without project ID only for instance admins
	if f.ProjectID == nil {
		return errors.New("filtering without project ID is allowed only for instance admins")
	}
	return nil
}

func AuditLogFilterFromRequest(r *http.Request) (*AuditLogFilter, error) {
	f := &AuditLogFilter{}
	query := r.URL.Query()

	if userID := query.Get("userId"); userID != "" {
		f.UserID = &userID
	}

	if projectIDStr := query.Get("projectId"); projectIDStr != "" {
		projectID, err := uuid.Parse(projectIDStr)
		if err != nil {
			return nil, errors.New("invalid projectId format")
		}
		f.ProjectID = &projectID
	}

	if resourceIDStr := query.Get("resourceId"); resourceIDStr != "" {
		resourceID, err := uuid.Parse(resourceIDStr)
		if err != nil {
			return nil, errors.New("invalid resourceId format")
		}
		f.ResourceID = &resourceID
	}

	if action := query.Get("action"); action != "" {
		actionEnum := AuditLogAction(action)
		f.Action = &actionEnum
	}

	if isSuccessStr := query.Get("isSuccess"); isSuccessStr != "" {
		isSuccess := isSuccessStr == "true"
		f.IsSuccess = &isSuccess
	}

	if traceID := query.Get("traceId"); traceID != "" {
		f.TraceID = &traceID
	}

	if startDateStr := query.Get("startDate"); startDateStr != "" {
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			return nil, errors.New("invalid startDate format, expected RFC3339")
		}
		f.StartDate = &startDate
	}

	if endDateStr := query.Get("endDate"); endDateStr != "" {
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			return nil, errors.New("invalid endDate format, expected RFC3339")
		}
		f.EndDate = &endDate
	}

	return f, nil
}

func GetTraceIDFromContext(ctx context.Context) *string {
	traceID, ok := ctx.Value("requestId").(string)
	if !ok {
		return nil
	}
	return &traceID
}
