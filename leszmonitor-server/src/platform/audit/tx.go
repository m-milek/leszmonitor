package audit

import (
	"context"
	"fmt"

	"github.com/m-milek/leszmonitor/platform/db"
)

// TxRunner runs a function inside a database transaction.
type TxRunner interface {
	WithTx(ctx context.Context, fn func(q db.Querier) error) error
}

// WithAuditedTx executes fn inside a database transaction and, if successful, records an audit log.
// It returns the result from fn and any error that occurred.
func WithAuditedTx[T any](
	ctx context.Context,
	client TxRunner,
	fn func(q db.Querier) (T, *AuditLogParams, error),
) (T, error) {
	var result T
	err := client.WithTx(ctx, func(q db.Querier) error {
		res, params, err := fn(q)
		if err != nil {
			return err
		}
		result = res
		if params != nil {
			if auditErr := NewAuditLogDAO(q).Record(ctx, *params); auditErr != nil {
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
	client TxRunner,
	fn func(q db.Querier) (*AuditLogParams, error),
) error {
	return client.WithTx(ctx, func(q db.Querier) error {
		params, err := fn(q)
		if err != nil {
			return err
		}
		if params != nil {
			if auditErr := NewAuditLogDAO(q).Record(ctx, *params); auditErr != nil {
				return fmt.Errorf("failed to record audit log: %w", auditErr)
			}
		}
		return nil
	})
}
