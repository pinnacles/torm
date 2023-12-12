package torm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Proc func(*sqlx.Tx) error

func Transaction(ctx context.Context, opts *sql.TxOptions, sql *sqlx.DB, proc Proc) (err error) {
	var tx *sqlx.Tx

	defer func() {
		if err != nil && tx != nil {
			if e := tx.Rollback(); e != nil {
				err = errors.Join(err, fmt.Errorf("failed to rollback transaction: %v", e))
			}
		}
	}()

	tx, err = sql.BeginTxx(ctx, opts)
	if err != nil {
		return
	}

	if err = proc(tx); err != nil {
		return
	}

	if err = tx.Commit(); err != nil {
		return
	}

	return nil
}
