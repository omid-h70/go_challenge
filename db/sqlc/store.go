package db

import (
	"context"
	"database/sql"
)

type Store interface {
	execTx(ctx context.Context)
}

type SqlStore struct {
	db *sql.DB
	//*Queries
}

func (store *SqlStore) execTx(ctx context.Context) {

}
