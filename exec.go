package torm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type insertBuilder struct {
	h      handler
	fields []string
}

func newInsert(h handler, f ...string) *insertBuilder {
	return &insertBuilder{
		h:      h,
		fields: f,
	}
}

func (e *insertBuilder) ToSQL(s schema) (*SQL, error) {
	meta := metas[s.TableName()]

	if len(e.fields) <= 0 {
		fs := make([]string, 0, len(meta.Fields))
		for _, f := range meta.Fields {
			fs = append(fs, f)
		}
		e.fields = fs
	}

	names := make([]string, 0, len(e.fields))
	for _, n := range e.fields {
		names = append(names, ":"+n)
	}
	syntax := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, meta.TableName, strings.Join(e.fields, ","), strings.Join(names, ","))

	log.Infof("SQL: %s value: %#v", syntax, s)
	return &SQL{
		Query: syntax,
		Args:  []interface{}{s},
	}, nil
}

func (e *insertBuilder) Exec(s schema) (sql.Result, error) {
	sql, err := e.ToSQL(s)
	if err != nil {
		return nil, err
	}
	return e.h.NamedExec(sql.Query, s)
}

type updateBuilder struct {
	h      handler
	fields []string
}

func newUpdate(h handler, f ...string) *updateBuilder {
	return &updateBuilder{
		h:      h,
		fields: f,
	}
}

func (u *updateBuilder) Where(clause string) *execUpdateBuilder {
	return &execUpdateBuilder{
		h:      u.h,
		fields: u.fields,
		clause: clause,
	}
}

type execUpdateBuilder struct {
	h      handler
	fields []string
	clause string
}

func (e *execUpdateBuilder) ToSQL(s schema) (*SQL, error) {
	meta := metas[s.TableName()]

	if len(e.fields) <= 0 {
		fs := make([]string, 0, len(meta.Fields))
		for _, f := range meta.Fields {
			fs = append(fs, f)
		}
		e.fields = fs
	}

	clauses := make([]string, 0, len(e.fields))
	for _, n := range e.fields {
		clauses = append(clauses, fmt.Sprintf("%s=:%s", n, n))
	}
	syntax := []string{fmt.Sprintf(`UPDATE %s SET %s`, meta.TableName, strings.Join(clauses, ","))}
	if e.clause != "" {
		syntax = append(syntax, fmt.Sprintf("WHERE %s", e.clause))
	}
	syntax = []string{strings.Join(syntax, " ")}

	log.Infof("SQL: %s value: %#v", syntax[0], s)
	return &SQL{
		Query: syntax[0],
		Args:  []interface{}{s},
	}, nil
}

func (e *execUpdateBuilder) Exec(s schema) (sql.Result, error) {
	sql, err := e.ToSQL(s)
	if err != nil {
		return nil, err
	}
	return e.h.NamedExec(sql.Query, s)
}

type deleteBuilder struct {
	h handler
}

func newDelete(h handler) *deleteBuilder {
	return &deleteBuilder{
		h: h,
	}
}

func (d *deleteBuilder) Where(clause string) *execDeleteBuilder {
	return &execDeleteBuilder{
		h:      d.h,
		clause: clause,
	}
}

type execDeleteBuilder struct {
	h      handler
	clause string
}

func (e *execDeleteBuilder) ToSQL(s schema) (*SQL, error) {
	meta := metas[s.TableName()]
	query, args, err := sqlx.Named(fmt.Sprintf("DELETE FROM %s WHERE %s", meta.TableName, e.clause), s)
	if err != nil {
		return nil, err
	}

	log.Infof("SQL: %s value: %#v", query, args)
	return &SQL{
		Query: query,
		Args:  args,
	}, nil
}

func (e *execDeleteBuilder) Exec(s schema) (sql.Result, error) {
	sql, err := e.ToSQL(s)
	if err != nil {
		return nil, err
	}
	return e.h.Exec(sql.Query, sql.Args...)
}
