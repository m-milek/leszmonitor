package services

import (
	"errors"
	"fmt"
	"github.com/m-milek/leszmonitor/db"
	monitors "github.com/m-milek/leszmonitor/uptime/monitor"
	"net/http"
)

type MonitorServiceT struct{}

// NewMonitorService creates a new instance of MonitorServiceT.
func newMonitorService() *MonitorServiceT {
	return &MonitorServiceT{}
}

var MonitorService = newMonitorService()

type MonitorCreateResponse struct {
	MonitorId string `json:"monitorId"`
}

func (s *MonitorServiceT) CreateMonitor(monitor monitors.IMonitor) (*MonitorCreateResponse, *ServiceError) {
	monitor.GenerateId()

	if err := monitor.Validate(); err != nil {
		return nil, &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration: %w", err),
		}
	}

	_, err := db.CreateMonitor(monitor)
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

func (s *MonitorServiceT) DeleteMonitor(id string) *ServiceError {
	wasDeleted, err := db.DeleteMonitor(id)
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

func (s *MonitorServiceT) GetAllMonitors() ([]monitors.IMonitor, *ServiceError) {
	monitorsList, err := db.GetAllMonitors()
	if err != nil {
		return nil, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Errorf("failed to retrieve monitors: %w", err),
		}
	}

	return monitorsList, nil
}

func (s *MonitorServiceT) GetMonitorById(id string) (monitors.IMonitor, *ServiceError) {
	monitor, err := db.GetMonitorById(id)
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

func (s *MonitorServiceT) UpdateMonitor(monitor monitors.IMonitor) *ServiceError {
	if err := monitor.Validate(); err != nil {
		return &ServiceError{
			Code: http.StatusBadRequest,
			Err:  fmt.Errorf("invalid monitor configuration for update: %w", err),
		}
	}

	wasUpdated, err := db.UpdateMonitor(monitor)
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
