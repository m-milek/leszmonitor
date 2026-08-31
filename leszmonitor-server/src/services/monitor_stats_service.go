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
	GetAverageLatencyByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (float64, *ServiceError)
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

func (s *MonitorStatsService) GetAverageLatencyByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (float64, *ServiceError) {
	logger := MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStats, "GetAverageLatencyByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID).
		Time("from", from).
		Time("to", to).
		Msg("Getting average latency by monitor ID")

	avgLatency, err := s.db.MonitorStats().GetAverageLatencyByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("monitorID", monitorID).Msg("No latency data found for the given monitor ID and time range")
			return 0, &ServiceError{
				Code: http.StatusNotFound,
				Err:  errors.New("no latency data found for the given monitor ID and time range"),
			}
		}
		logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to get average latency")
		return 0, &ServiceError{
			Code: http.StatusInternalServerError,
			Err:  errors.New("failed to get average latency: " + err.Error()),
		}
	}
	logger.Info().Float64("avgLatency", avgLatency).Msg("Average latency calculated successfully")

	return avgLatency, nil
}
