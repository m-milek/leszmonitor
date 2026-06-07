package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/api/authorization"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/util"
)

type IMonitorResultsService interface {
	GetLatestMonitorResultByMonitorID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitorID string) (monitorresult.IMonitorResult, *ServiceError)
	GetMonitorResultsByMonitorID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string, pagination *util.Pagination) ([]monitorresult.IMonitorResult, *ServiceError)
}

type MonitorResultsService struct {
	db   db.DB
	auth IAuthorizer
}

type MonitorResultsServiceDeps struct {
	DB   db.DB
	Auth IAuthorizer
}

func NewMonitorResultsService(deps MonitorResultsServiceDeps) *MonitorResultsService {
	return &MonitorResultsService{
		db:   deps.DB,
		auth: deps.Auth,
	}
}

func (s *MonitorResultsService) GetLatestMonitorResultByMonitorID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, monitorID string) (monitorresult.IMonitorResult, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorResultsService", "GetLatestMonitorResultByMonitorID")

	_, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		logger.Warn().Err(authErr).Msg("Unauthorized access to GetLatestMonitorResultByMonitorID")
		return nil, &ServiceError{Code: 403, Err: authErr}
	}

	result, err := s.db.MonitorResults().GetLatestMonitorResultByMonitorID(ctx, monitorID)
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

func (s *MonitorResultsService) GetMonitorResultsByMonitorID(ctx context.Context, projectAuth *authorization.ProjectAuthorization, id string, pagination *util.Pagination) ([]monitorresult.IMonitorResult, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorResultsService", "GetMonitorResultsByMonitorID")

	_, authErr := s.auth.authorizeProjectAction(ctx, projectAuth, models.PermissionMonitorReader)
	if authErr != nil {
		logger.Warn().Err(authErr).Msg("Unauthorized access to GetMonitorResultsByMonitorID")
		return nil, &ServiceError{Code: 403, Err: authErr}
	}

	results, err := s.db.MonitorResults().GetMonitorResultsByMonitorID(ctx, id, pagination)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Info().Str("monitorID", id).Msg("No monitor results found for given monitor ID")
			return nil, &ServiceError{Code: 404, Err: err}
		}
		logger.Error().Err(err).Msg("Failed to get monitor results by monitor ID")
		return nil, &ServiceError{Code: 500, Err: err}
	}

	return results, nil
}
