package globalparameters

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/m-milek/leszmonitor/platform/db"
)

type IGlobalParameterDAO interface {
	GetAllParameters(ctx context.Context) ([]GlobalParameterRecord, error)
	GetParameterByKey(ctx context.Context, key GlobalParameterKey) (*GlobalParameterRecord, error)
	UpsertParameter(ctx context.Context, record GlobalParameterRecord) (*GlobalParameterRecord, error)
}

type globalParameterDAO struct {
	pool db.Querier
}

func NewGlobalParameterDAO(pool db.Querier) IGlobalParameterDAO {
	return &globalParameterDAO{
		pool: pool,
	}
}

func (r *globalParameterDAO) GetAllParameters(ctx context.Context) ([]GlobalParameterRecord, error) {
	return db.Wrap(ctx, "GetAllGlobalParameters", func() ([]GlobalParameterRecord, error) {
		records := []GlobalParameterRecord{}

		err := sqlx.SelectContext(ctx, r.pool, &records, `
			SELECT key, value, updated_at
			FROM global_parameters
			ORDER BY key ASC`,
		)
		if err != nil {
			return nil, err
		}

		return records, nil
	})
}

func (r *globalParameterDAO) GetParameterByKey(ctx context.Context, key GlobalParameterKey) (*GlobalParameterRecord, error) {
	return db.Wrap(ctx, "GetGlobalParameterByKey", func() (*GlobalParameterRecord, error) {
		var record GlobalParameterRecord

		err := sqlx.GetContext(ctx, r.pool, &record, `
			SELECT key, value, updated_at
			FROM global_parameters
			WHERE key = $1`,
			key,
		)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		return &record, nil
	})
}

func (r *globalParameterDAO) UpsertParameter(ctx context.Context, record GlobalParameterRecord) (*GlobalParameterRecord, error) {
	return db.Wrap(ctx, "UpsertGlobalParameter", func() (*GlobalParameterRecord, error) {
		var saved GlobalParameterRecord

		err := r.pool.QueryRowxContext(ctx, `
			INSERT INTO global_parameters (key, value)
			VALUES ($1, $2)
			ON CONFLICT (key) DO UPDATE SET value = excluded.value, updated_at = CURRENT_TIMESTAMP
			RETURNING key, value, updated_at`,
			record.Key,
			record.Value,
		).StructScan(&saved)
		if err != nil {
			return nil, err
		}

		return &saved, nil
	})
}
