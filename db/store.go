package db

import (
	"context"
	"database/sql"
)

// Store defines all functions to execute db queries and transactions
type Store interface {
	Querier
	CreatePromotionTx(ctx context.Context, arg CreatePromotionTxParams) (CreatePromotionTxResult, error)
}

// SQLStore provides all functions to execute SQL queries and transactions
type SQLStore struct {
	db *sql.DB
	*Queries
}

// NewStore creates a new store
func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}
