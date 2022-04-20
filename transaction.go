package torm

import (
	"fmt"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/jmoiron/sqlx"
)

type Proc func(*sqlx.Tx) error

func Transaction(sql *sqlx.DB, proc Proc) (err error) {
	var tx *sqlx.Tx

	defer func() {
		if err != nil && tx != nil {
			if e := tx.Rollback(); e != nil {
				err = multierror.Append(err, fmt.Errorf("failed to rollback transaction: %v", e))
			}
		}
	}()

	tx, err = sql.Beginx()
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
