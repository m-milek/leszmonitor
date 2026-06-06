package db

import (
	"context"
	"fmt"
)

// WithTx executes fn inside a database transaction. If fn returns nil,
// the transaction is committed. If fn returns an error or panics,
// the transaction is rolled back. Nested transactions are not supported.
func (c *DBClient) WithTx(ctx context.Context, fn func(tx DB) error) error {
	if c.sqlxDB == nil {
		return fmt.Errorf("nested transactions are not supported")
	}

	tx, err := c.sqlxDB.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	txClient := newDBClientFromPool(tx)

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
