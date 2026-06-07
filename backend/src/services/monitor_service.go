package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/security"
)

type IMonitorService interface {
	CreateMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitor monitors.Monitor) (*MonitorCreateResponse, *ServiceError)
	DeleteMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string) *ServiceError
	GetMonitorsByProjectID(ctx context.Context, projectAuth *authorization.ProjectAuthorization) ([]monitors.Monitor, *ServiceError)
	GetMonitorByID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string) (*monitors.Monitor, *ServiceError)
	UpdateMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitor monitors.Monitor) *ServiceError
	GetMonitorBySlugByProject(ctx context.Context, auth *authorization.ProjectAuthorization, slug string) (*monitors.Monitor, *ServiceError)
	UpdateMonitorStateByID(ctx context.Context, auth *authorization.ProjectAuthorization, monitorID uuid.UUID, state monitors.MonitorState) *ServiceError
}

// MonitorService handles monitor-related CRUD operations.
type MonitorService struct {
	db   db.DB
	auth IAuthorizer
}

type MonitorServiceDeps struct {
	DB   db.DB
	Auth IAuthorizer
}

func NewMonitorService(deps MonitorServiceDeps) *MonitorService {
	return &MonitorService{
		db:   deps.DB,
		auth: deps.Auth,
	}
}

type MonitorCreateResponse struct {
	MonitorID string `json:"monitorId"`
}

// CreateMonitor creates a new monitor in the specified project.
func (s *MonitorService) CreateMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitor monitors.Monitor) (*MonitorCreateResponse, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "CreateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Creating new monitor")

	project, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorEditor)
	if authErr != nil {
		return nil, authErr
	}

	initializedMonitor := monitors.InitializeFromPayload(monitor, project.ID)

	if err := initializedMonitor.Validate(); err != nil {
		logger.Warn().Err(err).Msg("Invalid monitor configuration")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration: %w", err),
		}
	}

	var monitorFromDB *monitors.Monitor
	var createErr error
	txErr := s.db.WithTx(ctx, func(tx db.DB) error {
		monitorFromDB, createErr = tx.Monitors().InsertMonitor(ctx, *initializedMonitor)
		if createErr != nil {
			return createErr
		}

		monitorJSON, err := json.Marshal(monitorFromDB)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to marshal monitor for audit log")
			return fmt.Errorf("failed to marshal monitor for audit log: %w", err)
		}

		entry := security.AuditLogEntry{
			Username:   &projectAuth.Username,
			ProjectID:  &projectAuth.ProjectID,
			ResourceID: &monitorFromDB.ID,
			Action:     security.ActionCreateMonitor,
			IsSuccess:  true,
			Before:     nil,
			After:      new(string(monitorJSON)),
			Summary:    fmt.Sprintf("Monitor with ID %s created", monitorFromDB.ID),
			TraceID:    security.GetTraceIDFromContext(ctx),
		}
		entry.BeforeCreate()

		_, auditErr := tx.AuditLog().InsertAuditLogEntry(ctx, entry)
		return auditErr
	})
	if txErr != nil {
		logger.Error().Err(txErr).Msg("Failed to create monitor within transaction")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to create monitor within transaction: %w", txErr),
		}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitorFromDB.ID,
		Status:  monitors.Created,
		Monitor: monitorFromDB,
	})

	logger.Info().Str("id", monitor.ID.String()).Msg("Monitor created")
	return &MonitorCreateResponse{MonitorID: monitorFromDB.ID.String()}, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorService) DeleteMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string) *ServiceError {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	_, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorAdmin)
	if authErr != nil {
		return authErr
	}

	monitorUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Warn().Str("id", id).Msg("Invalid monitor ID format")
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid monitor ID format: %w", err)}
	}

	monitorBeforeDelete, err := s.db.Monitors().GetMonitorByID(ctx, monitorUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("id", id).Msg("Monitor not found in database")
			return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with ID %s not found", id)}
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor before deletion")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve monitor before deletion: %w", err)}
	}

	monitorBeforeDeleteJSON, err := json.Marshal(monitorBeforeDelete)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to marshal monitor state before deletion for audit log")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to marshal monitor state before deletion for audit log: %w", err)}
	}

	var deletedID *uuid.UUID
	if txErr := s.db.WithTx(ctx, func(tx db.DB) error {
		var err error
		deletedID, err = tx.Monitors().DeleteMonitorByID(ctx, monitorUUID)
		if err != nil {
			return err
		}
		if deletedID == nil {
			return db.ErrNotFound
		}

		entry := security.AuditLogEntry{
			Username:   &projectAuth.Username,
			ProjectID:  &projectAuth.ProjectID,
			ResourceID: &monitorUUID,
			Action:     security.ActionDeleteMonitor,
			IsSuccess:  true,
			Before:     new(string(monitorBeforeDeleteJSON)),
			After:      nil,
			Summary:    fmt.Sprintf("Monitor with ID %s deleted", monitorUUID.String()),
			TraceID:    security.GetTraceIDFromContext(ctx),
		}
		entry.BeforeCreate()

		_, err = tx.AuditLog().InsertAuditLogEntry(ctx, entry)
		return err
	}); txErr != nil {
		if errors.Is(txErr, db.ErrNotFound) {
			return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor not found or already deleted")}
		}
		logger.Error().Err(txErr).Str("id", id).Msg("Failed to delete monitor")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor: %w", txErr),
		}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      *deletedID,
		Status:  monitors.Deleted,
		Monitor: nil,
	})

	logger.Info().Str(id, id).Msg("Monitor deleted")
	return nil
}

