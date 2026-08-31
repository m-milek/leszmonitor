package db

import (
	"context"
	"database/sql"
	"time"
)

type IMonitorStatsDAO interface {
	GetAverageLatencyByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (float64, error)
}

type monitorStatsDAO struct {
	baseDAO
}

func newMonitorStatsDAO(base baseDAO) IMonitorStatsDAO {
	return &monitorStatsDAO{
		baseDAO: base,
	}
}

func (m *monitorStatsDAO) GetAverageLatencyByMonitorID(ctx context.Context, monitorID string, from time.Time, to time.Time) (float64, error) {
	query := `
		SELECT AVG(duration_ms)
		FROM monitor_results
		WHERE monitor_id = $1
		  AND created_at >= $2
		  AND created_at < $3
	`
	var avgLatency sql.NullFloat64
	row := m.pool.QueryRowxContext(ctx, query, monitorID, from.UTC().Format(time.RFC3339), to.UTC().Format(time.RFC3339))
	if err := row.Scan(&avgLatency); err != nil {
		return 0, err
	}
	if !avgLatency.Valid {
		return 0, ErrNotFound
	}
	return avgLatency.Float64, nil
}
