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

func TestSelectOne(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `*` FROM `test`")).
			WillReturnRows(sqlmock.NewRows([]string{"foo"}).AddRow(1))

		builder := NewBuilder(db)
		ts := test.TestSchema{}
		if err := builder.Select("*").Query(&ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		if ts.Foo != 1 {
			t.Fatal("ts.Foo is 1 was expected")
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestSelectMulti(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `id`,`foo`,`created_at`,`updated_at` FROM `test`")).
			WillReturnRows(sqlmock.NewRows([]string{"foo"}).AddRow(1).AddRow(2))

		builder := NewBuilder(db)
		ts := []test.TestSchema{}
		if err := builder.Select().Query(&ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		if len(ts) != 2 {
			t.Fatal("length of ts is 2 was expected")
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestSelectColumn(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `foo`,`bar` FROM `test`")).
			WillReturnRows(sqlmock.NewRows([]string{"foo"}).AddRow(1).AddRow(2))

		builder := NewBuilder(db)
		ts := []test.TestSchema{}
		if err := builder.Select("foo", "bar").Query(&ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		if len(ts) != 2 {
			t.Fatal("length of ts is 2 was expected")
		}
	}); err != nil {
		t.Fatal(err)
	}
}

func TestSeletWithWhere(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		mock.ExpectQuery(regexp.QuoteMeta("SELECT `*` FROM `test` WHERE foo = ?")).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"foo"}).AddRow(1).AddRow(1))

		builder := NewBuilder(db)
		ts := []test.TestSchema{}
		if err := builder.Select("*").Where("foo = :foo", KV{"foo": 1}).Query(&ts); err != nil {
			t.Fatal(err)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
		if len(ts) != 2 {
			t.Fatal("length of ts is 2 was expected")
		}
	}); err != nil {
		t.Fatal(err)
	}
}
