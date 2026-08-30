package resultsprocessor

import (
	"context"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/db"
	"github.com/m-milek/leszmonitor/events"
	"github.com/m-milek/leszmonitor/log"
	"github.com/m-milek/leszmonitor/models"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/models/monitors"
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

	results := events.MonitorRunChannel.Subscribe()
	defer events.MonitorRunChannel.Unsubscribe(results)

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
	previousResult, err := database.MonitorResults().GetLatestMonitorResultByMonitorID(ctx, msg.Monitor.ID.String())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			previousResult = nil
		} else {
			return errors.Wrap(err, "failed to retrieve previous monitor result")
		}
	}

	_, err = database.MonitorResults().InsertMonitorResult(ctx, msg.Result)
	if err != nil {
		return errors.Wrap(err, "failed to insert monitor result")
	}
	if isStatusChange(previousResult, msg.Result) {
		err = handleStatusChange(ctx, database, msg.Monitor, previousResult, msg.Result)
		if err != nil {
			return errors.Wrap(err, "failed to handle status change")
		}
	}
	return nil
}

func isStatusChange(previous monitorresult.IMonitorResult, current monitorresult.IMonitorResult) bool {
	if previous == nil {
		return false
	}
	return previous.GetStatus() != current.GetStatus()
}

func handleStatusChange(ctx context.Context, db db.DB, monitor monitors.Monitor, previous monitorresult.IMonitorResult, current monitorresult.IMonitorResult) error {
	logger := log.FromContext(ctx)
	monitorStatusChange := models.MonitorStatusChange{
		ID:             uuid.New(),
		MonitorID:      monitor.ID,
		CausedByID:     current.GetID(),
		PreviousStatus: string(previous.GetStatus()),
		NextStatus:     string(current.GetStatus()),
	}
	_, err := db.MonitorStatusChanges().InsertStatusChange(ctx, monitorStatusChange)
	if err != nil {
		return errors.Wrap(err, "failed to insert monitor status change")
	}
	logger.Debug().
		Str("monitor_id", monitor.ID.String()).
		Str("previous_status", string(previous.GetStatus())).
		Str("next_status", string(current.GetStatus())).
		Msg("Monitor status change recorded")
	return nil
}
