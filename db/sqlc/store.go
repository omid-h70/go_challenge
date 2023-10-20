package db

import (
	"context"
	"database/sql"
)

// Store will provide all functions to execute db queries and transactions
type Store interface {
	Querier
	//execTx(ctx context.Context, fn func(q *Queries) error) error
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParam) (VerifyEmailTxResult, error)
}

type SqlStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) Store {
	return &SqlStore{
		db:      db,
		Queries: New(db),
	}
}
