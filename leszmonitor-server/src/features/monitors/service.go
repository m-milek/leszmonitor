package monitors

import (
	"context"
	"errors"
	"fmt"

	"github.com/m-milek/leszmonitor/features/users"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/audit"
	"github.com/m-milek/leszmonitor/platform/auth"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

type IMonitorService interface {
	CreateMonitor(
		ctx context.Context,
		monitor Monitor,
	) (*MonitorCreateResponse, *apperr.ServiceError)
	DeleteMonitor(ctx context.Context, id string) *apperr.ServiceError
	GetAllMonitors(ctx context.Context) ([]Monitor, *apperr.ServiceError)
	GetMonitorByID(ctx context.Context, id string) (*Monitor, *apperr.ServiceError)
	UpdateMonitor(ctx context.Context, monitor Monitor) *apperr.ServiceError
	UpdateMonitorStateByID(
		ctx context.Context,
		monitorID uuid.UUID,
		state MonitorRunState,
	) *apperr.ServiceError
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

// CreateMonitor creates a new monitor.
func (s *MonitorService) CreateMonitor(
	ctx context.Context,
	monitor Monitor,
) (*MonitorCreateResponse, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "CreateMonitor")
	logger.Trace().
		Interface("monitor", monitor).
		Msg("Creating monitor")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewUnauthorizedError("user claims not found in context")
	}

	owner, err := users.NewUserDAO(s.db.Querier()).GetUserByUsername(ctx, userClaims.Username)
	if err != nil {
		logger.Error().Err(err).Str("username", userClaims.Username).Msg("Failed to find creating user")
		return nil, apperr.NewInternalError("failed to find creating user: %w", err)
	}

	initializedMonitor := InitializeFromPayload(monitor, owner.ID)

	if err := initializedMonitor.Validate(); err != nil {
		logger.Error().Err(err).Msg("Invalid monitor configuration")
		return nil, apperr.NewBadRequestError("invalid monitor configuration: %w", err)
	}

	monitorFromDB, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*Monitor, *audit.AuditLogParams, error) {
		m, createErr := NewMonitorDAO(q).InsertMonitor(ctx, *initializedMonitor)
		if createErr != nil {
			return nil, nil, createErr
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &m.ID,
			Action:     audit.ActionCreateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s created", m.ID),
			After:      m,
		}
		return m, params, nil
	})
	if txErr != nil {
		logger.Error().Err(txErr).Msg("Failed to create monitor within transaction")
		return nil, apperr.NewInternalError("failed to create monitor within transaction: %w", txErr)
	}

	MonitorLifecycleChannel.Broadcast(MonitorLifecycleMessage{
		ID:      monitorFromDB.ID,
		Status:  Created,
		Monitor: monitorFromDB,
	})

	logger.Debug().Str("id", monitorFromDB.ID.String()).Msg("Monitor created")
	return &MonitorCreateResponse{MonitorID: monitorFromDB.ID.String()}, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorService) DeleteMonitor(ctx context.Context, id string) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewUnauthorizedError("user claims not found in context")
	}

	monitorUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Error().Str("id", id).Msg("Invalid monitor ID format")
		return apperr.NewBadRequestError("invalid monitor ID format: %w", err)
	}

	monitorBeforeDelete, err := NewMonitorDAO(s.db.Querier()).GetMonitorByID(ctx, monitorUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Monitor not found in database")
			return apperr.NewNotFoundError("monitor with ID %s not found", id)
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor before deletion")
		return apperr.NewInternalError("failed to retrieve monitor before deletion: %w", err)
	}

	deletedID, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*uuid.UUID, *audit.AuditLogParams, error) {
		delID, err := NewMonitorDAO(q).DeleteMonitorByID(ctx, monitorUUID)
		if err != nil {
			return nil, nil, err
		}
		if delID == nil {
			return nil, nil, db.ErrNotFound
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &monitorUUID,
			Action:     audit.ActionDeleteMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s deleted", monitorUUID.String()),
			Before:     monitorBeforeDelete,
		}

		return delID, params, nil
	})
	if txErr != nil {
		if errors.Is(txErr, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Monitor not found or already deleted")
			return apperr.NewNotFoundError("monitor not found or already deleted")
		}
		logger.Error().Err(txErr).Str("id", id).Msg("Failed to delete monitor")
		return apperr.NewInternalError("failed to delete monitor: %w", txErr)
	}

	MonitorLifecycleChannel.Broadcast(MonitorLifecycleMessage{
		ID:      *deletedID,
		Status:  Deleted,
		Monitor: nil,
	})

	logger.Debug().Str("id", id).Msg("Monitor deleted")
	return nil
}

// GetAllMonitors retrieves every monitor in the instance.
func (s *MonitorService) GetAllMonitors(ctx context.Context) ([]Monitor, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "GetAllMonitors")
	logger.Trace().Msg("Retrieving all monitors")

	allMonitors, err := NewMonitorDAO(s.db.Querier()).GetAllMonitors(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return nil, apperr.NewInternalError("failed to retrieve monitors: %w", err)
	}

	logger.Debug().Int("count", len(allMonitors)).Msg("Monitors retrieved successfully")
	return allMonitors, nil
}

