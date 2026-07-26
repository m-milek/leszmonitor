package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/security"
)

type IMonitorService interface {
	CreateMonitor(ctx context.Context, projectSlug string, username string, monitor monitors.Monitor) (*MonitorCreateResponse, *ServiceError)
	DeleteMonitor(ctx context.Context, username string, id string) *ServiceError
	GetMonitorsByProjectSlug(ctx context.Context, projectSlug string) ([]monitors.Monitor, *ServiceError)
	GetMonitorByID(ctx context.Context, id string) (*monitors.Monitor, *ServiceError)
	UpdateMonitor(ctx context.Context, username string, monitor monitors.Monitor) *ServiceError
	GetMonitorBySlugByProject(ctx context.Context, projectSlug string, slug string) (*monitors.Monitor, *ServiceError)
	UpdateMonitorStateByID(ctx context.Context, username string, monitorID uuid.UUID, state monitors.MonitorState) *ServiceError
}

// MonitorService handles monitor-related CRUD operations.
type MonitorService struct {
	db db.DB
}

type MonitorServiceDeps struct {
	DB db.DB
}

func NewMonitorService(deps MonitorServiceDeps) *MonitorService {
	return &MonitorService{
		db: deps.DB,
	}
}

type MonitorCreateResponse struct {
	MonitorID string `json:"monitorId"`
}

// CreateMonitor creates a new monitor in the specified project.
func (s *MonitorService) CreateMonitor(ctx context.Context, projectSlug string, username string, monitor monitors.Monitor) (*MonitorCreateResponse, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "CreateMonitor")
	logger.Trace().Str("projectSlug", projectSlug).Interface("monitor", monitor).Str("username", username).Msg("Creating monitor")

	project, err := s.db.Projects().GetProjectBySlug(ctx, projectSlug)
	if err != nil {
		logger.Error().Err(err).Str("projectSlug", projectSlug).Msg("Failed to find project by slug")
		return nil, NewNotFoundError("failed to find project: %w", err)
	}

	initializedMonitor := monitors.InitializeFromPayload(monitor, project.ID)

	if err := initializedMonitor.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid monitor configuration")
		return nil, NewBadRequestError("invalid monitor configuration: %w", err)
	}

	var monitorFromDB *monitors.Monitor
	var createErr error
	txErr := s.db.WithTx(ctx, func(tx db.DB) error {
		monitorFromDB, createErr = tx.Monitors().InsertMonitor(ctx, *initializedMonitor)
		if createErr != nil {
			return createErr
		}

		entry, err := security.NewAuditLogEntry(
			ctx,
			&username,
			&project.ID,
			&monitorFromDB.ID,
			security.ActionCreateMonitor,
			true,
			fmt.Sprintf("Monitor with ID %s created", monitorFromDB.ID),
			nil,
			monitorFromDB,
		)
		if err != nil {
			return NewInternalError(FormatFailedToCreateAuditLog, err)
		}

		_, auditErr := tx.AuditLog().InsertAuditLogEntry(ctx, entry)
		return auditErr
	})
	if txErr != nil {
		logger.Error().Err(txErr).Msg("Failed to create monitor within transaction")
		return nil, NewInternalError("failed to create monitor within transaction: %w", txErr)
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitorFromDB.ID,
		Status:  monitors.Created,
		Monitor: monitorFromDB,
	})

	logger.Debug().Str("id", monitor.ID.String()).Msg("Monitor created")
	return &MonitorCreateResponse{MonitorID: monitorFromDB.ID.String()}, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorService) DeleteMonitor(ctx context.Context, username string, id string) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "DeleteMonitor")
	logger.Trace().Str("id", id).Str("username", username).Msg("Deleting monitor")

	monitorUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error().Str("id", id).Msg("Invalid monitor ID format")
		return NewBadRequestError("invalid monitor ID format: %w", err)
	}

	monitorBeforeDelete, err := s.db.Monitors().GetMonitorByID(ctx, monitorUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Monitor not found in database")
			return NewNotFoundError("monitor with ID %s not found", id)
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor before deletion")
		return NewInternalError("failed to retrieve monitor before deletion: %w", err)
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

		entry, err := security.NewAuditLogEntry(
			ctx,
			&username,
			&monitorBeforeDelete.ProjectID,
			&monitorUUID,
			security.ActionDeleteMonitor,
			true,
			fmt.Sprintf("Monitor with ID %s deleted", monitorUUID.String()),
			monitorBeforeDelete,
			nil,
		)
		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("Failed to create audit log entry for monitor deletion")
			return NewInternalError(FormatFailedToCreateAuditLog, err)
		}

		_, err = tx.AuditLog().InsertAuditLogEntry(ctx, entry)

		if err != nil {
			logger.Error().Err(err).Str("id", id).Msg("Failed to insert audit log entry for monitor deletion")
			return NewInternalError("failed to insert audit log entry for monitor deletion: %w", err)
		}

		return nil
	}); txErr != nil {
		if errors.Is(txErr, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Monitor not found or already deleted")
			return NewNotFoundError("monitor not found or already deleted")
		}
		logger.Error().Err(txErr).Str("id", id).Msg("Failed to delete monitor")
		return NewInternalError("failed to delete monitor: %w", txErr)
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      *deletedID,
		Status:  monitors.Deleted,
		Monitor: nil,
	})

	logger.Debug().Str(id, id).Msg("Monitor deleted")
	return nil
}

