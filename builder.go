package torm

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

type querier interface {
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
}

type handler interface {
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Rebind(query string) string
	PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error)
}

type SQL struct {
	Query string
	Args  []interface{}
}

type Builder struct {
	h  handler
	ts *time.Time
}

func NewBuilder(h handler) *Builder {
	return &Builder{
		h: h,
	}
}

func (t Builder) Querier() querier {
	return t.h
}

func (t Builder) Select(f ...string) *selectBuilder {
	return newSelect(t.h, f...)
}

func (t Builder) Insert(f ...string) *insertBuilder {
	return newInsert(t.h, t.ts, f...)
}

func (t Builder) Update(f ...string) *updateBuilder {
	return newUpdate(t.h, t.ts, f...)
}

func (t Builder) Delete() *deleteBuilder {
	return newDelete(t.h)
}

func (t *Builder) SetTime(ts *time.Time) {
	t.ts = ts
}
