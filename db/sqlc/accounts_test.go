package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go_challenge/util"
	"testing"
	"time"
)

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func createRandomAccount(t *testing.T) Account {
	randomUser := createRandomUser(t)

	arg := CreateAccountParams{
		Owner:    randomUser.UserName,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, account.Balance, arg.Balance)
	require.Equal(t, account.Currency, arg.Currency)

	require.NotZero(t, account.AccountID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)
	dbAccount, err := testQueries.GetAccount(context.Background(), account.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, dbAccount)
	require.Equal(t, account.AccountID, dbAccount.AccountID)
	require.WithinDuration(t, account.CreatedAt.Time, dbAccount.CreatedAt.Time, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateAccountParams{
		AccountID: account.AccountID,
		Balance:   util.RandomMoney(),
	}
	dbAccount, err := testQueries.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, dbAccount.Balance, arg.Balance)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.AccountID)
	require.NoError(t, err)

	//Just to make sure
	dbAccount, err := testQueries.GetAccount(context.Background(), account.AccountID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, dbAccount)
}

func TestListAccounts(t *testing.T) {
	//var lastAccount Account
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		//Owner : lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accounts)

	for _, account := range accounts {
		require.NotEmpty(t, account)
		//require.Equal(t, lastAccount.Owner, account.Owner)
	}

}
