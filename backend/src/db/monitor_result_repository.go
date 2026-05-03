package db

import (
	"context"
	"encoding/json"

	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/models/monitors"
)

type IMonitorResultRepository interface {
	InsertMonitorResult(ctx context.Context, result *monitorresult.IMonitorResult) (interface{}, error)
	GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (*monitors.IMonitorResult, error)
}

type monitorResultRepository struct {
	baseRepository
}

func newMonitorResultRepository(repository baseRepository) IMonitorResultRepository {
	return &monitorResultRepository{
		baseRepository: repository,
	}
}

func (r *monitorResultRepository) InsertMonitorResult(ctx context.Context, result *monitorresult.IMonitorResult) (interface{}, error) {
	return dbWrap(ctx, "InsertMonitorResult", func() (interface{}, error) {
		res := *result

		detailsJson, err := json.Marshal(res.GetDetails())
		if err != nil {
			return nil, err
		}
		_, err = r.pool.ExecContext(ctx,
			`INSERT INTO monitor_results (monitor_id, is_success, is_manually_triggered, duration_ms, error_message, details, created_at) 
            VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			res.GetMonitorID(),
			res.GetIsSuccess(),
			res.GetIsManuallyTriggered(),
			res.GetDurationMs(),
			res.GetErrorMessage(),
			detailsJson,
			res.GetCreatedAt(),
		)

		return nil, err
	})
}

func (r *monitorResultRepository) GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (*monitors.IMonitorResult, error) {
	return dbWrap(ctx, "GetLatestMonitorResultByMonitorID", func() (*monitors.IMonitorResult, error) {
		return nil, nil
	})
}
