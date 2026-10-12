package monitors

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/probe"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/tags"
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
	) (*Monitor, *apperr.ServiceError)
	DeleteMonitor(ctx context.Context, id string) *apperr.ServiceError
	GetAllMonitors(ctx context.Context) ([]Monitor, *apperr.ServiceError)
	GetMonitorByID(ctx context.Context, id string) (*Monitor, *apperr.ServiceError)
	UpdateMonitor(ctx context.Context, monitor Monitor) (*Monitor, *apperr.ServiceError)
	UpdateMonitorStateByID(
		ctx context.Context,
		monitorID uuid.UUID,
		state MonitorRunState,
	) (*Monitor, *apperr.ServiceError)
	RunMonitorManuallyByID(ctx context.Context, monitorUUID uuid.UUID) *apperr.ServiceError
	ReceivePush(ctx context.Context, monitorID uuid.UUID, payload probe.PushProbePayload) *apperr.ServiceError
	ExportToYAML(ctx context.Context, monitorID uuid.UUID, w io.Writer) *apperr.ServiceError
}

// MonitorService handles monitor-related CRUD operations.
type MonitorService struct {
	db         db.DB
	monitorDAO IMonitorDAO
	userDAO    users.IUserDAO
}

type MonitorServiceDeps struct {
	DB db.DB
}

func NewMonitorService(deps MonitorServiceDeps) *MonitorService {
	return &MonitorService{
		db:         deps.DB,
		monitorDAO: NewMonitorDAO(deps.DB.Querier()),
		userDAO:    users.NewUserDAO(deps.DB.Querier()),
	}
}

// CreateMonitor creates a new monitor.
func (s *MonitorService) CreateMonitor(
	ctx context.Context,
	monitor Monitor,
) (*Monitor, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "CreateMonitor")
	logger.Trace().
		Interface("monitor", monitor).
		Msg("Creating monitor")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	owner, err := s.userDAO.GetUserByUsername(ctx, userClaims.Username)
	if err != nil {
		logger.Error().Err(err).Str("username", userClaims.Username).Msg("Failed to find creating user")
		return nil, apperr.NewInternalError("failed to find creating user: %w", err)
	}

	initializedMonitor := NewMonitorFromPayload(monitor, owner.ID)

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
	return monitorFromDB, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorService) DeleteMonitor(ctx context.Context, id string) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return apperr.NewInternalError("user claims not found in context")
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

	if monitorBeforeDelete.Source == MonitorSourceConfig {
		return apperr.NewConflictError("monitor %s is managed by config file", id)
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
func (s *MonitorService) UpdateMonitor(ctx context.Context, monitor Monitor) (*Monitor, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	updatedMonitor, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*Monitor, *audit.AuditLogParams, error) {
		existingMonitor, err := NewMonitorDAO(q).GetMonitorByID(ctx, monitor.ID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitor.ID.String()).Msg("Monitor not found")
				return nil, nil, apperr.NewNotFoundError("monitor with ID %s not found", monitor.ID)
			}
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to retrieve existing monitor for update")
			return nil, nil, fmt.Errorf("failed to retrieve existing monitor for update: %w", err)
		}

		if existingMonitor.Source == MonitorSourceConfig {
			return nil, nil, apperr.NewConflictError("monitor %s is managed by config file", monitor.ID)
		}

		if monitor.Type != existingMonitor.Type {
			return nil, nil, apperr.NewBadRequestError("cannot edit monitor type")
		}

		monitor.RunState = existingMonitor.RunState
		monitor.OwnerID = existingMonitor.OwnerID

		if err := monitor.Validate(); err != nil {
			logger.Error().
				Err(err).
				Str("id", monitor.ID.String()).
				Interface("monitor", monitor).
				Msg("Invalid monitor configuration")
			return nil, nil, apperr.NewBadRequestError("invalid monitor configuration: %w", err)
		}

		_, err = NewMonitorDAO(q).UpdateMonitor(ctx, monitor)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to update monitor in database")
			return nil, nil, fmt.Errorf("failed to update monitor in database: %w", err)
		}

		updatedMonitor, err := NewMonitorDAO(q).GetMonitorByID(ctx, monitor.ID)
		if err != nil {
			logger.Error().Err(err).Str("id", monitor.ID.String()).Msg("Failed to retrieve updated monitor")
			return nil, nil, fmt.Errorf("failed to retrieve updated monitor: %w", err)
		}

		params := &audit.AuditLogParams{
			Username:   &userClaims.Username,
			ResourceID: &monitor.ID,
			Action:     audit.ActionUpdateMonitor,
			IsSuccess:  true,
			Summary:    fmt.Sprintf("Monitor with ID %s updated", monitor.ID),
			Before:     existingMonitor,
			After:      updatedMonitor,
		}

		return updatedMonitor, params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return nil, serviceErr
		}
		logger.Error().Err(txErr).Str("id", monitor.ID.String()).Msg("Failed to update monitor within transaction")
		return nil, apperr.NewInternalError("failed to update monitor within transaction: %w", txErr)
	}

	MonitorLifecycleChannel.Broadcast(MonitorLifecycleMessage{
		ID:      monitor.ID,
		Status:  Edited,
		Monitor: nil,
	})

	logger.Debug().Str("id", monitor.ID.String()).Msg("Monitor updated")
	return updatedMonitor, nil
}

