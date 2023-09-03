package gapi

import (
	"database/sql"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	db "go_challenge/db/sqlc"
	"go_challenge/pb"
	"go_challenge/util"
	"go_challenge/worker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"testing"
)

//mockdb , mockwk are aliases

/*
COOOOOOOOOOOOOOOOOOOOOOOOOOOOOOL !
Stronger Unit test with Custom gomock Matcher !
We implement a Custom Matcher based on Go mock Interfaces
*/
type eqCreatUserTxParamsMatcher struct {
	arg      db.CreateUserTxParams
	password string
	user     db.User
}

func (e eqCreatUserTxParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserTxParams)

	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
		return false
	}
	arg.AfterCreate(e.user)
	return err == nil
}

func (e eqCreatUserTxParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %w", e.password)
}

func EqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
	return eqCreatUserTxParamsMatcher{arg, password, user}
}

func randomUser(t *testing.T) (user db.User, password string) {

	password = util.RandomString(6)
	hashedPassword, err := util.GetHashedPassword(password)
	require.NoError(t, err)

	user = db.User{
		UserName:       util.RandomOwner(),
		HashedPassword: hashedPassword,
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	return
}

func TestCreateUserAPI(t *testing.T) {

	user, password := randomUser(t)

	testCases := []struct {
		name          string
		req           *pb.CreateUserRequest
		buildStubs    func(store *mockdb.MockStore, tsk *mockwr.MockDistributor)
		checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.CreateUserRequest{
				UserName: user.UserName,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, tsk *mockwr.MockDistributor) {
				arg := db.CreateUserTxParams{
					CreateUserParams: db.CreateUserParams{
						UserName: user.UserName,
						FullName: user.FullName,
						Email:    user.Email,
					},
				}
				store.EXPECT().
					CreateUserTx(gomock.Any(), EqCreateUserTxParams(arg, password, user)).
					Times(1).
					Return(db.CreateUserTxResult{user}, nil)

				payload := &worker.PayLoadSendVerifyEmail{
					UserName: user.UserName,
				}

				tsk.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), payload, gomock.Any()).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				createUser := res.GetUser()
				require.Equal(t, user.UserName, createUser.UserName)
				require.Equal(t, user.FullName, createUser.FullName)
				require.Equal(t, user.Email, createUser.Email)
			},
		},
		{
			name: "Internal Error",
			req: &pb.CreateUserRequest{
				UserName: user.UserName,
				Password: password,
				FullName: user.FullName,
				Email:    user.Email,
			},
			buildStubs: func(store *mockdb.MockStore, tsk *mockwr.MockDistributor) {

				store.EXPECT().
					CreateUserTx(gomock.Any(), gomock.Any()).
					Times(1).
					Return(db.CreateUserTxResult{}, sql.ErrConnDone)

				tsk.EXPECT().
					DistributeTaskSendVerifyEmail(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(0).
					Return(nil)
			},
			checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
				require.NoError(t, err)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Internal, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			/*if we use one ctrl here for both we will get a deadlock !*/
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)

			wrctrl := gomock.NewController(t)
			defer wrctrl.Finish()
			dst := mockwk.NewMockTaskDistributor(wrctrl)

			tc.buildStubs(store, dst)

			server := newTestServer(t, store, dst)
			res, err := server.CreateUser(ctx.Background(), tc.req)
			tc.checkResponse(t, res, err)
		})
	}
}
