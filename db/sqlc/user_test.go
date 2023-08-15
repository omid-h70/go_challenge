package db

import (
	"testing"
)

func createRandomUser(t *testing.T) {

	/*
			arg := CreateUserParams{
				Username: 		util.RandomOwner(),
				HashedPassword: "secret",
				FullName:	    util.RandomOwner(),
				Email: 			util.RandomEmail()
			}
			user, err := testQueries.CreateUser(context.Background()arg)
			require.NoError(t, err)
			require.NoError(t, user)

			require.Equal(t, arg.Username, user.Username)
			require.Equal(t, arg.HashedPassword, user.HashedPassword)
		    require.Equal(t, arg.FullName, user.FullName)
		    require.Equal(t, arg.Email, user.Email)

			require.TRUE(t, user.PasswordChangedAt.IsZero())
		    require.NotZero(t, user.CreatedAt)
	*/
}

func TestGetUser(t *testing.T) {
	/*
		user1 := createRandomUser(t)
		user2 := testQueries.GetUser(context.Background(), user1.Username)
		require.NoError(t, err)
		require.NoError(t, user2)

		require.Equal(t, user1.Username, user2.Username)
		require.Equal(t, user1.HashedPassword, user2.HashedPassword)
		require.Equal(t, user1.FullName, user2.FullName)
		require.Equal(t, user1.Email, user2.Email)

		require.WithinDurationf(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
		require.WithinDurationf(t, user1.CreatedAt, user2.CreatedAt, time.Second)
		*
	*/
}
