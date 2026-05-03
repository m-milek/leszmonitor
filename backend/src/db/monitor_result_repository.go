package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	consts "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
)

type IMonitorResultRepository interface {
	InsertMonitorResult(ctx context.Context, result monitorresult.IMonitorResult) (interface{}, error)
	GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (monitorresult.IMonitorResult, error)
}

type monitorResultRepository struct {
	baseRepository
}

func newMonitorResultRepository(repository baseRepository) IMonitorResultRepository {
	return &monitorResultRepository{
		baseRepository: repository,
	}
}

func (r *monitorResultRepository) InsertMonitorResult(ctx context.Context, result monitorresult.IMonitorResult) (interface{}, error) {
	return dbWrap(ctx, "InsertMonitorResult", func() (interface{}, error) {
		detailsJson, err := json.Marshal(result.GetDetails())
		if err != nil {
			return nil, err
		}

		var errorDetailsJson []byte
		if ed := result.GetErrorDetails(); ed.ErrorMessage != "" || len(ed.Errors) > 0 || len(ed.Failures) > 0 {
			var err error
			errorDetailsJson, err = json.Marshal(ed)
			if err != nil {
				return nil, err
			}
		}

		_, err = r.pool.ExecContext(ctx,
			`INSERT INTO monitor_results (id, monitor_id, is_success, is_manually_triggered, duration_ms, error_details, details, created_at) 
            VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			uuid.New().String(),
			result.GetMonitorID(),
			result.GetIsSuccess(),
			result.GetIsManuallyTriggered(),
			result.GetDurationMs(),
			errorDetailsJson,
			detailsJson,
			result.GetCreatedAt(),
		)

		return nil, err
	})
}

func (r *monitorResultRepository) GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (monitorresult.IMonitorResult, error) {
	return dbWrap(ctx, "GetLatestMonitorResultByMonitorID", func() (monitorresult.IMonitorResult, error) {
		var result monitorresult.MonitorResult

		err := r.pool.GetContext(ctx, &result, `
            SELECT mr.monitor_id, m.kind, mr.is_success, mr.is_manually_triggered, mr.duration_ms, mr.error_details, mr.details, mr.created_at
            FROM monitor_results mr
            JOIN monitors m ON m.id = mr.monitor_id
            WHERE mr.monitor_id = $1
            ORDER BY mr.created_at DESC LIMIT 1`, monitorID)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		details, err := monitorresult.ParseResultDetails(consts.MonitorConfigType(result.MonitorType), result.DetailsJSON)
		if err != nil {
			return nil, err
		}
		result.Details = details

		if len(result.ErrorDetailsJSON) > 0 {
			var errorDetails monitorresult.ErrorDetails
			if err := json.Unmarshal(result.ErrorDetailsJSON, &errorDetails); err == nil {
				if errorDetails.ErrorMessage != "" || len(errorDetails.Errors) > 0 || len(errorDetails.Failures) > 0 {
					result.ErrorDetails = &errorDetails
				}
			}
		}

		return &result, nil
	})
}
