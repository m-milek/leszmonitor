package statuschange

import (
	"context"
	"errors"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/platform/apperr"
	"github.com/m-milek/leszmonitor/platform/constants"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
)

type IMonitorStatusChangeService interface {
	GetStatusHistoryByMonitorID(ctx context.Context, monitorID uuid.UUID, pagination *util.Pagination, from time.Time, to time.Time) ([]MonitorStatusPeriod, *apperr.ServiceError)
}

type MonitorStatusChangeService struct {
	db                db.DB
	statusChangeDAO   IMonitorStatusChangeDAO
	monitorDAO        monitors.IMonitorDAO
	monitorResultsDAO results.IMonitorResultDAO
}

type MonitorStatusChangeServiceDeps struct {
	DB db.DB
}

func NewMonitorStatusChangeService(deps MonitorStatusChangeServiceDeps) MonitorStatusChangeService {
	statusChangeDAO := NewMonitorStatusChangeDAO(deps.DB.Querier())
	monitorDAO := monitors.NewMonitorDAO(deps.DB.Querier())
	monitorResultsDAO := results.NewMonitorResultDAO(deps.DB.Querier())
	return MonitorStatusChangeService{
		db:                deps.DB,
		statusChangeDAO:   statusChangeDAO,
		monitorDAO:        monitorDAO,
		monitorResultsDAO: monitorResultsDAO,
	}
}

func (s *MonitorStatusChangeService) GetStatusHistoryByMonitorID(ctx context.Context, monitorID uuid.UUID, pagination *util.Pagination, from time.Time, to time.Time) ([]MonitorStatusPeriod, *apperr.ServiceError) {
	logger := log.MethodLoggerFromContext(ctx, constants.ServiceNameMonitorStatusChange, "GetStatusHistoryByMonitorID")
	logger.Trace().
		Str("monitorID", monitorID.String()).
		Interface("pagination", pagination).
		Time("from", from).
		Time("to", to).
		Msg("Getting status history by monitor ID")

	monitor, err := s.monitorDAO.GetMonitorByID(ctx, monitorID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, apperr.NewNotFoundError("monitor with ID %s not found", monitorID.String())
		}
		logger.Error().Err(err).Msg("Failed to get monitor by ID")
		return nil, apperr.NewInternalError("failed to get monitor: %w", err)
	}

	previousStatusChange, err := s.statusChangeDAO.GetLatestStatusChangeByMonitorID(ctx, monitorID, from)
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		logger.Error().Err(err).Msg("Failed to get latest status change by monitor ID")
		return nil, apperr.NewInternalError("failed to get latest status change: %w", err)
	}

	statusChanges, err := s.statusChangeDAO.GetStatusChangesByMonitorID(ctx, monitorID, from, to)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to get status changes by monitor ID")
		return nil, apperr.NewInternalError("failed to get status changes: %w", err)
	}

	history := buildStatusHistory(monitor.CreatedAt, previousStatusChange, statusChanges, to)
	slices.Reverse(history)

	start := min(pagination.Offset(), len(history))
	end := min(start+pagination.PerPage, len(history))

	return history[start:end], nil
}

// buildStatusHistory turns status changes into consecutive status periods, oldest first.
// Without a change before the window, the first period starts when the monitor was created.
func buildStatusHistory(
	monitorCreatedAt time.Time,
	previousStatusChange *MonitorStatusChange,
	statusChanges []MonitorStatusChange,
	to time.Time,
) []MonitorStatusPeriod {
	current := MonitorStatusPeriod{Status: kind.MonitorStatusUnknown, StartedAt: monitorCreatedAt}
	if previousStatusChange != nil {
		current = MonitorStatusPeriod{
			Status:    kind.MonitorStatus(previousStatusChange.NextStatus),
			StartedAt: previousStatusChange.CreatedAt,
		}
	} else if len(statusChanges) > 0 {
		current.Status = kind.MonitorStatus(statusChanges[0].PreviousStatus)
	}

	history := []MonitorStatusPeriod{}
	for _, statusChange := range statusChanges {
		if previousStatusChange != nil && statusChange.ID == previousStatusChange.ID {
			continue
		}
		current.EndedAt = &statusChange.CreatedAt
		current.DurationSeconds = int64(statusChange.CreatedAt.Sub(current.StartedAt).Seconds())
		history = append(history, current)

		current = MonitorStatusPeriod{
			Status:    kind.MonitorStatus(statusChange.NextStatus),
			StartedAt: statusChange.CreatedAt,
		}
	}
	current.DurationSeconds = int64(to.Sub(current.StartedAt).Seconds())

	return append(history, current)
}
