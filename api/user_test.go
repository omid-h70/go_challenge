package api

import (
	"fmt"
	"github.com/golang/mock/gomock"
	db "go_challenge/db/sqlc"
	"go_challenge/util"
	"reflect"
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
	return fmt.Sprintf("is equal to %w", e.password)
}

func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreatUserParamsMatcher{arg, password}
}
