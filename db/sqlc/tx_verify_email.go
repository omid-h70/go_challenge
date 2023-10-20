package db

import (
	"context"
	"database/sql"
)

// VerifyEmailTxParam Contains input parameter for Transfer Transaction
type VerifyEmailTxParam struct {
	EmailId    int64
	SecretCode string
}

type VerifyEmailTxResult struct {
	User        User
	VerifyEmail VerifyEmail
}

// VerifyEmailTx outbox outbox outbox
// we are performing outbox pattern here
// vvvvvvvvvvvvvvoooooooooooolllllllllllllllllllaaaaaaa
// =>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>
func (store *SqlStore) VerifyEmailTx(ctx context.Context, arg VerifyEmailTxParam) (VerifyEmailTxResult, error) {
	var result VerifyEmailTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		result.VerifyEmail, err = q.UpdateVerifyEmail(ctx, UpdateVerifyEmailParams{
			ID:         arg.EmailId,
			SecretCode: arg.SecretCode,
		})
		if err != nil {
			return err
		}
		result.User, err = q.UpdateUser(ctx, UpdateUserParams{
			UserName: result.VerifyEmail.Username,
			IsEmailVerified: sql.NullBool{
				Valid: true,
				Bool:  true,
			},
		})
		return err
	})
	return result, err
}
