package torm

import (
	"errors"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pinnacles/torm/internal/test"
)

func TestTransactionCommit(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO test").WithArgs(1).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := Transaction(db, func(tx *sqlx.Tx) error {
			builder := NewBuilder(tx)
			return builder.Insert().Exec(&test.TestSchema{Foo: 1, Bar: 2})
		}); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestTransactionRollback(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO test").WithArgs(1).WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		if err := Transaction(db, func(tx *sqlx.Tx) error {
			builder := NewBuilder(tx)
			return builder.Insert().Exec(&test.TestSchema{Foo: 1, Bar: 2})
		}); err == nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}
