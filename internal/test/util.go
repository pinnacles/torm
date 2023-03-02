package test

import (
	"time"
	"context"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

type TestSchema struct {
	ID        int       `db:"id" torm:"autoIncrement"`
	Foo       int       `db:"foo"`
	Bar       int       `json:"bar"`
	CreatedAt time.Time `db:"created_at" torm:"autoCreateTime"`
	UpdatedAt time.Time `db:"updated_at" torm:"autoUpdateTime"`
}

func (s TestSchema) TableName() string {
	return "test"
}

func WithSqlxMock(proc func(ctx context.Context, db *sqlx.DB, mock sqlmock.Sqlmock)) error {
	mdb, mock, err := sqlmock.New()
	if err != nil {
		return err
	}
	defer mdb.Close()
	db := sqlx.NewDb(mdb, "mysql")
	defer db.Close()

	ctx := context.Background()
	proc(ctx, db, mock)
	return nil
}
