package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logging"
	monitors "github.com/m-milek/leszmonitor/uptime/monitor"
	"net/http"
)

type MonitorServiceT struct {
	BaseService
}

// NewMonitorService creates a new instance of MonitorServiceT.
func newMonitorService() *MonitorServiceT {
	return &MonitorServiceT{
		BaseService{
			serviceLogger: logging.NewServiceLogger("MonitorService"),
		},
	}
}

var MonitorService = newMonitorService()

type MonitorCreateResponse struct {
	MonitorId string `json:"monitorId"`
}

func (s *MonitorServiceT) CreateMonitor(ctx context.Context, monitor monitors.IMonitor) (*MonitorCreateResponse, *ServiceError) {
	logger := s.getMethodLogger("CreateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Creating new monitor")

	monitor.GenerateId()

	if err := monitor.Validate(); err != nil {
		logger.Warn().Err(err).Msg("Invalid monitor configuration")
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration: %w", err),
		}
	}

	_, err := db.CreateMonitor(ctx, monitor)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to add monitor to database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to add monitor to database: %w", err),
		}
	}

	logger.Info().Str("id", monitor.GetId()).Msg("Monitor created")
	return &MonitorCreateResponse{
		MonitorId: monitor.GetId(),
	}, nil
}

func (s *MonitorServiceT) DeleteMonitor(ctx context.Context, id string) *ServiceError {
	logger := s.getMethodLogger("DeleteMonitor")
	logger.Trace().Str("id", id).Msg("Deleting monitor")

	wasDeleted, err := db.DeleteMonitor(ctx, id)
	if err != nil {
		logger.Error().Err(err).Str("id", id).Msg("Failed to delete monitor from database")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor: %w", err),
		}
	}

	if !wasDeleted {
		logger.Warn().Str("id", id).Msg("Monitor not found or already deleted")
		return &ServiceError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("monitor not found or already deleted"),
		}
	}

	logger.Info().Str("id", id).Msg("Monitor deleted")
	return nil
}

func (s *MonitorServiceT) GetAllMonitors(ctx context.Context) ([]monitors.IMonitor, *ServiceError) {
	logger := s.getMethodLogger("GetAllMonitors")
	logger.Trace().Msg("Retrieving all monitors")

	monitorsList, err := db.GetAllMonitors(ctx)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to retrieve monitors from database")
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitors: %w", err),
		}
	}

	return monitorsList, nil
}

func (s *MonitorServiceT) GetMonitorById(ctx context.Context, id string) (monitors.IMonitor, *ServiceError) {
	logger := s.getMethodLogger("GetMonitorById")
	logger.Trace().Str("id", id).Msg("Retrieving monitor by ID")

	monitor, err := db.GetMonitorById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("id", id).Msg("Monitor not found")
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor with ID %s not found", id),
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

func (s *MonitorServiceT) UpdateMonitor(ctx context.Context, monitor monitors.IMonitor) *ServiceError {
	logger := s.getMethodLogger("UpdateMonitor")
	logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	if err := monitor.Validate(); err != nil {
		logger.Warn().Err(err).Msg("Invalid monitor configuration for update")
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration for update: %w", err),
		}
	}

	wasUpdated, err := db.UpdateMonitor(ctx, monitor)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Interface("monitor", monitor).Msg("Monitor not found for update")
			return &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor with ID %s not found", monitor.GetId()),
			}
		}
		logger.Error().Err(err).Interface("monitor", monitor).Msg("Failed to update monitor in database")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update monitor in database: %w", err),
		}
	}

	if !wasUpdated {
		logger.Warn().Interface("monitor", monitor).Msg("Monitor was not updated, possibly not found")
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("monitor with ID %s was not updated", monitor.GetId()),
		}
	}

	logger.Info().Str("id", monitor.GetId()).Msg("Monitor updated")
	return nil
}
