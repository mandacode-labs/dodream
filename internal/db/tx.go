// Package db provides database utilities.
package db

import (
	"context"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"github.com/mandacode-labs/dodream/ent"
)

// WithTx runs fn inside a transaction.
// If fn returns an error, the transaction is rolled back.
// If fn succeeds, the transaction is committed.
func WithTx(ctx context.Context, client *ent.Client, fn func(*ent.Client) error) error {
	tx, err := client.Tx(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx.Client()); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("rollback transaction: %v (original error: %w)", rerr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}

// WithTxOptions runs fn inside a transaction with the given options.
// If fn returns an error, the transaction is rolled back.
// If fn succeeds, the transaction is committed.
func WithTxOptions(ctx context.Context, client *ent.Client, opts *sql.TxOptions, fn func(*ent.Client) error) error {
	tx, err := client.BeginTx(ctx, opts)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if v := recover(); v != nil {
			_ = tx.Rollback()
			panic(v)
		}
	}()

	if err := fn(tx.Client()); err != nil {
		if rerr := tx.Rollback(); rerr != nil {
			return fmt.Errorf("rollback transaction: %v (original error: %w)", rerr, err)
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
