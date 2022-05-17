package torm

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
)

type querier interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	NamedExec(query string, arg interface{}) (sql.Result, error)
}

type handler interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Rebind(query string) string
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

type SQL struct {
	Query string
	Args  []interface{}
}

type Builder struct {
	h handler
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
	return newInsert(t.h, f...)
}

func (t Builder) Update(f ...string) *updateBuilder {
	return newUpdate(t.h, f...)
}

func (t Builder) Delete() *deleteBuilder {
	return newDelete(t.h)
}
