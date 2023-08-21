package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store will provide all functions to execute db queries and transactions
type Store interface {
	execTx(ctx context.Context, fn func(q *Queries) error) error
	TransferTx(ctx context.Context, arg TransferTxParam) (TransferTxResult, error)
}

type SqlStore struct {
	db *sql.DB
	*Queries
}

func NewStore(db *sql.DB) *SqlStore {
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

		//------------------------------ From Account
		createEntryArg := CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		}
		result.FromEntry, err = q.CreateEntry(ctx, createEntryArg)
		if err != nil {
			return err
		}

		/* V1 DeadLock and too many Operations
		fromAccount, err := q.GetAccount(ctx, arg.FromAccountID)
		if err != nil {
			return err
		}

		result.FromAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			AccountID: fromAccount.AccountID,
			Balance:   fromAccount.Balance - arg.Amount,
		})
		*/

		/* V2 DeadLock and too many Operations
		result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		*/

		if err != nil {
			return err
		}

		//----------------------------- To Account
		toEntryArg := CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		}
		result.ToEntry, err = q.CreateEntry(ctx, toEntryArg)
		if err != nil {
			return err
		}

		/* V1 DeadLock and too many Operations
		toAccount, err := q.GetAccount(ctx, arg.ToAccountID)
		if err != nil {
			return err
		}

		result.ToAccount, err = q.UpdateAccount(ctx, UpdateAccountParams{
			AccountID: toAccount.AccountID,
			Balance:   toAccount.Balance + arg.Amount,
		})
		*/

		/* V2 DeadLock and too many Operations
		result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		*/

		/*Order is important to avoid DeadLock*/
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return err

	})
	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		AccountID: accountID1,
		Amount:    amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		AccountID: accountID2,
		Amount:    amount2,
	})
	return
}
