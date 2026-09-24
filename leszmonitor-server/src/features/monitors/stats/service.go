package stats

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/statuschange"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
)

type IMonitorStatsService interface {
	GetStatsByMonitorID(ctx context.Context, monitorID uuid.UUID, from time.Time, to time.Time) (MonitorStats, *apperr.ServiceError)
}

type MonitorStatsService struct {
	db                db.DB
	statusChangeDAO   statuschange.IMonitorStatusChangeDAO
	monitorDAO        monitors.IMonitorDAO
	monitorResultsDAO results.IMonitorResultDAO
}

type MonitorStatsServiceDeps struct {
	DB db.DB
}

func NewMonitorStatsService(deps MonitorStatsServiceDeps) MonitorStatsService {
	statusChangeDAO := statuschange.NewMonitorStatusChangeDAO(deps.DB.Querier())
	monitorDAO := monitors.NewMonitorDAO(deps.DB.Querier())
	monitorResultsDAO := results.NewMonitorResultDAO(deps.DB.Querier())
	return MonitorStatsService{
		db:                deps.DB,
		statusChangeDAO:   statusChangeDAO,
		monitorDAO:        monitorDAO,
		monitorResultsDAO: monitorResultsDAO,
	}
}

func (s *MonitorStatsService) GetStatsByMonitorID(ctx context.Context, monitorID uuid.UUID, from time.Time, to time.Time) (MonitorStats, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStats, "GetStatsByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID.String()).
		Time("from", from).
		Time("to", to).
		Msg("Getting stats by monitor ID")

	monitor, err := s.monitorDAO.GetMonitorByID(ctx, monitorID)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get monitor by ID")
		return MonitorStats{}, apperr.NewInternalError("failed to get monitor: %w", err)
	}

	monitorResults, err := s.monitorResultsDAO.GetMonitorResultsByMonitorIDInTimeWindow(ctx, monitorID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get monitor results by monitor ID in time window")
		return MonitorStats{}, apperr.NewInternalError("failed to get monitor results: %w", err)
	}
	latestStatusChange, err := s.statusChangeDAO.GetLatestStatusChangeByMonitorID(ctx, monitorID, to)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		logger.Error().Err(err).Msg("Failed to get latest status change by monitor ID")
		return MonitorStats{}, apperr.NewInternalError("failed to get status changes: %w", err)
	}

	latencyStats := calculateLatencyStats(monitorResults)
	statusChangeStats := calculateStatusChangeStats(latestStatusChange, to)
	uptimeStats := calculateUptimeStats(monitorResults)
	probeSpecificStats, err := getProbeSpecificStats(monitor.Type, monitorResults)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get probe specific stats")
		return MonitorStats{}, apperr.NewInternalError("failed to get probe specific stats: %w", err)
	}

	return MonitorStats{
		ProbeType:     monitor.Type,
		Latency:       latencyStats,
		StatusChange:  statusChangeStats,
		Uptime:        uptimeStats,
		ProbeSpecific: probeSpecificStats,
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

func calculateStatusChangeStats(latestStatusChange *statuschange.MonitorStatusChange, to time.Time) StatusChangeStats {
	var secondsInCurrentStatus int64
	if latestStatusChange != nil {
		secondsInCurrentStatus = int64(to.Sub(latestStatusChange.CreatedAt).Seconds())
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

func getProbeSpecificStats(probeType kind.ProbeType, monitorResults []results.IMonitorResult) (any, error) {
	return mapToProbeStats(probeType, monitorResults)
}
