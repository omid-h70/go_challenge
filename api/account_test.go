package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "go_challenge/db/mock"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountApi(t *testing.T) {

	user, _ := randomUser(t)
	account := randomAccount(user.UserName)

	//defining a nameless struct and filling it !
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.AccountID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.AccountID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.AccountID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.AccountID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//start test server and send request
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
		})
	}
}

func TestCreateAccountApi(t *testing.T) {

	//GET IT FROM GIT !!!!

	/*
		acount := randomAccount()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)

		//build stubs
		store.EXPECT().
			GetAccount(gomock.Any(), gomock.Eq(account.ID)).
			Times(1).
			Return(account, nil)

		//start test server and send request
		server := NewServer(store)
		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d". account.ID)
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		server.router.ServeHTTP(recorder, request)
		//check response
		require.Equal(t, http.StatusOK, recorder.Code)
		requireBodyMatchAccount(t,
	*/
}

func TestListAccountsApi(t *testing.T) {

}

func randomAccount(owner string) db.Account {

	return db.Account{
		AccountID: util.RandomInt(1, 1000),
		Owner:     owner,
		Balance:   util.RandomMoney(),
		Currency:  util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, account []db.Account) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts db.Account
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Equal(t, account, gotAccounts)
}
