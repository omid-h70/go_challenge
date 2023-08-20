// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1

package db

import (
	"time"
)

type Account struct {
	AccountID int64
	Owner     string
	Balance   int64
	Currency  string
	CreatedAt time.Time
}

type Entry struct {
	EntryID   int64
	AccountID int64
	// it can be negative or positive
	Amount    int64
	CreatedAt time.Time
}

type Transfer struct {
	TransferID    int64
	FromAccountID int64
	ToAccountID   int64
	// it most be positive
	Amount    int64
	CreatedAt time.Time
}

type User struct {
	UserName          string
	HashedPassword    string
	FullName          string
	Email             string
	PasswordChangedAt time.Time
	CreatedAt         time.Time
}
