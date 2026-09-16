package db

import (
	"context"
	"fmt"
)

// WithTx executes fn inside a database transaction. If fn returns nil,
// the transaction is committed. If fn returns an error or panics,
// the transaction is rolled back.
func (c *Client) WithTx(ctx context.Context, fn func(q Querier) error) error {
	tx, err := c.sqlxDB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
