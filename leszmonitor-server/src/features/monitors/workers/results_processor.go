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
	previousResult, err := results.NewMonitorResultDAO(database.Querier()).
		GetLatestMonitorResultByMonitorID(ctx, msg.Monitor.ID)
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			previousResult = nil
		} else {
			return errors.Wrap(err, "failed to retrieve previous monitor result")
		}
	}

	_, err = results.NewMonitorResultDAO(database.Querier()).InsertMonitorResult(ctx, msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to insert monitor result")
	}
	monitors.MonitorResultSavedChannel.Broadcast(msg)

	previousStatus := kind.MonitorStatusUnknown
	if previousResult != nil {
		previousStatus = previousResult.GetStatus()
	}

	if previousStatus != msg.Result.GetStatus() {
		err = handleStatusChange(ctx, database, msg.Monitor, previousStatus, msg.Result)
		if err != nil {
			return errors.Wrap(err, "failed to handle status change")
		}
	}
	return nil
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
		CausedByID:     current.GetID(),
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
