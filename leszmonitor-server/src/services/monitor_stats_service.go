package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/models"
)

type IMonitorStatsService interface {
	GetStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (models.MonitorStats, *ServiceError)
}

type MonitorStatsService struct {
	db db.DB
}

type MonitorStatsServiceDeps struct {
	DB db.DB
}

func NewMonitorStatsService(deps MonitorStatsServiceDeps) MonitorStatsService {
	return MonitorStatsService{
		db: deps.DB,
	}
}

func (s *MonitorStatsService) GetStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (models.MonitorStats, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStats, "GetStatsByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID).
		Time("from", from).
		Time("to", to).
		Msg("Getting stats by monitor ID")

	latencyStats, err := s.db.MonitorStats().GetLatencyStatsByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("monitorID", monitorID).Msg("No data found for the given monitor ID and time range")
			return models.MonitorStats{}, &ServiceError{
				Code: http.StatusNotFound,
				Err:  errors.New("no data found for the given monitor ID and time range"),
			}
		}
		logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to get stats")
		return models.MonitorStats{}, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("failed to get stats: " + err.Error()),
		}
	}

	statusChangeStats, err := s.db.MonitorStats().GetStatusChangesByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("monitorID", monitorID).Msg("No status change data found for the given monitor ID and time range")
			return models.MonitorStats{}, nil
		}
		logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to get status change stats")
		return models.MonitorStats{}, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("failed to get status change stats: " + err.Error()),
		}
	}

	return models.MonitorStats{
		Latency:      latencyStats,
		StatusChange: statusChangeStats,
		Uptime:       models.UptimeStats{},
	}, nil
}
