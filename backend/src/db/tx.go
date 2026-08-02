package db

import (
	"context"
	"fmt"

	"github.com/m-milek/leszmonitor/security"
)

// WithTx executes fn inside a database transaction. If fn returns nil,
// the transaction is committed. If fn returns an error or panics,
// the transaction is rolled back. Nested transactions are not supported.
func (c *Client) WithTx(ctx context.Context, fn func(tx DB) error) error {
	if c.sqlxDB == nil {
		return fmt.Errorf("nested transactions are not supported")
	}

	tx, err := c.sqlxDB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txClient := newClientFromPool(tx)

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(txClient); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}

// WithAuditedTx executes fn inside a database transaction and, if successful, records an audit log.
// It returns the result from fn and any error that occurred.
func WithAuditedTx[T any](
	ctx context.Context,
	client DB,
	fn func(tx DB) (T, *security.AuditLogParams, error),
) (T, error) {
	var result T
	err := client.WithTx(ctx, func(tx DB) error {
		res, params, err := fn(tx)
		if err != nil {
			return err
		}
		result = res
		if params != nil {
			if auditErr := tx.AuditLog().Record(ctx, *params); auditErr != nil {
				return fmt.Errorf("failed to record audit log: %w", auditErr)
			}
		}
		return nil
	})
	return result, err
}

// WithAuditedVoidTx executes fn inside a database transaction and, if successful, records an audit log.
// This is a convenience wrapper for operations that do not return a result.
func WithAuditedVoidTx(
	ctx context.Context,
	client DB,
	fn func(tx DB) (*security.AuditLogParams, error),
) error {
	return client.WithTx(ctx, func(tx DB) error {
		params, err := fn(tx)
		if err != nil {
			return err
		}
		if params != nil {
			if auditErr := tx.AuditLog().Record(ctx, *params); auditErr != nil {
				return fmt.Errorf("failed to record audit log: %w", auditErr)
			}
		}
		return nil
	})
}
