package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	consts "github.com/m-milek/leszmonitor/models/consts"
	"github.com/m-milek/leszmonitor/models/monitorresult"
	"github.com/m-milek/leszmonitor/util"
)

type IMonitorResultRepository interface {
	InsertMonitorResult(ctx context.Context, result monitorresult.IMonitorResult) (interface{}, error)
	GetLatestMonitorResultByMonitorID(ctx context.Context, monitorID string) (monitorresult.IMonitorResult, error)
	GetMonitorResultsByMonitorID(ctx context.Context, id string, pagination *util.Pagination) ([]monitorresult.IMonitorResult, error)
	DeleteMonitorResultsOlderThanDuration(ctx context.Context, monitorID uuid.UUID, duration time.Duration) (int64, error)
}

type monitorResultRepository struct {
	baseRepository
}

func (r *monitorResultRepository) GetMonitorResultsByMonitorID(ctx context.Context, id string, pagination *util.Pagination) ([]monitorresult.IMonitorResult, error) {
	return dbWrap(ctx, "GetMonitorResultsByMonitorID", func() ([]monitorresult.IMonitorResult, error) {
		var results []monitorresult.MonitorResult

		err := sqlx.SelectContext(ctx, r.pool, &results, `
			SELECT mr.id, mr.monitor_id, m.kind, mr.is_success, mr.is_manually_triggered, mr.duration_ms, mr.error_details, mr.details, mr.created_at
			FROM monitor_results mr
			JOIN monitors m ON m.id = mr.monitor_id
			WHERE mr.monitor_id = $1
			ORDER BY mr.created_at DESC
			LIMIT $2 OFFSET $3`, id, pagination.PerPage, pagination.Offset())

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, ErrNotFound
			}
			return nil, err
		}

		var monitorResults []monitorresult.IMonitorResult
		for _, r := range results {
			err = processResultDetails(r)
			if err != nil {
				return nil, err
			}

			monitorResults = append(monitorResults, &r)
		}

		return monitorResults, nil
	})
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
			result.GetID(),
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

		err := sqlx.GetContext(ctx, r.pool, &result, `
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

		err = processResultDetails(result)
		if err != nil {
			return nil, err
		}

		return &result, nil
	})
}

func (r *monitorResultRepository) DeleteMonitorResultsOlderThanDuration(ctx context.Context, monitorID uuid.UUID, duration time.Duration) (int64, error) {
	return dbWrap(ctx, "DeleteMonitorResultsOlderThanDuration", func() (int64, error) {
		cutoffTime := time.Now().Add(-duration).Format(time.RFC3339)
		result, err := r.pool.ExecContext(ctx,
			`DELETE FROM monitor_results WHERE monitor_id = $1 AND created_at < $2`,
			monitorID,
			cutoffTime,
		)
		if err != nil {
			return 0, err
		}
		rowsAffected, err := result.RowsAffected()
		return rowsAffected, err
	})
}

func processResultDetails(result monitorresult.MonitorResult) error {
	details, err := monitorresult.ParseResultDetails(consts.ProbeType(result.MonitorType), result.DetailsJSON)
	if err != nil {
		return err
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

	return nil
}
