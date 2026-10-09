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

// DowntimeHandler is called in the same transaction that records a detected downtime.
type DowntimeHandler func(ctx context.Context, q db.Querier, downtime AppDowntime) error

type HeartbeatWorker struct {
	db         db.DB
	onDowntime DowntimeHandler
	logger     zerolog.Logger
}

func NewHeartbeatWorker(database db.DB, onDowntime DowntimeHandler) *HeartbeatWorker {
	return &HeartbeatWorker{
		db:         database,
		onDowntime: onDowntime,
	}
}

func (w *HeartbeatWorker) Run(ctx context.Context) {
	w.logger = log.FromContext(ctx).With().Str("component", "liveness_worker").Logger()

	ctx = log.WithContext(ctx, &w.logger)

	w.logger.Info().Msg("Starting liveness worker...")

	ticker := time.NewTicker(durationBetweenBeats)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			w.logger.Info().Msg("Liveness worker shutting down...")
			return
		case <-ticker.C:
			if err := w.Beat(ctx); err != nil {
				w.logger.Error().Err(err).Msg("Failed to record heartbeat")
			}
		}
	}
}

// Beat records a heartbeat and, if the previous one is too old, the downtime window between them.
func (w *HeartbeatWorker) Beat(ctx context.Context) error {
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

			if w.onDowntime != nil {
				if err := w.onDowntime(ctx, q, *downtime); err != nil {
					return err
				}
			}

			log.FromContext(ctx).Warn().
				Time("started_at", downtime.StartedAt).
				Time("ended_at", downtime.EndedAt).
				Msg("Detected application downtime")
		}

		_, err = appDowntimeDAO.InsertHeartbeat(ctx, &HeartbeatRecord{BeatAt: now})
		return err
	})
}
