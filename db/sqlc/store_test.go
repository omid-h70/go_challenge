package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before Account 1", account1.Balance, ">> before Account 2", account2.Balance)

	n := 10
	amount := int64(10)

	errChannel := make(chan error)
	transferChannel := make(chan TransferTxResult)

	//run n concurrent transaction between account1 and account2 to see the results
	for i := 0; i < n; i++ {
		go func() {
			arg := TransferTxParam{
				FromAccountID: account1.AccountID,
				ToAccountID:   account2.AccountID,
				Amount:        amount,
			}

			result, err := store.TransferTx(context.Background(), arg)
			errChannel <- err
			transferChannel <- result
		}()
	}

	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errChannel
		require.NoError(t, err)

		transferResult := <-transferChannel
		require.NotEmpty(t, transferResult)

		transfer := transferResult.Transfer
		require.NotEmpty(t, transferResult)
		require.Equal(t, transfer.FromAccountID, account1.AccountID)
		require.Equal(t, transfer.ToAccountID, account2.AccountID)
		require.Equal(t, transfer.Amount, amount)
		require.NotZero(t, transfer.TransferID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), transfer.TransferID)
		require.NoError(t, err)

		fromEntry := transferResult.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, fromEntry.AccountID, account1.AccountID)
		require.Equal(t, fromEntry.Amount, -amount) //withdraw money
		require.NotEmpty(t, fromEntry.EntryID)
		require.NotZero(t, fromEntry.CreatedAt)
		_, err = store.GetEntry(context.Background(), fromEntry.EntryID)
		require.NoError(t, err)

		toEntry := transferResult.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, toEntry.AccountID, account2.AccountID)
		require.Equal(t, toEntry.Amount, amount) //money is going in
		require.NotZero(t, toEntry.EntryID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.EntryID)
		require.NoError(t, err)

		//TODO: Check Balances AS well

		fromAccount := transferResult.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.AccountID, account1.AccountID)

		toAccount := transferResult.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.AccountID, account2.AccountID)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0)

		k := int(diff1 / amount)
		//require.True(t, k >= 1 && k < n)
		require.NotContains(t, existed, k)
		fmt.Println(">> k", k, "Added")
		existed[k] = true
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.AccountID)
	require.NoError(t, err)

	fmt.Println(">> after Account 1", updatedAccount1.Balance, ">> after Account 2", updatedAccount2.Balance)
	require.Equal(t, account1.Balance-(int64(n)*amount), updatedAccount1.Balance)
	require.Equal(t, account2.Balance+(int64(n)*amount), updatedAccount2.Balance)
}

func TestTransferTxDeadLock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before Account 1", account1.Balance, ">> before Account 2", account2.Balance)

	n := 10
	amount := int64(10)

	errChannel := make(chan error)

	//run n concurrent transaction between account1 and account2 to see the results
	for i := 0; i < n; i++ {
		fromAccountId := account1.AccountID
		toAccountId := account2.AccountID

		if i%2 == 1 {
			fromAccountId = account2.AccountID
			toAccountId = account1.AccountID
		}

		go func() {
			arg := TransferTxParam{
				FromAccountID: fromAccountId,
				ToAccountID:   toAccountId,
				Amount:        amount,
			}

			_, err := store.TransferTx(context.Background(), arg)
			errChannel <- err

		}()
	}

	for i := 0; i < n; i++ {
		err := <-errChannel
		require.NoError(t, err)
	}

	updatedAccount1, err := store.GetAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)

	updatedAccount2, err := store.GetAccount(context.Background(), account2.AccountID)
	require.NoError(t, err)

	fmt.Println(">> after Account 1", updatedAccount1.Balance, ">> after Account 2", updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)
}
