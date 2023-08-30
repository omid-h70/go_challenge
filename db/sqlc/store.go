package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store will provide all functions to execute db queries and transactions
type Store interface {
	Querier
	execTx(ctx context.Context, fn func(q *Queries) error) error
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParam) (CreateUserTxResult, error)
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

func (store *SqlStore) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	/*
		&sql.TxOptions{
			Isolation: sql.LevelWriteCommitted,
			ReadOnly:  false,
		})
	*/

	if err != nil {
		return err
	}

	q := New(tx)
	execErr := fn(q)
	if execErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error:%v rb error:%v", rbErr, execErr)
		}
	}
	return tx.Commit()
}
