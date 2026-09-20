package downtime

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/m-milek/leszmonitor/platform/db"
	"github.com/m-milek/leszmonitor/platform/log"
	"github.com/rs/zerolog"
)

const durationBetweenBeats = time.Duration(5) * time.Second

const downtimeThreshold = 3 * durationBetweenBeats

type HeartbeatWorker struct {
	db     db.DB
	logger zerolog.Logger
}

func NewHeartbeatWorker(database db.DB) *HeartbeatWorker {
	return &HeartbeatWorker{
		db: database,
	}
}

func (w *HeartbeatWorker) Run(ctx context.Context) {
	w.logger = log.FromContext(ctx).With().Str("component", "liveness_worker").Logger()

	w.logger.Info().Msg("Starting liveness worker...")

	if err := w.beat(ctx); err != nil {
		w.logger.Error().Err(err).Msg("Failed to record heartbeat")
	}

	ticker := time.NewTicker(durationBetweenBeats)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("Liveness worker shutting down...")
			return
		case <-ticker.C:
			if err := w.beat(ctx); err != nil {
				w.logger.Error().Err(err).Msg("Failed to record heartbeat")
			}
		}
	}
}

func (w *HeartbeatWorker) beat(ctx context.Context) error {
	now := time.Now().UTC()

	return w.db.WithTx(ctx, func(q db.Querier) error {
		appDowntimeDAO := NewAppDowntimeDAO(q)

		previous, err := appDowntimeDAO.GetHeartbeat(ctx)
		if err != nil && !errors.Is(err, db.ErrNotFound) {
			return err
		}

		if previous != nil && now.Sub(previous.BeatAt) > downtimeThreshold {
			downtime := &AppDowntime{
				ID:        uuid.New(),
				StartedAt: previous.BeatAt,
				EndedAt:   now,
			}

			if _, err := appDowntimeDAO.InsertDowntime(ctx, downtime); err != nil {
				return err
			}

			w.logger.Warn().
				Time("started_at", downtime.StartedAt).
				Time("ended_at", downtime.EndedAt).
				Msg("Detected application downtime")
		}

		_, err = appDowntimeDAO.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: now})
		return err
	})
}
