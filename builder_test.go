package torm

import (
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/pinnacles/torm/internal/test"
)

func TestBuilder(t *testing.T) {
	if err := test.WithSqlxMock(func(db *sqlx.DB, mock sqlmock.Sqlmock) {
		builder := NewBuilder(db)

		if builder.Insert() == nil {
			t.Error("insertBuilder is nil")
		}
		if builder.Select() == nil {
			t.Error("selectBuilder is nil")
		}
		if builder.Update() == nil {
			t.Error("updateBuilder is nil")
		}
		if builder.Delete() == nil {
			t.Error("deleteBuilder is nil")
		}
		if builder.Querier() == nil {
			t.Error("querier is nil")
		}
	}); err != nil {
		t.Fatal(err)
	}
}
