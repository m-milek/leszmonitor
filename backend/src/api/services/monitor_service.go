package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logging"
	"github.com/m-milek/leszmonitor/models"
	monitors "github.com/m-milek/leszmonitor/uptime/monitor"
	"net/http"
)

// MonitorServiceT handles monitor-related CRUD operations.
type MonitorServiceT struct {
	baseService
}

// NewMonitorService creates a new instance of MonitorServiceT.
func newMonitorService() *MonitorServiceT {
	return &MonitorServiceT{
		baseService{
			authService:   newAuthorizationService(),
			serviceLogger: logging.NewServiceLogger("MonitorService"),
		},
	}
}

var MonitorService = newMonitorService()

type MonitorCreateResponse struct {
	MonitorID string `json:"monitorId"`
}

// CreateMonitor creates a new monitor in the specified group.
func (s *MonitorServiceT) CreateMonitor(ctx context.Context, teamAuth *middleware.TeamAuth, groupID string, monitor monitors.IConcreteMonitor) (*MonitorCreateResponse, *ServiceError) {
	logger := s.getMethodLogger("InsertMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Creating new monitor")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionMonitorEditor)
	if authErr != nil {
		return nil, authErr
	}

	group, err := GroupService.internalGetMonitorGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	monitor.SetGroupID(group.ID)
	monitor.SetTeamID(team.ID)
	monitor.GenerateDisplayID()

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

	// Broadcast that a monitor has been added
	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		ID:      dbRes.GetID(),
		Status:  monitors.Created,
		Monitor: &dbRes,
	})

	logger.Info().Str("id", monitor.GetID().String()).Msg("Monitor created")
	return &MonitorCreateResponse{
		MonitorID: monitor.GetID().String(),
	}, nil
}

// DeleteMonitor deletes a monitor by its DisplayID.
func (s *MonitorServiceT) DeleteMonitor(ctx context.Context, teamAuth *middleware.TeamAuth, id string) *ServiceError {
	logger := s.getMethodLogger("DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	_, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionMonitorAdmin)
	if authErr != nil {
		return authErr
	}

	deletedID, err := s.getDB().Monitors().DeleteMonitorByDisplayID(ctx, id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to delete monitor from database")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor: %w", err),
		}
	}

	if deletedID == nil {
		logger.Warn().Str("id", id).Msg("Monitor not found or already deleted")
		return &ServiceError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("monitor not found or already deleted"),
		}
	}

	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		ID:      *deletedID,
		Status:  monitors.Deleted,
		Monitor: nil,
	})

	logger.Info().Str("id", id).Msg("Monitor deleted")
	return nil
}

// GetMonitorsByTeamID retrieves all monitors for the team in the provided TeamAuth.
func (s *MonitorServiceT) GetMonitorsByTeamID(ctx context.Context, teamAuth *middleware.TeamAuth) ([]monitors.IConcreteMonitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorsByTeamID")
	logger.Trace().Msg("Retrieving all monitors")

	team, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitorsList, err := s.getDB().Monitors().GetMonitorsByTeamID(ctx, team.ID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitors: %w", err),
		}
	}

	return monitorsList, nil
}

// GetMonitorByID retrieves a specific monitor by its DisplayID.
func (s *MonitorServiceT) GetMonitorByID(ctx context.Context, teamAuth *middleware.TeamAuth, id string) (monitors.IMonitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorByID")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by DisplayID")

	_, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionMonitorReader)
	if authErr != nil {
		return nil, authErr
	}

	monitor, err := s.getDB().Monitors().GetMonitorByID(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("id", id).Msg("Monitor not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor with DisplayID %s not found", id),
			}
		}
		logger.Error().Err(err).Str("id", id).Msg("Failed to retrieve monitor from database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitor: %w", err),
		}
	}

	return monitor, nil
}

// UpdateMonitor updates an existing monitor's configuration.
func (s *MonitorServiceT) UpdateMonitor(ctx context.Context, teamAuth *middleware.TeamAuth, monitor monitors.IConcreteMonitor) *ServiceError {
	logger := s.getMethodLogger("UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	if err := monitor.Validate(); err != nil {
		logger.Warn().Err(err).Msg("Invalid monitor configuration for update")
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration for update: %w", err),
		}
	}

	_, authErr := s.authService.authorizeTeamAction(ctx, teamAuth, models.PermissionMonitorEditor)
	if authErr != nil {
		return authErr
	}

	updatedMonitor, err := s.getDB().Monitors().UpdateMonitor(ctx, monitor)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Interface("monitor", monitor).Msg("Monitor not found for update")
			return &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor with DisplayID %s not found", monitor.GetID()),
			}
		}
		logger.Error().Err(err).Interface("monitor", monitor).Msg("Failed to update monitor in database")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update monitor in database: %w", err),
		}
	}

	if updatedMonitor == nil {
		logger.Warn().Interface("monitor", monitor).Msg("Monitor was not updated, possibly not found")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("monitor with DisplayID %s was not updated", monitor.GetID()),
		}
	}

	monitors.MessageBroadcaster.Broadcast(monitors.MonitorMessage{
		ID:      updatedMonitor.GetID(),
		Status:  monitors.Edited,
		Monitor: &updatedMonitor,
	})

	logger.Info().Str("id", monitor.GetID().String()).Msg("Monitor updated")
	return nil
}
