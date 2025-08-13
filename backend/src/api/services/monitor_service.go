package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/logger"
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
			logger: logger.NewServiceLogger("MonitorService"),
		},
	}
}

var MonitorService = newMonitorService()

type MonitorCreateResponse struct {
	MonitorId string `json:"monitorId"`
}

func (s *MonitorServiceT) CreateMonitor(ctx context.Context, monitor monitors.IMonitor) (*MonitorCreateResponse, *ServiceError) {
	s.logger.Trace().Interface("monitor", monitor).Msg("Creating new monitor")

	monitor.GenerateId()

	if err := monitor.Validate(); err != nil {
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration: %w", err),
		}
	}

	_, err := db.CreateMonitor(ctx, monitor)
	if err != nil {
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to add monitor to database: %w", err),
		}
	}

	return &MonitorCreateResponse{
		MonitorId: monitor.GetId(),
	}, nil
}

func (s *MonitorServiceT) DeleteMonitor(ctx context.Context, id string) *ServiceError {
	s.logger.Trace().Str("id", id).Msg("Deleting monitor")

	wasDeleted, err := db.DeleteMonitor(ctx, id)
	if err != nil {
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to delete monitor: %w", err),
		}
	}

	if !wasDeleted {
		return &ServiceError{
			Code: http.StatusNotFound,
			Err:  fmt.Errorf("monitor not found or already deleted"),
		}
	}

	return nil
}

func (s *MonitorServiceT) GetAllMonitors(ctx context.Context) ([]monitors.IMonitor, *ServiceError) {
	s.logger.Trace().Msg("Retrieving all monitors")

	monitorsList, err := db.GetAllMonitors(ctx)
	if err != nil {
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitors: %w", err),
		}
	}

	return monitorsList, nil
}

func (s *MonitorServiceT) GetMonitorById(ctx context.Context, id string) (monitors.IMonitor, *ServiceError) {
	s.logger.Trace().Str("id", id).Msg("Retrieving monitor by ID")

	monitor, err := db.GetMonitorById(ctx, id)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor with ID %s not found", id),
			}
		}
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitor: %w", err),
		}
	}

	return monitor, nil
}

func (s *MonitorServiceT) UpdateMonitor(ctx context.Context, monitor monitors.IMonitor) *ServiceError {
	s.logger.Trace().Interface("monitor", monitor).Msg("Updating monitor")

	if err := monitor.Validate(); err != nil {
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration for update: %w", err),
		}
	}

	wasUpdated, err := db.UpdateMonitor(ctx, monitor)
	if err != nil {
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update monitor in database: %w", err),
		}
	}

	if !wasUpdated {
		if errors.Is(err, db.ErrNotFound) {
			return &ServiceError{
				Code: http.StatusNotFound,
				Err:  fmt.Errorf("monitor with ID %s not found", monitor.GetId()),
			}
		}
		return &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to update monitor: %w", err),
		}
	}

	return nil
}
