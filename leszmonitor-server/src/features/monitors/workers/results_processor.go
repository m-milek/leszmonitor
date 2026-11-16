package workers

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/features/monitors"
	"github.com/m-milek/leszmonitor/features/monitors/kind"
	"github.com/m-milek/leszmonitor/features/monitors/results"
	"github.com/m-milek/leszmonitor/features/monitors/statuschange"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/m-milek/leszmonitor/platform/util"
	"github.com/pkg/errors"
)

type ResultsProcessor struct {
	db db.DB
}

func NewResultsProcessor(database db.DB) *ResultsProcessor {
	return &ResultsProcessor{
		db: database,
	}
}

func (p *ResultsProcessor) Run(ctx context.Context) {
	logger := log.FromContext(ctx).With().Str("component", "results_processor").Logger()

	logger.Info().Msg("Starting results processor...")

	results := monitors.MonitorRunChannel.Subscribe()
	defer monitors.MonitorRunChannel.Unsubscribe(results)

	for {
		select {
		case <-ctx.Done():
			logger.Info().Msg("Stopping results processor")
			return
		case msg := <-results:
			err := processMonitorRunMessage(ctx, p.db, msg)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to process monitor run result")
			}
		}
	}
}

func processMonitorRunMessage(ctx context.Context, database db.DB, msg monitors.MonitorRunMessage) error {
	_, err := results.NewMonitorResultDAO(database.Querier()).InsertMonitorResult(ctx, msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to insert monitor result")
	}
	monitors.MonitorResultSavedChannel.Broadcast(msg)

	latestStatusChange, err := statuschange.NewMonitorStatusChangeDAO(database.Querier()).
		GetLatestStatusChangeByMonitorID(ctx, msg.Monitor.ID, time.Now().UTC())
	if err != nil && !errors.Is(err, db.ErrNotFound) {
		return errors.Wrap(err, "failed to retrieve latest status change")
	}

	previousStatus := kind.MonitorStatusUnknown
	if latestStatusChange != nil {
		previousStatus = kind.MonitorStatus(latestStatusChange.NextStatus)
	}

	if previousStatus == msg.Result.GetStatus() {
		return nil
	}

	if previousStatus != kind.MonitorStatusUnknown {
		confirmed, err := isStatusConfirmed(ctx, database, msg.Monitor, msg.Result.GetStatus())
		if err != nil {
			return errors.Wrap(err, "failed to check status threshold")
		}
		if !confirmed {
			return nil
		}
	}

	err = handleStatusChange(ctx, database, msg.Monitor, previousStatus, msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to handle status change")
	}
	return nil
}

func isStatusConfirmed(
	ctx context.Context,
	database db.DB,
	monitor monitors.Monitor,
	status kind.MonitorStatus,
) (bool, error) {
	threshold := monitor.RecoveryThreshold
	if status == kind.MonitorStatusDown {
		threshold = monitor.FailureThreshold
	}
	if threshold <= 1 {
		return true, nil
	}

	latest, err := results.NewMonitorResultDAO(database.Querier()).
		GetMonitorResultsByMonitorID(ctx, monitor.ID, &util.Pagination{Page: 1, PerPage: threshold}, nil)
	if err != nil {
		return false, err
	}
	if len(latest) < threshold {
		return false, nil
	}
	for _, result := range latest {
		if result.GetStatus() != status {
			return false, nil
		}
	}
	return true, nil
}

func handleStatusChange(
	ctx context.Context,
	database db.DB,
	monitor monitors.Monitor,
	previousStatus kind.MonitorStatus,
	current results.IMonitorResult,
) error {
	logger := log.FromContext(ctx)
	monitorStatusChange := statuschange.MonitorStatusChange{
		ID:             uuid.New(),
		MonitorID:      monitor.ID,
		CausedByID:     new(current.GetID()),
		PreviousStatus: string(previousStatus),
		NextStatus:     string(current.GetStatus()),
		CreatedAt:      time.Now().UTC(),
	}
	_, err := statuschange.NewMonitorStatusChangeDAO(database.Querier()).InsertStatusChange(ctx, monitorStatusChange)
	if err != nil {
		return errors.Wrap(err, "failed to insert monitor status change")
	}
	logger.Debug().
		Str("monitor_id", monitor.ID.String()).
		Str("previous_status", string(previousStatus)).
		Str("next_status", string(current.GetStatus())).
		Msg("Monitor status change recorded")
	return nil
}