func (s *MonitorService) UpdateMonitorStateByID(
	ctx context.Context,
	monitorID uuid.UUID,
	state MonitorRunState,
) (*Monitor, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "UpdateMonitorStateByID")
	logger.Trace().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Updating monitor state by ID")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		logger.Error().Msg("User claims not found in context")
		return nil, apperr.NewInternalError("user claims not found in context")
	}

	if !isValidMonitorState(string(state)) {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Invalid monitor state provided")
		return nil, apperr.NewBadRequestError("invalid monitor state: %s", state)
	}

	updatedMonitor, txErr := audit.WithAuditedTx(ctx, s.db, func(q db.Querier) (*Monitor, *audit.AuditLogParams, error) {
		monitor, err := NewMonitorDAO(q).GetMonitorByID(ctx, monitorID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Error().Str("id", monitorID.String()).Msg("Monitor not found for state update")
				return nil, nil, apperr.NewNotFoundError("monitor with ID %s not found", monitorID.String())
			}
			logger.Error().Err(err).Str("id", monitorID.String()).Msg("Failed to retrieve monitor for state update")
			return nil, nil, apperr.NewInternalError("failed to retrieve monitor for state update: %w", err)
		}

		if monitor.Source == MonitorSourceConfig {
			return nil, nil, apperr.NewConflictError("monitor %s is managed by config file", monitorID)
		}

		if monitor.RunState == state {
			logger.Warn().
				Str("id", monitorID.String()).
				Str("state", string(state)).
				Msg("Nothing to update in monitor, no update needed")
			return monitor, nil, nil
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
			return nil, nil, apperr.NewInternalError("failed to update monitor state in database: %w", updateErr)
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

		return monitor, params, nil
	})
	if txErr != nil {
		if serviceErr, ok2 := errors.AsType[*apperr.ServiceError](txErr); ok2 {
			return nil, serviceErr
		}
		return nil, apperr.NewInternalError("failed to update monitor state: %w", txErr)
	}

	MonitorLifecycleChannel.Broadcast(MonitorLifecycleMessage{
		ID:     monitorID,
		Status: Edited,
	})

	logger.Debug().Str("id", monitorID.String()).Str("newState", string(state)).Msg("Monitor state updated")
	return updatedMonitor, nil
}

