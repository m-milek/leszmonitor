package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/security"
)

type IMonitorService interface {
	CreateMonitor(
		ctx context.Context,
		projectSlug string,
		monitor monitors.Monitor,
	) (*MonitorCreateResponse, *ServiceError)
	DeleteMonitor(ctx context.Context, id string) *ServiceError
	GetMonitorsByProjectSlug(ctx context.Context, projectSlug string) ([]monitors.Monitor, *ServiceError)
	GetMonitorByID(ctx context.Context, id string) (*monitors.Monitor, *ServiceError)
	UpdateMonitor(ctx context.Context, monitor monitors.Monitor) *ServiceError
	GetMonitorBySlugByProject(ctx context.Context, projectSlug string, slug string) (*monitors.Monitor, *ServiceError)
	UpdateMonitorStateByID(
		ctx context.Context,
		monitorID uuid.UUID,
		state monitors.MonitorState,
	) *ServiceError
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
func (s *MonitorService) CreateMonitor(
	ctx context.Context,
	projectSlug string,
	monitor monitors.Monitor,
) (*MonitorCreateResponse, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "CreateMonitor")
	logger.Trace().
		Str("projectSlug", projectSlug).
		Interface("monitor", monitor).
		Msg("Creating monitor")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, NewUnauthorizedError("user claims not found in context")
	}

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

	monitorFromDB, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (*monitors.Monitor, *security.AuditLogParams, error) {
		m, createErr := tx.Monitors().InsertMonitor(ctx, *initializedMonitor)
		if createErr != nil {
			return nil, nil, createErr
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ProjectID:  &project.ID,
			ResourceID: &m.ID,
			Action:     security.ActionCreateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s created", m.ID),
			After:      m,
		}
		return m, params, nil
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

	logger.Debug().Str("id", monitorFromDB.ID.String()).Msg("Monitor created")
	return &MonitorCreateResponse{MonitorID: monitorFromDB.ID.String()}, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorService) DeleteMonitor(ctx context.Context, id string) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

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

	deletedID, txErr := db.WithAuditedTx(ctx, s.db, func(tx db.DB) (*uuid.UUID, *security.AuditLogParams, error) {
		delID, err := tx.Monitors().DeleteMonitorByID(ctx, monitorUUID)
		if err != nil {
			return nil, nil, err
		}
		if delID == nil {
			return nil, nil, db.ErrNotFound
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ProjectID:  &monitorBeforeDelete.ProjectID,
			ResourceID: &monitorUUID,
			Action:     security.ActionDeleteMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s deleted", monitorUUID.String()),
			Before:     monitorBeforeDelete,
		}

		return delID, params, nil
	})
	if txErr != nil {
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

	logger.Debug().Str("id", id).Msg("Monitor deleted")
	return nil
}

