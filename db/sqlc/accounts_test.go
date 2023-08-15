package db

import (
	"testing"
)

func TestListAccounts(t *testing.T) {
	//var lastAccount Account
	for i := 0; i < 10; i++ {
		//lastAccount = createRandoAccount(t)
	}

	arg := ListAccountsParams{
		//Owner : lastAccount.Owner,
		Limit:  5,
		Offset: 0,
	}

	/*
		accounts, err := testQueries.ListAccounts(context.Background(), arg)
		require.NoError(t, err)
		require.NotEmpty(t, accounts)

		for _, account := range accounts {
			require.NotEmpty(t, account)
			require.Equal(t, lastAccount.Owner, account.Owner)
		}
	*/
}