func (s *MonitorService) RunMonitorManuallyByID(ctx context.Context, monitorUUID uuid.UUID) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "RunMonitorManuallyByID")
	logger.Trace().Str("id", monitorUUID.String()).Msg("Running monitor manually by ID")

	userClaims, ok := auth.GetUserClaimsFromContext(ctx)
	if !ok {
		return apperr.NewInternalError("user claims not found in context")
	}

	monitor, err := NewMonitorDAO(s.db.Querier()).GetMonitorByID(ctx, monitorUUID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apperr.NewNotFoundError("monitor with ID %s not found", monitorUUID.String())
		}
		return apperr.NewInternalError("failed to retrieve monitor for manual run: %w", err)
	}

	if monitor.Type == kind.PushConfigType {
		return apperr.NewBadRequestError("cannot manually run a push monitor")
	}

	if err := monitor.Validate(); err != nil {
		return apperr.NewBadRequestError("monitor validation failed: %w", err)
	}

	if _, err := probe.Parse[probe.Probe](monitor.Type, monitor.ProbeConfig); err != nil {
		if errors.Is(err, probe.ErrInvalidProbeConfig) {
			return apperr.NewBadRequestError("probe config validation failed for manual run: %w", err)
		}
		return apperr.NewInternalError("failed to unmarshal probe config for manual run: %w", err)
	}

	auditErr := audit.NewAuditLogDAO(s.db.Querier()).Record(ctx, audit.AuditLogParams{
		Username:   &userClaims.Username,
		ResourceID: &monitor.ID,
		Action:     audit.ActionRunMonitorManually,
		IsSuccess:  true,
		Summary:    fmt.Sprintf("Monitor with ID %s run manually", monitor.ID),
	})
	if auditErr != nil {
		return apperr.NewInternalError("failed to record audit log for manual run: %w", auditErr)
	}

	traceID, _ := log.TraceIDFromContext(ctx)
	MonitorExecuteChannel.Broadcast(MonitorExecuteMessage{Monitor: *monitor, TraceID: traceID, IsManuallyTriggered: true})

	return nil
}

func (s *MonitorService) ReceivePush(ctx context.Context, monitorID uuid.UUID, payload probe.PushProbePayload) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "ReceivePush")
	logger.Trace().Str("id", monitorID.String()).Any("payload", payload).Msg("Received push payload")

	monitor, err := NewMonitorDAO(s.db.Querier()).GetMonitorByID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apperr.NewNotFoundError("monitor with ID %s not found", monitorID.String())
		}
		return apperr.NewInternalError("failed to retrieve monitor for push: %w", err)
	}
	if monitor.Type != kind.PushConfigType {
		return apperr.NewNotFoundError("push monitor with ID %s not found", monitorID.String())
	}

	if err := monitor.Validate(); err != nil {
		return apperr.NewBadRequestError("monitor validation failed: %w", err)
	}

	if _, err := probe.Parse[*probe.PushProbe](monitor.Type, monitor.ProbeConfig); err != nil {
		if errors.Is(err, probe.ErrInvalidProbeConfig) {
			return apperr.NewBadRequestError("probe config validation failed: %w", err)
		}
		return apperr.NewInternalError("failed to unmarshal probe config: %w", err)
	}

	if monitor.RunState != MonitorRunStateActive {
		return apperr.NewConflictError("push to a paused monitor")
	}

	result := results.NewMonitorResult(monitor.ID, monitor.Type, kind.MonitorStatusUp, false, payload.Latency, &results.PushResultDetails{RawMessage: payload.Body})
	if payload.Status == kind.MonitorStatusDown {
		result.AddFailure(results.FailureReasonPushReportedDown, nil, nil)
	}

	MonitorRunChannel.Broadcast(MonitorRunMessage{
		Monitor: *monitor,
		Result:  &result,
	})

	return nil
}

func (s *MonitorService) ExportToYAML(ctx context.Context, monitorID uuid.UUID, w io.Writer) *apperr.ServiceError {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitor, "ExportToYAML")
	logger.Trace().Str("id", monitorID.String()).Msg("Exporting to YAML")

	monitorsDAO := NewMonitorDAO(s.db.Querier())
	tagsDAO := tags.NewTagDAO(s.db.Querier())

	monitor, err := monitorsDAO.GetMonitorByID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return apperr.NewNotFoundError("monitor with ID %s not found", monitorID.String())
		}
		return apperr.NewInternalError("failed to retrieve monitor: %w", err)
	}

	if err = exportToYAML(ctx, *monitor, tagsDAO, w); err != nil {
		return apperr.NewInternalError("failed to export monitor %s to YAML", monitor.Slug)
	}

	return nil
}
