package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/monitors"
	"github.com/m-milek/leszmonitor/security"
)

// MonitorServiceT handles monitor-related CRUD operations.
type MonitorServiceT struct {
	baseService
}

func newMonitorService() *MonitorServiceT {
	return &MonitorServiceT{
		baseService{
			authService:     newAuthorizationService(),
			auditLogService: newAuditLogService(),
			serviceLogger:   log.NewServiceLogger("MonitorService"),
		},
	}
}

var MonitorService = newMonitorService()

type MonitorCreateResponse struct {
	MonitorID string `json:"monitorId"`
}

// CreateMonitor creates a new monitor in the specified project.
func (s *MonitorServiceT) CreateMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitor monitors.Monitor) (*MonitorCreateResponse, *ServiceError) {
	logger := s.getMethodLogger("CreateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Creating new monitor")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorEditor)
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

	dbRes, createErr := s.getDB().Monitors().InsertMonitor(ctx, *initializedMonitor)
	if createErr != nil {
		logger.Error().Err(createErr).Msg("Failed to add monitor to database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to add monitor to database: %w", createErr),
		}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      dbRes.ID,
		Status:  monitors.Created,
		Monitor: dbRes,
	})

	logger.Info().Str("id", monitor.ID.String()).Msg("Monitor created")
	return &MonitorCreateResponse{MonitorID: dbRes.ID.String()}, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorServiceT) DeleteMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string) *ServiceError {
	logger := s.getMethodLogger("DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	_, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorAdmin)
	if authErr != nil {
		return authErr
	}

	monitorUUID, err := uuid.Parse(id)
	if err != nil {
		logger.Warn().Str("id", id).Msg("Invalid monitor ID format")
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid monitor ID format: %w", err)}
	}

	deletedID, err := s.getDB().Monitors().DeleteMonitorByID(ctx, monitorUUID)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to delete monitor from database")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor: %w", err),
		}
	}

	if deletedID == nil {
		return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor not found or already deleted")}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      *deletedID,
		Status:  monitors.Deleted,
		Monitor: nil,
	})

	err = s.auditLogService.Record(ctx, security.AuditLogEntry{
		Username:   &projectAuth.Username,
		ProjectID:  &projectAuth.ProjectID,
		ResourceID: deletedID,
		Action:     security.ActionDeleteMonitor,
		IsSuccess:  true,
		Summary:    fmt.Sprintf("Monitor with ID %s deleted", monitorUUID.String()),
		Before:     nil,
		After:      nil,
		TraceID:    security.GetTraceIDFromContext(ctx),
	})
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to record audit log entry for monitor deletion")
		return nil
	}

	logger.Info().Str("id", id).Msg("Monitor deleted")
	return nil
}

// GetMonitorsByProjectID retrieves all monitors for the project in the provided ProjectAuth.
func (s *MonitorServiceT) GetMonitorsByProjectID(ctx context.Context, projectAuth *authorization.ProjectAuthorization) ([]monitors.Monitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorsByProjectID")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitorsList, err := s.getDB().Monitors().GetMonitorsByProjectID(ctx, project.ID)
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
func (s *MonitorServiceT) GetMonitorByID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string) (*monitors.Monitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorByID")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by slug")

	_, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	// TODO fix - monitor ID uses slug - bug. This method is not used by the frontend
	monitor, err := s.getDB().Monitors().GetMonitorBySlug(ctx, id, projectAuth.ProjectID)
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
func (s *MonitorServiceT) UpdateMonitor(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitor monitors.Monitor) *ServiceError {
	logger := s.getMethodLogger("UpdateMonitor")

	if err := monitor.Validate(); err != nil {
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid monitor configuration: %w", err)}
	}

	_, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorEditor)
	if authErr != nil {
		return authErr
	}

	_, err := s.getDB().Monitors().UpdateMonitor(ctx, monitor)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with ID %s not found", monitor.ID)}
		}
		logger.Error().Err(err).Msg("Failed to update monitor in database")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to update monitor: %w", err)}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitor.ID,
		Status:  monitors.Edited,
		Monitor: nil,
	})

	logger.Info().Str("id", monitor.ID.String()).Msg("Monitor updated")
	return nil
}

func (s *MonitorServiceT) GetMonitorBySlugByProject(ctx context.Context, auth *authorization.ProjectAuthorization, slug string) (*monitors.Monitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorBySlugByProject")

	_, authErr := s.authService.authorizeProjectAction(ctx, auth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitor, err := s.getDB().Monitors().GetMonitorBySlugByProject(ctx, slug, auth.ProjectID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with slug %s not found in project", slug)}
		}
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve monitor by slug and project: %w", err)}
	}

	logger.Debug().Str("slug", slug).Str("projectID", auth.ProjectID.String()).Msg("Monitor retrieved by slug and project")
	return monitor, nil
}

func (s *MonitorServiceT) UpdateMonitorStateByID(ctx context.Context, auth *authorization.ProjectAuthorization, monitorID uuid.UUID, state monitors.MonitorState) *ServiceError {
	logger := s.getMethodLogger("UpdateMonitorStateByID")

	if !monitors.IsValidMonitorState(string(state)) {
		logger.Warn().Str("id", monitorID.String()).Str("state", string(state)).Msg("Invalid monitor state provided")
		return &ServiceError{Code: http.StatusBadRequest, Err: fmt.Errorf("invalid monitor state: %s", state)}
	}

	_, authErr := s.authService.authorizeProjectAction(ctx, auth, models.PermissionMonitorEditor)
	if authErr != nil {
		return authErr
	}

	monitor, err := s.getDB().Monitors().GetMonitorByID(ctx, monitorID)
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

	_, updateErr := s.getDB().Monitors().UpdateMonitor(ctx, *monitor)
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
