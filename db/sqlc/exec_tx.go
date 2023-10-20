package db

import (
	"context"
	"fmt"
)

func (store *SqlStore) execTx(ctx context.Context, fn func(q *Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	/*
		&sql.TxOptions{
			Isolation: sql.LevelWriteCommitted,
			ReadOnly:  false,
		})
	*/

	if err != nil {
		return err
	}

	q := New(tx)
	execErr := fn(q)
	if execErr != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx error:%v rb error:%v", rbErr, execErr)
		}
	}
	return tx.Commit()
}
