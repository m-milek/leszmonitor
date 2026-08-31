package services

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/m-milek/leszmonitor/constants"
	"github.com/m-milek/leszmonitor/db"
)

type IMonitorStatsService interface {
	GetLatencyStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (db.LatencyStats, *ServiceError)
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

func (s *MonitorStatsService) GetLatencyStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (db.LatencyStats, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStats, "GetLatencyStatsByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID).
		Time("from", from).
		Time("to", to).
		Msg("Getting latency stats by monitor ID")

	stats, err := s.db.MonitorStats().GetLatencyStatsByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("monitorID", monitorID).Msg("No latency data found for the given monitor ID and time range")
			return db.LatencyStats{}, &ServiceError{
				Code: http.StatusNotFound,
				Err:  errors.New("no latency data found for the given monitor ID and time range"),
			}
		}
		logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to get latency stats")
		return db.LatencyStats{}, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("failed to get latency stats: " + err.Error()),
		}
	}
	logger.Info().
		Float64("avg", stats.Avg).
		Float64("min", stats.Min).
		Float64("max", stats.Max).
		Msg("Latency stats calculated successfully")

	return stats, nil
}
