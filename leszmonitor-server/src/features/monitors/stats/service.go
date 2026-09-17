package stats

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

type IMonitorStatsService interface {
	GetStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (MonitorStats, *apperr.ServiceError)
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

func (s *MonitorStatsService) GetStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (MonitorStats, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStats, "GetStatsByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID).
		Time("from", from).
		Time("to", to).
		Msg("Getting stats by monitor ID")

	statsDAO := NewMonitorStatsDAO(s.db.Querier())

	latencyStats, err := statsDAO.GetLatencyStatsByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		if !errors.Is(err, db.ErrNotFound) {
			logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to get stats")
			return MonitorStats{}, &apperr.ServiceError{
				Code: http.StatusInternalServerError,
				Err:  errors.New("failed to get stats: " + err.Error()),
			}
		}
		logger.Warn().Str("monitorID", monitorID).Msg("No latency data found for the given monitor ID and time range")
	}

	hasNoStatusChanges := false
	var statusChangeStats StatusChangeStats
	statusChangeStats, err = statsDAO.GetStatusChangeStatsByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			logger.Warn().Str("monitorID", monitorID).Msg("No status change data found for the given monitor ID and time range")
			hasNoStatusChanges = true
		} else {
			return MonitorStats{}, &apperr.ServiceError{
				Code: http.StatusInternalServerError,
				Err:  errors.New("failed to get status change stats: " + err.Error()),
			}
		}
	}

	if hasNoStatusChanges {
		oldestResult, err := NewMonitorResultDAO(s.db.Querier()).GetOldestMonitorResultByMonitorID(ctx, monitorID)
		if err != nil {
			if errors.Is(err, db.ErrNotFound) {
				logger.Warn().Str("monitorID", monitorID).Msg("No monitor results found for the given monitor ID")
				return MonitorStats{
					Latency:      latencyStats,
					StatusChange: StatusChangeStats{},
					Uptime:       UptimeStats{},
				}, nil
			}
			logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to get oldest monitor result")
			return MonitorStats{}, &apperr.ServiceError{
				Code: http.StatusInternalServerError,
				Err:  errors.New("failed to get oldest monitor result: " + err.Error()),
			}
		}
		createdAt, err := time.Parse(time.RFC3339, oldestResult.GetCreatedAt())
		if err != nil {
			logger.Error().Err(err).Str("monitorID", monitorID).Msg("Failed to parse created_at of oldest monitor result")
			return MonitorStats{}, &apperr.ServiceError{
				Code: http.StatusInternalServerError,
				Err:  errors.New("failed to parse created_at of oldest monitor result: " + err.Error()),
			}
		}
		secondsInCurrentStatus := time.Since(createdAt).Seconds()
		statusChangeStats = StatusChangeStats{
			SecondsInCurrentStatus: int64(secondsInCurrentStatus),
		}
	}

	return MonitorStats{
		Latency:      latencyStats,
		StatusChange: statusChangeStats,
		Uptime:       UptimeStats{},
	}, nil
}
