package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/monitors"
)

// MonitorServiceT handles monitor-related CRUD operations.
type MonitorServiceT struct {
	baseService
}

func newMonitorService() *MonitorServiceT {
	return &MonitorServiceT{
		baseService{
			authService:   newAuthorizationService(),
			serviceLogger: log.NewServiceLogger("MonitorService"),
		},
	}
}

var MonitorService = newMonitorService()

type MonitorCreateResponse struct {
	MonitorID string `json:"monitorId"`
}

// CreateMonitor creates a new monitor in the specified project.
func (s *MonitorServiceT) CreateMonitor(ctx context.Context, projectAuth *middleware.ProjectAuth, monitor monitors.IConcreteMonitor) (*MonitorCreateResponse, *ServiceError) {
	logger := s.getMethodLogger("CreateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Creating new monitor")

	project, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorEditor)
	if authErr != nil {
		return nil, authErr
	}

	monitor.SetProjectSlug(project.Slug)
	monitor.GenerateSlug()

	if err := monitor.Validate(); err != nil {
		logger.Warn().Err(err).Msg("Invalid monitor configuration")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration: %w", err),
		}
	}

	dbRes, createErr := s.getDB().Monitors().InsertMonitor(ctx, monitor)
	if createErr != nil {
		logger.Error().Err(createErr).Msg("Failed to add monitor to database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to add monitor to database: %w", createErr),
		}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      dbRes.GetID(),
		Status:  monitors.Created,
		Monitor: &dbRes,
	})

	logger.Info().Str("id", monitor.GetID().String()).Msg("Monitor created")
	return &MonitorCreateResponse{MonitorID: dbRes.GetID().String()}, nil
}

// DeleteMonitor deletes a monitor by its slug.
func (s *MonitorServiceT) DeleteMonitor(ctx context.Context, projectAuth *middleware.ProjectAuth, id string) *ServiceError {
	logger := s.getMethodLogger("DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	_, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorAdmin)
	if authErr != nil {
		return authErr
	}

	deletedID, err := s.getDB().Monitors().DeleteMonitorBySlug(ctx, id)
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

	logger.Info().Str("id", id).Msg("Monitor deleted")
	return nil
}

// GetMonitorsByProjectID retrieves all monitors for the project in the provided ProjectAuth.
func (s *MonitorServiceT) GetMonitorsByProjectID(ctx context.Context, projectAuth *middleware.ProjectAuth) ([]monitors.IConcreteMonitor, *ServiceError) {
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
func (s *MonitorServiceT) GetMonitorByID(ctx context.Context, projectAuth *middleware.ProjectAuth, id string) (monitors.IMonitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorByID")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by slug")

	_, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitor, err := s.getDB().Monitors().GetMonitorBySlug(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("id", id).Msg("Monitor not found in database")
			return nil, &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with slug %s not found", id)}
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor from database")
		return nil, &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to retrieve monitor: %w", err)}
	}

	if monitor.GetProjectSlug() != projectAuth.ProjectID {
		logger.Warn().Str("id", id).Msg("Monitor does not belong to the authorized project")
		return nil, &ServiceError{Code: http.StatusForbidden, Err: fmt.Errorf("monitor with slug %s does not belong to the authorized project", id)}
	}

	return monitor, nil
}

// UpdateMonitor updates an existing monitor's configuration.
func (s *MonitorServiceT) UpdateMonitor(ctx context.Context, projectAuth *middleware.ProjectAuth, monitor monitors.IConcreteMonitor) *ServiceError {
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
			return &ServiceError{Code: http.StatusNotFound, Err: fmt.Errorf("monitor with slug %s not found", monitor.GetID())}
		}
		logger.Error().Err(err).Msg("Failed to update monitor in database")
		return &ServiceError{Code: http.StatusInternalServerError, Err: fmt.Errorf("failed to update monitor: %w", err)}
	}

	events.MonitorLifecycleChannel.Broadcast(monitors.MonitorLifecycleMessage{
		ID:      monitor.GetID(),
		Status:  monitors.Edited,
		Monitor: nil,
	})

	logger.Info().Str("id", monitor.GetID().String()).Msg("Monitor updated")
	return nil
}
