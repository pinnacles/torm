package test

import (
	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type TestSchema struct {
	Foo int `db:"foo"`
	Bar int `json:"bar"`
}

func (s TestSchema) TableName() string {
	return "test"
}

func WithSqlxMock(proc func(db *sqlx.DB, mock sqlmock.Sqlmock)) error {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		return err
	}
	defer mdb.Close()
	db := sqlx.NewDb(mdb, "mysql")
	defer db.Close()

	proc(db, mock)
	return nil
}
