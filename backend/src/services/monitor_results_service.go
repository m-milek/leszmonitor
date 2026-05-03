package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/api/middleware"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/monitorresult"
)

type MonitorResultsServiceT struct {
	baseService
}

func newMonitorResultsService(service baseService) *MonitorResultsServiceT {
	return &MonitorResultsServiceT{baseService: service}
}

var MonitorResultsService = newMonitorResultsService(newBaseService(newAuthorizationService(), "MonitorResultsService"))

func (s *MonitorResultsServiceT) GetLatestMonitorResultByMonitorID(ctx context.Context, projectAuth *middleware.ProjectAuth, monitorID string) (monitorresult.IMonitorResult, *ServiceError) {
	logger := s.getMethodLogger("GetLatestMonitorResultByMonitorID")

	_, authErr := s.authService.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		logger.Warn().Err(authErr).Msg("Unauthorized access to GetLatestMonitorResultByMonitorID")
		return nil, &ServiceError{Code: 403, Err: authErr}
	}

	result, err := s.getDB().MonitorResults().GetLatestMonitorResultByMonitorID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Info().Str("monitorID", monitorID).Msg("No monitor result found for given monitor ID")
			return nil, &ServiceError{Code: 404, Err: err}
		}
		logger.Error().Err(err).Msg("Failed to get latest monitor result by monitor ID")
		return nil, &ServiceError{Code: 500, Err: err}
	}

	return result, nil
}
