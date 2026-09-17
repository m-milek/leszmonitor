package monitors

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/m-milek/leszmonitor/platform/db"
)

type IMonitorStatsDAO interface {
	GetLatencyStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (LatencyStats, error)
	GetStatusChangeStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (StatusChangeStats, error)
}

type monitorStatsDAO struct {
	pool db.Querier
}

func NewMonitorStatsDAO(pool db.Querier) IMonitorStatsDAO {
	return &monitorStatsDAO{
		pool: pool,
	}
}

func (m *monitorStatsDAO) GetLatencyStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (LatencyStats, error) {
	query := `
		SELECT AVG(duration_ms), MIN(duration_ms), MAX(duration_ms)
		FROM monitor_results
		WHERE monitor_id = $1
		  AND created_at >= $2
		  AND created_at < $3
	`
	var avgLatency, minLatency, maxLatency sql.NullFloat64
	row := m.pool.QueryRowxContext(ctx, query, monitorID, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
	if err := row.Scan(&avgLatency, &minLatency, &maxLatency); err != nil {
		return LatencyStats{}, err
	}
	if !avgLatency.Valid {
		return LatencyStats{}, db.ErrNotFound
	}
	return LatencyStats{
		Avg: avgLatency.Float64,
		Min: minLatency.Float64,
		Max: maxLatency.Float64,
	}, nil
}

func (m *monitorStatsDAO) GetStatusChangeStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (StatusChangeStats, error) {
	var lastCreatedAt string
	query := `
		SELECT created_at
		FROM monitor_status_changes
		WHERE monitor_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	row := m.pool.QueryRowxContext(ctx, query, monitorID)
	if err := row.Scan(&lastCreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return StatusChangeStats{}, db.ErrNotFound
		}
		return StatusChangeStats{}, err
	}

	secondsInCurrentStatus, err := time.Parse(time.RFC3339, lastCreatedAt)
	if err != nil {
		return StatusChangeStats{}, err
	}

	duration := time.Since(secondsInCurrentStatus)

	return StatusChangeStats{
		SecondsInCurrentStatus: int64(duration.Seconds()),
	}, nil
}