// GetMonitorByID retrieves a specific monitor either by its UUID or, if id does not parse as a
// UUID, by its slug.
func (s *MonitorService) GetMonitorByID(ctx context.Context, id string) (*Monitor, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "GetMonitorByID")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by ID or slug")

	var (
		monitor *Monitor
		err     error
	)

	if monitorUUID, parseErr := uuid.Parse(id); parseErr == nil {
		monitor, err = NewMonitorDAO(s.db.Querier()).GetMonitorByID(ctx, monitorUUID)
	} else {
		monitor, err = NewMonitorDAO(s.db.Querier()).GetMonitorBySlug(ctx, id)
	}

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("id", id).Msg("Monitor not found in database")
			return nil, apperr.NewNotFoundError("monitor with id %s not found", id)
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor from database")
		return nil, apperr.NewInternalError("failed to retrieve monitor: %w", err)
	}

	logger.Debug().Str("id", id).Interface("monitor", monitor).Msg("Monitor retrieved successfully")
	return monitor, nil
}

// UpdateMonitor updates an existing monitor's configuration.
func (s *MonitorService) UpdateMonitor(ctx context.Context, monitor Monitor) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewUnauthorizedError("user claims not found in context")
	}

	txErr := audit.WithAuditedVoidTx(ctx, s.db, func(q db.Querier) (*audit.AuditLogParams, error) {
		existingMonitor, err := NewMonitorDAO(q).GetMonitorByID(ctx, monitor.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitor.ID.String()).Msg("Monitor not found")
				return nil, apperr.NewNotFoundError("monitor with ID %s not found", monitor.ID)
			}
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to retrieve existing monitor for update")
			return nil, fmt.Errorf("failed to retrieve existing monitor for update: %w", err)
		}

		monitor.RunState = existingMonitor.RunState
		monitor.OwnerID = existingMonitor.OwnerID

		if err := monitor.Validate(); err != nil {
			logger.Error().
				Err(err).
				Str("id", monitor.ID.String()).
				Interface("monitor", monitor).
				Msg("Invalid monitor configuration")
			return nil, apperr.NewBadRequestError("invalid monitor configuration: %w", err)
		}

		_, err = NewMonitorDAO(q).UpdateMonitor(ctx, monitor)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to update monitor in database")
			return nil, fmt.Errorf("failed to update monitor in database: %w", err)
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &monitor.ID,
			Action:     audit.ActionUpdateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s updated", monitor.ID),
			Before:     existingMonitor,
			After:      monitor,
		}

		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return serviceErr
		}
		logger.Error().Err(txErr).Str("id", monitor.ID.String()).Msg("Failed to update monitor within transaction")
		return apperr.NewInternalError("failed to update monitor within transaction: %w", txErr)
	}

	MonitorLifecycleChannel.Broadcast(MonitorLifecycleMessage{
		ID:      monitor.ID,
		Status:  Edited,
		Monitor: nil,
	})

	logger.Debug().Str("id", monitor.ID.String()).Msg("Monitor updated")
	return nil
}

func (s *MonitorService) UpdateMonitorStateByID(
	ctx context.Context,
	monitorID uuid.UUID,
	state MonitorRunState,
) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitorStateByID")
	logger.Trace().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Updating monitor state by ID")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewUnauthorizedError("user claims not found in context")
	}

	if !IsValidMonitorState(string(state)) {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Invalid monitor state provided")
		return apperr.NewBadRequestError("invalid monitor state: %s", state)
	}

	txErr := audit.WithAuditedVoidTx(ctx, s.db, func(q db.Querier) (*audit.AuditLogParams, error) {
		monitor, err := NewMonitorDAO(q).GetMonitorByID(ctx, monitorID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitorID.String()).Msg("Monitor not found for state update")
				return nil, apperr.NewNotFoundError("monitor with ID %s not found", monitorID.String())
			}
			logger.Error().Err(err).Str("id", monitorID.String()).Msg("Failed to retrieve monitor for state update")
			return nil, apperr.NewInternalError("failed to retrieve monitor for state update: %w", err)
		}

		if monitor.RunState == state {
			logger.Warn().
				Str("id", monitorID.String()).
				Str("state", string(state)).
				Msg("Nothing to update in monitor, no update needed")
			return nil, nil
		}

		oldMonitor := *monitor
		monitor.RunState = state

		_, updateErr := NewMonitorDAO(q).UpdateMonitor(ctx, *monitor)
		if updateErr != nil {
			logger.Error().
				Err(updateErr).
				Str("id", monitorID.String()).
				Str("newState", string(state)).
				Msg("Failed to update monitor state in database")
			return nil, apperr.NewInternalError("failed to update monitor state in database: %w", updateErr)
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &monitor.ID,
			Action:     audit.ActionUpdateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s state updated to %s", monitor.ID, state),
			Before:     oldMonitor,
			After:      *monitor,
		}

		return params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return serviceErr
		}
		return apperr.NewInternalError("failed to update monitor state: %w", txErr)
	}

	MonitorLifecycleChannel.Broadcast(MonitorLifecycleMessage{
		ID:     monitorID,
		Status: Edited,
	})

	logger.Debug().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Monitor state updated")
	return nil
}
