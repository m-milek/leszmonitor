package stats

import (
	"context"
	"time"

	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/statuschange"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

type IMonitorStatsService interface {
	GetStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (MonitorStats, *apperr.ServiceError)
}

type MonitorStatsService struct {
	db                db.DB
	statusChangeDAO   statuschange.IMonitorStatusChangeDAO
	monitorResultsDAO results.IMonitorResultDAO
}

type MonitorStatsServiceDeps struct {
	DB db.DB
}

func NewMonitorStatsService(deps MonitorStatsServiceDeps) MonitorStatsService {
	statusChangeDAO := statuschange.NewMonitorStatusChangeDAO(deps.DB.Querier())
	monitorResultsDAO := results.NewMonitorResultDAO(deps.DB.Querier())
	return MonitorStatsService{
		db:                deps.DB,
		statusChangeDAO:   statusChangeDAO,
		monitorResultsDAO: monitorResultsDAO,
	}
}

func (s *MonitorStatsService) GetStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (MonitorStats, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStats, "GetStatsByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID).
		Time("from", from).
		Time("to", to).
		Msg("Getting stats by monitor ID")

	monitorResults, err := s.monitorResultsDAO.GetMonitorResultsByMonitorIDInTimeWindow(ctx, monitorID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get monitor results by monitor ID in time window")
		return MonitorStats{}, apperr.NewInternalError("failed to get monitor results: %w", err)
	}
	statusChanges, err := s.statusChangeDAO.GetStatusChangesByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get status changes by monitor ID in time window")
		return MonitorStats{}, apperr.NewInternalError("failed to get status changes: %w", err)
	}

	latencyStats := calculateLatencyStats(monitorResults)
	statusChangeStats := calculateStatusChangeStats(statusChanges)
	uptimeStats := calculateUptimeStats(monitorResults)

	return MonitorStats{
		Latency:      latencyStats,
		StatusChange: statusChangeStats,
		Uptime:       uptimeStats,
	}, nil
}

func calculateLatencyStats(monitorResults []results.IMonitorResult) LatencyStats {
	var minLatency, maxLatency, totalLatency float64
	for _, result := range monitorResults {
		duration := float64(result.GetDurationMs())
		if minLatency == 0 || duration < minLatency {
			minLatency = duration
		}
		if duration > maxLatency {
			maxLatency = duration
		}
		totalLatency += duration
	}

	var avgLatency float64
	if len(monitorResults) > 0 {
		avgLatency = totalLatency / float64(len(monitorResults))
	}

	return LatencyStats{
		Min: minLatency,
		Max: maxLatency,
		Avg: avgLatency,
	}
}

func calculateStatusChangeStats(statusChanges []statuschange.MonitorStatusChange) StatusChangeStats {
	var secondsInCurrentStatus int64
	if len(statusChanges) > 0 {
		latestStatusChange := statusChanges[len(statusChanges)-1]
		secondsInCurrentStatus = int64(time.Since(latestStatusChange.CreatedAt).Seconds())
	}

	return StatusChangeStats{
		SecondsInCurrentStatus: secondsInCurrentStatus,
	}
}

func calculateUptimeStats(monitorResults []results.IMonitorResult) UptimeStats {
	if len(monitorResults) == 0 {
		return UptimeStats{}
	}

	statusToCount := make(map[kind.MonitorStatus]int)
	for _, result := range monitorResults {
		statusToCount[result.GetStatus()]++
	}

	total := float64(len(monitorResults))
	statusToPercentage := make(map[kind.MonitorStatus]float64, len(statusToCount))
	for status, count := range statusToCount {
		statusToPercentage[status] = float64(count) / total * 100
	}

	return UptimeStats{
		StatusToCount:      statusToCount,
		StatusToPercentage: statusToPercentage,
	}
}
