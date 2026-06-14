package services

import (
	"context"
	"errors"

	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/util"
)

type IMonitorResultsService interface {
	GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (monitorresult.IMonitorResult, *ServiceError)
	GetMonitorResultsByMonitorID(ctx context.Context, id string, pagination *util.Pagination) ([]monitorresult.IMonitorResult, *ServiceError)
}

type MonitorResultsService struct {
	db db.DB
}

type MonitorResultsServiceDeps struct {
	DB db.DB
}

func NewMonitorResultsService(deps MonitorResultsServiceDeps) *MonitorResultsService {
	return &MonitorResultsService{
		db: deps.DB,
	}
}

func (s *MonitorResultsService) GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (monitorresult.IMonitorResult, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorResultsService", "GetLatestMonitorResultByMonitorID")

	result, err := s.db.MonitorResults().GetLatestMonitorResultByMonitorID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Info().Str("monitorID", monitorID).Msg("No monitor result found for given monitor ID")
			return nil, NewNotFoundError("no monitor result found: %w", err)
		}
		logger.Error().Err(err).Msg("Failed to get latest monitor result by monitor ID")
		return nil, NewInternalError("failed to get latest monitor result: %w", err)
	}

	return result, nil
}

func (s *MonitorResultsService) GetMonitorResultsByMonitorID(ctx context.Context, id string, pagination *util.Pagination) ([]monitorresult.IMonitorResult, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, "MonitorResultsService", "GetMonitorResultsByMonitorID")

	results, err := s.db.MonitorResults().GetMonitorResultsByMonitorID(ctx, id, pagination)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Info().Str("monitorID", id).Msg("No monitor results found for given monitor ID")
			return nil, NewNotFoundError("no monitor results found: %w", err)
		}
		logger.Error().Err(err).Msg("Failed to get monitor results by monitor ID")
		return nil, NewInternalError("failed to get monitor results: %w", err)
	}

	return results, nil
}