// GetMonitorsByProjectID retrieves all monitors for the project in the provided ProjectAuth.
func (s *MonitorService) GetMonitorsByProjectID(ctx context.Context, projectAuth *authorization.ProjectAuthorization) ([]monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "GetMonitorsByProjectID")
	logger.Trace().Msg("Retrieving monitors for project")

	project, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitorsList, err := s.db.Monitors().GetMonitorsByProjectID(ctx, project.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitors: %w", err),
		}
	}

	return monitorsList, nil
}

// GetMonitorByID retrieves a specific monitor by its slug.
func (s *MonitorService) GetMonitorByID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string) (*monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "GetMonitorByID")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by slug")

	_, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	// TODO fix - monitor ID uses slug - bug. This method is not used by the frontend
	monitor, err := s.db.Monitors().GetMonitorBySlug(ctx, id, projectAuth.ProjectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("id", id).Msg("Monitor not found in database")
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with slug %s not found", id)}
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor from database")
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve monitor: %w", err)}
	}

	if monitor.ProjectID != projectAuth.ProjectID {
		logger.Warn().Str("id", id).Msg("Monitor does not belong to the authorized project")
		return nil, &ServiceError{Code: http.StatusForbidden, Err: fmt.Errorf("monitor with slug %s does not belong to the authorized project", id)}
	}

	return monitor, nil
}

