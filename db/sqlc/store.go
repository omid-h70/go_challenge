package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store will provide all functions to execute db queries and transactions
type Store interface {
	execTx(ctx context.Context, fn func(q *Queries) error) error
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

// TransferTxParam Contains input parameter for Transfer Transaction
type TransferTxParam struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs money transfer from one account to the other
// it creates a transfer record, add account entries, update accounts balance within a single transaction
func (store *SqlStore) TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		createTransferArg := CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		}
		var err error
		result.Transfer, err = q.CreateTransfer(ctx, createTransferArg)
		if err != nil {
			return err
		}

		createEntryArg := CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}
		result.FromEntry, err = q.CreateEntry(ctx, createEntryArg)
		if err != nil {
			return err
		}

		toEntryArg := CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}
		result.ToEntry, err = q.CreateEntry(ctx, toEntryArg)
		if err != nil {
			return err
		}

		return err
	})
	return result, err
}
