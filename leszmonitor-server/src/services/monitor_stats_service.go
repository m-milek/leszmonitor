package services

import (
	"context"
	"errors"
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
		return 0, &ServiceError{
			Code: 500,
			Err:  errors.New("failed to get average latency: " + err.Error()),
		}
	}

	return avgLatency, nil
}
