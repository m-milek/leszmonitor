package db

import (
	"context"
	"database/sql"
	"time"
)

type LatencyStats struct {
	Avg float64
	Min float64
	Max float64
}

type IMonitorStatsDAO interface {
	GetLatencyStatsByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (LatencyStats, error)
}

type monitorStatsDAO struct {
	baseDAO
}

func newMonitorStatsDAO(base baseDAO) IMonitorStatsDAO {
	return &monitorStatsDAO{
		baseDAO: base,
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
	var avg, min, max sql.NullFloat64
	row := m.pool.QueryRowxContext(ctx, query, monitorID, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
	if err := row.Scan(&avg, &min, &max); err != nil {
		return LatencyStats{}, err
	}
	if !avg.Valid {
		return LatencyStats{}, ErrNotFound
	}
	return LatencyStats{
		Avg: avg.Float64,
		Min: min.Float64,
		Max: max.Float64,
	}, nil
}
