package db

import (
	"context"
	"database/sql"
	"fmt"
)

// TxOption defines a function type for applying options
type TxOption func(*sql.TxOptions)

// WithIsolationLevel is an option to set the transaction's isolation level
func WithIsolationLevel(level sql.IsolationLevel) TxOption {
	return func(opts *sql.TxOptions) {
		opts.Isolation = level
	}
}

// ExecTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error, opts ...TxOption) error {
	txOptions := &sql.TxOptions{}

	for _, opt := range opts {
		opt(txOptions)
	}

	tx, err := store.db.BeginTx(ctx, txOptions)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