// GetMonitorsByProjectSlug retrieves all monitors for the project.
func (s *MonitorService) GetMonitorsByProjectSlug(ctx context.Context, projectSlug string) ([]monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "GetMonitorsByProjectSlug")
	logger.Trace().Str("projectSlug", projectSlug).Msg("Retrieving monitors by project slug")

	project, err := s.db.Projects().GetProjectBySlug(ctx, projectSlug)
	if err != nil {
		logger.Error().Err(err).Str("projectSlug", projectSlug).Msg("Failed to find project by slug")
		return nil, NewInternalError("failed to find project: %w", err)
	}

	monitorsList, err := s.db.Monitors().GetMonitorsByProjectID(ctx, project.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return nil, NewInternalError("failed to retrieve monitors: %w", err)
	}

	logger.Debug().Str("projectSlug", projectSlug).Int("monitorCount", len(monitorsList)).Msg("Monitors retrieved successfully")
	return monitorsList, nil
}

// GetMonitorByID retrieves a specific monitor by its ID.
func (s *MonitorService) GetMonitorByID(ctx context.Context, id string) (*monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "GetMonitorByID")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by ID")

	monitorUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error().Str("id", id).Msg("Invalid monitor ID format (should be uuid)")
		return nil, NewBadRequestError("invalid monitor ID format: %w", err)
	}

	monitor, err := s.db.Monitors().GetMonitorByID(ctx, monitorUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Monitor not found in database")
			return nil, NewNotFoundError("monitor with id %s not found", id)
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor from database")
		return nil, NewInternalError("failed to retrieve monitor: %w", err)
	}

	logger.Debug().Str("id", id).Interface("monitor", monitor).Msg("Monitor retrieved successfully")
	return monitor, nil
}

