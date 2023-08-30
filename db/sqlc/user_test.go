package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/require"
	"go_challenge/util"
	"testing"
	"time"
)

func createRandomUser(t *testing.T) User {

	arg := CreateUserParams{
		UserName:       util.RandomOwner(),
		HashedPassword: "secret",
		FullName:       util.RandomOwner(),
		Email:          util.RandomEmail(),
	}
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)

	require.Equal(t, arg.UserName, user.UserName)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestGetUser(t *testing.T) {

	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.UserName)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.UserName, user2.UserName)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.Email, user2.Email)

	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}

func TestUpdateUserOnlyFullName(t *testing.T) {

	oldUser := createRandomUser(t)
	newFullName := util.RandomOwner()
	usr, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		UserName: oldUser.UserName,
		FullName: sql.NullString{
			Valid:  true,
			String: oldUser.HashedPassword,
		},
	})

	require.NotEmpty(t, usr)
	require.NoError(t, err)
	require.Equal(t, usr.UserName, newFullName)
	require.NotEqual(t, usr.FullName, oldUser.FullName)
}

func TestUpdateUserOnlyEmail(t *testing.T) {

	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	usr, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		UserName: oldUser.UserName,
		Email: sql.NullString{
			Valid:  true,
			String: newEmail,
		},
	})

	require.NotEmpty(t, usr)
	require.NoError(t, err)
	require.Equal(t, usr.Email, newEmail)
	require.NotEqual(t, usr.Email, oldUser.Email)
}

func TestUpdateUserOnlyHashedPassword(t *testing.T) {
	//TODO
}

func TestUpdateUser(t *testing.T) {

	oldUser := createRandomUser(t)
	newEmail := util.RandomEmail()
	newFullName := util.RandomOwner()
	hashPassword, err := util.GetHashedPassword(util.RandomString(5))
	require.NoError(t, err)

	usr, err := testQueries.UpdateUser(context.Background(), UpdateUserParams{
		UserName: oldUser.UserName,
		FullName: sql.NullString{
			Valid:  true,
			String: newFullName,
		},
		Email: sql.NullString{
			Valid:  true,
			String: newEmail,
		},
		HashedPassword: sql.NullString{
			Valid:  true,
			String: hashPassword,
		},
	})

	require.NotEmpty(t, usr)
	require.NoError(t, err)

	require.NotEqual(t, usr.Email, oldUser.Email)
	require.NotEqual(t, usr.FullName, oldUser.FullName)
	require.NotEqual(t, usr.HashedPassword, oldUser.HashedPassword)

	require.Equal(t, usr.Email, newEmail)
	require.Equal(t, usr.HashedPassword, hashPassword)
	require.Equal(t, usr.FullName, newFullName)
}