// GetMonitorsByProjectSlug retrieves all monitors for the project.
func (s *MonitorService) GetMonitorsByProjectSlug(
	ctx context.Context,
	projectSlug string,
) ([]monitors.Monitor, *ServiceError) {
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

	logger.Debug().
		Str("projectSlug", projectSlug).
		Int("monitorCount", len(monitorsList)).
		Msg("Monitors retrieved successfully")
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
func (s *MonitorService) UpdateMonitor(ctx context.Context, monitor monitors.Monitor) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	txErr := db.WithAuditedVoidTx(ctx, s.db, func(tx db.DB) (*security.AuditLogParams, error) {
		existingMonitor, err := tx.Monitors().GetMonitorByID(ctx, monitor.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitor.ID.String()).Msg("Monitor not found")
				return nil, NewNotFoundError("monitor with ID %s not found", monitor.ID)
			}
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to retrieve existing monitor for update")
			return nil, fmt.Errorf("failed to retrieve existing monitor for update: %w", err)
		}

		monitor.State = existingMonitor.State
		monitor.ProjectID = existingMonitor.ProjectID

		if err := monitor.Validate(); err != nil {
			logger.Error().
				Err(err).
				Str("id", monitor.ID.String()).
				Interface("monitor", monitor).
				Msg("Invalid monitor configuration")
			return nil, NewBadRequestError("invalid monitor configuration: %w", err)
		}

		_, err = tx.Monitors().UpdateMonitor(ctx, monitor)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to update monitor in database")
			return nil, fmt.Errorf("failed to update monitor in database: %w", err)
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ProjectID:  &existingMonitor.ProjectID,
			ResourceID: &monitor.ID,
			Action:     security.ActionUpdateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s updated", monitor.ID),
			Before:     existingMonitor,
			After:      monitor,
		}

		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*ServiceError](txErr); ok2 {
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

func (s *MonitorService) GetMonitorBySlugByProject(
	ctx context.Context,
	projectSlug string,
	slug string,
) (*monitors.Monitor, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "GetMonitorBySlugByProject")
	logger.Trace().Str("slug", slug).Msg("Retrieving monitor by slug and project")

	project, getErr := internalGetProjectBySlug(ctx, s.db, projectSlug)
	if getErr != nil {
		logger.Error().
			Err(getErr).
			Str("slug", slug).
			Str("projectSlug", projectSlug).
			Msg("Failed to retrieve project by slug")
		return nil, getErr
	}

	monitor, err := s.db.Monitors().GetMonitorBySlugByProject(ctx, slug, project.ID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("slug", slug).Str("projectSlug", projectSlug).Msg("Monitor not found in project")
			return nil, NewNotFoundError("monitor with slug %s not found in project", slug)
		}
		logger.Error().
			Err(err).
			Str("slug", slug).
			Str("projectSlug", projectSlug).
			Msg("Failed to retrieve monitor by slug and project")
		return nil, NewInternalError("failed to retrieve monitor by slug and project: %w", err)
	}

	logger.Debug().Str("slug", slug).Msg("Monitor retrieved by slug and project")
	return monitor, nil
}

func (s *MonitorService) UpdateMonitorStateByID(
	ctx context.Context,
	monitorID uuid.UUID,
	state monitors.MonitorState,
) *ServiceError {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitorStateByID")
	logger.Trace().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Updating monitor state by ID")

	userClaims, ok := authorization.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return NewUnauthorizedError("user claims not found in context")
	}

	if !monitors.IsValidMonitorState(string(state)) {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Invalid monitor state provided")
		return NewBadRequestError("invalid monitor state: %s", state)
	}

	txErr := db.WithAuditedVoidTx(ctx, s.db, func(tx db.DB) (*security.AuditLogParams, error) {
		monitor, err := tx.Monitors().GetMonitorByID(ctx, monitorID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitorID.String()).Msg("Monitor not found for state update")
				return nil, NewNotFoundError("monitor with ID %s not found", monitorID.String())
			}
			logger.Error().Err(err).Str("id", monitorID.String()).Msg("Failed to retrieve monitor for state update")
			return nil, NewInternalError("failed to retrieve monitor for state update: %w", err)
		}

		if monitor.State == state {
			logger.Warn().
				Str("id", monitorID.String()).
				Str("state", string(state)).
				Msg("Nothing to update in monitor, no update needed")
			return nil, nil
		}

		oldMonitor := *monitor
		monitor.State = state

		_, updateErr := tx.Monitors().UpdateMonitor(ctx, *monitor)
		if updateErr != nil {
			logger.Error().
				Err(updateErr).
				Str("id", monitorID.String()).
				Str("newState", string(state)).
				Msg("Failed to update monitor state in database")
			return nil, NewInternalError("failed to update monitor state in database: %w", updateErr)
		}

		params := &security.AuditLogParams{
			Username:   &userClaims.Username,
			ProjectID:  &monitor.ProjectID,
			ResourceID: &monitor.ID,
			Action:     security.ActionUpdateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s state updated to %s", monitor.ID, state),
			Before:     oldMonitor,
			After:      *monitor,
		}

		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*ServiceError](txErr); ok2 {
			return serviceErr
		}
		return NewInternalError("failed to update monitor state: %w", txErr)
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:     monitorID,
		Status: monitors.Edited,
	})

	logger.Debug().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Monitor state updated")
	return nil
}