// UpdateMonitor updates an existing monitor's configuration.
func (s *MonitorService) UpdateMonitor(ctx context.Context, username string, monitor monitors.Monitor) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	txErr := s.db.WithTx(ctx, func(tx db.DB) error {
		existingMonitor, err := tx.Monitors().GetMonitorByID(ctx, monitor.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitor.ID.String()).Msg("Monitor not found")
				return NewNotFoundError("monitor with ID %s not found", monitor.ID)
			}
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to retrieve existing monitor for update")
			return fmt.Errorf("failed to retrieve existing monitor for update: %w", err)
		}

		monitor.State = existingMonitor.State
		monitor.ProjectID = existingMonitor.ProjectID

		if err := monitor.Validate(); err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Interface("monitor", monitor).Msg("Invalid monitor configuration")
			return NewBadRequestError("invalid monitor configuration: %w", err)
		}

		_, err = tx.Monitors().UpdateMonitor(ctx, monitor)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to update monitor in database")
			return fmt.Errorf("failed to update monitor in database: %w", err)
		}

		entry, err := security.NewAuditLogEntry(
			ctx,
			&username,
			&existingMonitor.ProjectID,
			&monitor.ID,
			security.ActionUpdateMonitor,
			true,
			fmt.Sprintf("Monitor with ID %s updated", monitor.ID),
			existingMonitor,
			monitor,
		)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to create audit log entry for monitor update")
			return NewInternalError(FormatFailedToCreateAuditLog, err)
		}

		_, auditErr := tx.AuditLog().InsertAuditLogEntry(ctx, entry)

		if auditErr != nil {
			logger.Error().Err(auditErr).Str("id", monitor.ID.String()).Msg("Failed to insert audit log entry for monitor update")
			return NewInternalError("failed to insert audit log entry for monitor update: %w", auditErr)
		}

		return nil
	})
	if txErr != nil {
		if serviceErr, ok := txErr.(*ServiceError); ok {
			return serviceErr
		}
		logger.Error().Err(txErr).Str("id", monitor.ID.String()).Msg("Failed to update monitor within transaction")
		return NewInternalError("failed to update monitor within transaction: %w", txErr)
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitor.ID,
		Status:  monitors.Edited,
		Monitor: nil,
	})

	logger.Debug().Str("id", monitor.ID.String()).Msg("Monitor updated")
	return nil
}

func (s *MonitorService) GetMonitorBySlugByProject(ctx context.Context, projectSlug string, slug string) (*monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "GetMonitorBySlugByProject")
	logger.Trace().Str("slug", slug).Msg("Retrieving monitor by slug and project")

	project, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		logger.Error().Err(getErr).Str("slug", slug).Str("projectSlug", projectSlug).Msg("Failed to retrieve project by slug")
		return nil, getErr
	}

	monitor, err := s.db.Monitors().GetMonitorBySlugByProject(ctx, slug, project.ID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("slug", slug).Str("projectSlug", projectSlug).Msg("Monitor not found in project")
			return nil, NewNotFoundError("monitor with slug %s not found in project", slug)
		}
		logger.Error().Err(err).Str("slug", slug).Str("projectSlug", projectSlug).Msg("Failed to retrieve monitor by slug and project")
		return nil, NewInternalError("failed to retrieve monitor by slug and project: %w", err)
	}

	logger.Debug().Str("slug", slug).Msg("Monitor retrieved by slug and project")
	return monitor, nil
}

func (s *MonitorService) UpdateMonitorStateByID(ctx context.Context, username string, monitorID uuid.UUID, state monitors.MonitorState) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitorStateByID")
	logger.Trace().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Updating monitor state by ID")

	if !monitors.IsValidMonitorState(string(state)) {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Invalid monitor state provided")
		return NewBadRequestError("invalid monitor state: %s", state)
	}

	monitor, err := s.db.Monitors().GetMonitorByID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", monitorID.String()).Msg("Monitor not found for state update")
			return NewNotFoundError("monitor with ID %s not found", monitorID.String())
		}
		logger.Error().Err(err).Str("id", monitorID.String()).Msg("Failed to retrieve monitor for state update")
		return NewInternalError("failed to retrieve monitor for state update: %w", err)
	}

	if monitor.State == state {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Nothing to update in monitor, no update needed")
		return nil
	}

	monitor.State = state

	_, updateErr := s.db.Monitors().UpdateMonitor(ctx, *monitor)
	if updateErr != nil {
		logger.Error().Err(updateErr).Str("id", monitorID.String()).Str("newState", string(state)).Msg("Failed to update monitor state in database")
		return NewInternalError("failed to update monitor state in database: %w", updateErr)
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitor.ID,
		Monitor: monitor,
		Status:  monitors.Edited,
	})

	logger.Debug().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Monitor state updated")
	return nil
}
