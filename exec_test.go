package torm

import (
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pinnacles/torm/internal/test"
)

func init() {
	Register(test.TestSchema{})
}

func TestInsert(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO test (foo) VALUES`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		builder := NewBuilder(db)
		ts := test.TestSchema{
			Foo: 1,
			Bar: 2,
		}
		if _, err := builder.Insert().Exec(ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestInsertColumn(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO test (bar) VALUES`)).
			WithArgs(2).
			WillReturnResult(sqlmock.NewResult(1, 1))

		builder := NewBuilder(db)
		ts := test.TestSchema{
			Foo: 1,
			Bar: 2,
		}
		if _, err := builder.Insert("bar").Exec(ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE test SET foo=? WHERE foo = ?`)).
			WithArgs(1, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		builder := NewBuilder(db)
		ts := test.TestSchema{
			Foo: 1,
			Bar: 2,
		}
		if _, err := builder.Update().Where("foo = :foo").Exec(ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestUpdateColumn(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE test SET bar=? WHERE foo = ?`)).
			WithArgs(2, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		builder := NewBuilder(db)
		ts := test.TestSchema{
			Foo: 1,
			Bar: 2,
		}
		if _, err := builder.Update("bar").Where("foo = :foo").Exec(ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM test WHERE foo = ?`)).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		builder := NewBuilder(db)
		ts := test.TestSchema{
			Foo: 1,
			Bar: 2,
		}
		if _, err := builder.Delete().Where("foo = :foo").Exec(ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	}); err != nil {
		t.Fatal(err)
	}
}
