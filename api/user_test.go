package api

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
	"reflect"
	"testing"
)

/*
COOOOOOOOOOOOOOOOOOOOOOOOOOOOOOL !
Stronger Unit test with Custom gomock Matcher !
We implement a Custom Matcher based on Go mock Interfaces
*/
type eqCreatUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

func (e eqCreatUserParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)

	if !ok {
		return false
	}

	err := util.CheckPassword(e.password, arg.HashedPassword)
	if err != nil {
		return false
	}
	e.arg.HashedPassword = arg.HashedPassword
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreatUserParamsMatcher) String() string {
	return fmt.Sprintf("is equal to %s", e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreatUserParamsMatcher{arg, password}
}

func TestCreateUser(t *testing.T) {

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
