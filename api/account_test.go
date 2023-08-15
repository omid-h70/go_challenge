package api

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	db "go_challenge/db/sqlc"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAccountApi1(t *testing.T) {

	//defining a nameless struct and filling it !
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.Mockstore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOk, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.Mockstore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}

	for i:= range testCases{
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T){
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			store := mockdb.NewMockStore(ctrl)
			tc.buildStubs(store)

			//start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := fmt.Sprintf("/accounts/%d". account.ID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
		}}
	}
}

func TestGetAccountApi(t *testing.T) {

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

func randomAccount() /*db.Account*/ {
	/*
		return db.Account{
		ID: util.RandomInt(1, 1000),
		Owner: util.RandomOwner(),
		Balance: util.RandomMoney(),
		Currency: util.RandomCurrency(),
		}
	*/
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer) {
	//data, err := io.ReadAll(body)
	//require.NoError(t, err)

	/*
		var getAccount, db.Account
		err = json.Unmarshal(data, &gotAccount)
		require.NoError(t, err)
		require.Equal(t, account, gotAccount)
	*/
}
