package results

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
)

type IMonitorResultsService interface {
	GetLatestMonitorResultByMonitorID(
		ctx context.Context,
		monitorID uuid.UUID,
	) (IMonitorResult, *apperr.ServiceError)
	GetMonitorResultsByMonitorID(
		ctx context.Context,
		id uuid.UUID,
		pagination *util.Pagination,
		from *time.Time,
	) ([]IMonitorResult, *apperr.ServiceError)
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

func (s *MonitorResultsService) GetLatestMonitorResultByMonitorID(
	ctx context.Context,
	monitorID uuid.UUID,
) (IMonitorResult, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitorResults, "GetLatestMonitorResultByMonitorID")
	logger.Trace().Str("monitorID", monitorID.String()).Msg("Retrieving latest monitor result by monitor ID")

	result, err := NewMonitorResultDAO(s.db.Querier()).GetLatestMonitorResultByMonitorID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("monitorID", monitorID.String()).Msg("No monitor result found for given monitor ID")
			return nil, apperr.NewNotFoundError("no monitor result found: %w", err)
		}
		logger.Error().Err(err).Msg("Failed to get latest monitor result by monitor ID")
		return nil, apperr.NewInternalError("failed to get latest monitor result: %w", err)
	}

	logger.Debug().Str("monitorID", monitorID.String()).Msg("Latest monitor result retrieved successfully")
	return result, nil
}

func (s *MonitorResultsService) GetMonitorResultsByMonitorID(
	ctx context.Context,
	id uuid.UUID,
	pagination *util.Pagination,
	from *time.Time,
) ([]IMonitorResult, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitorResults, "GetMonitorResultsByMonitorID")
	logger.Trace().
		Str("monitorID", id.String()).
		Interface("pagination", pagination).
		Interface("from", from).
		Msg("Retrieving monitor results by monitor ID")

	results, err := NewMonitorResultDAO(s.db.Querier()).GetMonitorResultsByMonitorID(ctx, id, pagination, from)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Error().Str("monitorID", id.String()).Msg("No monitor results found for given monitor ID")
			return nil, apperr.NewNotFoundError("no monitor results found: %w", err)
		}
		logger.Error().Err(err).Msg("Failed to get monitor results by monitor ID")
		return nil, apperr.NewInternalError("failed to get monitor results: %w", err)
	}

	logger.Debug().Str("monitorID", id.String()).Int("resultCount", len(results)).Msg("Monitor results retrieved successfully")
	return results, nil
}