// UpdateMonitor updates an existing monitor's configuration.
func (s *MonitorService) UpdateMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitor monitors.Monitor) *ServiceError {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	_, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorEditor)
	if authErr != nil {
		return authErr
	}

	txErr := s.db.WithTx(ctx, func(tx db.DB) error {
		existingMonitor, err := tx.Monitors().GetMonitorByID(ctx, monitor.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with ID %s not found", monitor.ID)}
			}
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to retrieve existing monitor for update")
			return fmt.Errorf("failed to retrieve existing monitor for update: %w", err)
		}

		if existingMonitor.ProjectID != projectAuth.ProjectID {
			logger.Warn().Str("id", monitor.ID.String()).Msg("Monitor does not belong to the authorized project")
			return &ServiceError{Code: http.StatusForbidden, Err: fmt.Errorf("monitor with ID %s does not belong to the authorized project", monitor.ID)}
		}

		monitor.State = existingMonitor.State
		monitor.ProjectID = existingMonitor.ProjectID

		if err := monitor.Validate(); err != nil {
			return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid monitor configuration: %w", err)}
		}

		_, err = tx.Monitors().UpdateMonitor(ctx, monitor)
		if err != nil {
			return fmt.Errorf("failed to update monitor in database: %w", err)
		}

		beforeJSON, err := json.Marshal(existingMonitor)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to marshal existing monitor state for audit log")
			return fmt.Errorf("failed to marshal existing monitor state for audit log: %w", err)
		}
		afterJSON, err := json.Marshal(monitor)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to marshal updated monitor state for audit log")
			return fmt.Errorf("failed to marshal updated monitor state for audit log: %w", err)
		}

		entry := security.AuditLogEntry{
			Username:   &projectAuth.Username,
			ProjectID:  &projectAuth.ProjectID,
			ResourceID: &monitor.ID,
			Action:     security.ActionUpdateMonitor,
			IsSuccess:  true,
			Before:     new(string(beforeJSON)),
			After:      new(string(afterJSON)),
			Summary:    fmt.Sprintf("Monitor with ID %s updated", monitor.ID),
			TraceID:    security.GetTraceIDFromContext(ctx),
		}
		entry.BeforeCreate()

		_, auditErr := tx.AuditLog().InsertAuditLogEntry(ctx, entry)

		return auditErr
	})
	if txErr != nil {
		if serviceErr, ok := txErr.(*ServiceError); ok {
			return serviceErr
		}
		logger.Error().Err(txErr).Str("id", monitor.ID.String()).Msg("Failed to update monitor within transaction")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to update monitor within transaction: %w", txErr)}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitor.ID,
		Status:  monitors.Edited,
		Monitor: nil,
	})

	logger.Info().Str("id", monitor.ID.String()).Msg("Monitor updated")
	return nil
}

func (s *MonitorService) GetMonitorBySlugByProject(ctx context.Context, auth *authorization.ProjectAuthorization, slug string) (*monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "GetMonitorBySlugByProject")
	logger.Trace().Str("slug", slug).Str("projectID", auth.ProjectID.String()).Msg("Retrieving monitor by slug and project")

	_, authErr := s.auth.authorizeProjectAction(ctx, auth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitor, err := s.db.Monitors().GetMonitorBySlugByProject(ctx, slug, auth.ProjectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with slug %s not found in project", slug)}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve monitor by slug and project: %w", err)}
	}

	logger.Debug().Str("slug", slug).Str("projectID", auth.ProjectID.String()).Msg("Monitor retrieved by slug and project")
	return monitor, nil
}

func (s *MonitorService) UpdateMonitorStateByID(ctx context.Context, auth *authorization.ProjectAuthorization, monitorID uuid.UUID, state monitors.MonitorState) *ServiceError {
	logger := MethodLoggerFromContext(ctx, "MonitorService", "UpdateMonitorStateByID")
	logger.Trace().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Updating monitor state by ID")

	if !monitors.IsValidMonitorState(string(state)) {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Invalid monitor state provided")
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid monitor state: %s", state)}
	}

	_, authErr := s.auth.authorizeProjectAction(ctx, auth, models.PermissionMonitorEditor)
	if authErr != nil {
		return authErr
	}

	monitor, err := s.db.Monitors().GetMonitorByID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with ID %s not found", monitorID.String())}
		}
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve monitor for state update: %w", err)}
	}

	if monitor.State == state {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Monitor state is already in the desired state, no update needed")
		return nil
	}

	monitor.State = state

	_, updateErr := s.db.Monitors().UpdateMonitor(ctx, *monitor)
	if updateErr != nil {
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to update monitor state in database: %w", updateErr)}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitor.ID,
		Monitor: monitor,
		Status:  monitors.Edited,
	})

	logger.Info().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Monitor state updated")
	return nil
}
