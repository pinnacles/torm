package torm

import (
	"context"
	"errors"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pinnacles/torm/internal/test"
)

func TestTransactionCommit(t *testing.T) {
	if err := test.WithSqlxMock(func(ctx context.Context, db *sqlx.DB, mock sqlmock.Sqlmock) {
		tm := time.Now()
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `test`").WithArgs(1, tm, tm).WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		if err := Transaction(context.Background(), nil, db, func(tx *sqlx.Tx) error {
			builder := NewBuilder(tx)
			builder.SetTime(&tm)
			_, err := builder.Insert().Exec(ctx, &test.TestSchema{Foo: 1, Bar: 2})
			return err
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
	if err := test.WithSqlxMock(func(ctx context.Context, db *sqlx.DB, mock sqlmock.Sqlmock) {
		tm := time.Now()
		mock.ExpectBegin()
		mock.ExpectExec("INSERT INTO `test`").WithArgs(1, tm, tm).WillReturnError(errors.New("insert error"))
		mock.ExpectRollback()

		if err := Transaction(context.Background(), nil, db, func(tx *sqlx.Tx) error {
			builder := NewBuilder(tx)
			builder.SetTime(&tm)
			_, err := builder.Insert().Exec(ctx, &test.TestSchema{Foo: 1, Bar: 2})
			return err
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
