package torm

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type insertBuilder struct {
	h      handler
	fields []string
	ts     *time.Time
}

func newInsert(h handler, ts *time.Time, f ...string) *insertBuilder {
	return &insertBuilder{
		h:      h,
		fields: f,
		ts:     ts,
	}
}

func (b *insertBuilder) ToSQL(s schema) (*SQL, error) {
	meta := metas[s.TableName()]

	specifyAutoCreateTimeCol := map[string]bool{}
	specifyAutoUpdateTimeCol := map[string]bool{}
	if len(b.fields) <= 0 {
		fs := []string{}
		for _, f := range meta.Fields {
			if meta.IsAutoIncrement(f) {
				continue
			}
			fs = append(fs, f)
		}
		b.fields = fs
	} else {
		for _, fn := range b.fields {
			if meta.IsAutoCreateTime(fn) {
				specifyAutoCreateTimeCol[fn] = true
				continue
			}
			if meta.IsAutoUpdateTime(fn) {
				specifyAutoUpdateTimeCol[fn] = true
				continue
			}
		}
		for k := range meta.AutoCreateTimeColumns {
			if _, ok := specifyAutoCreateTimeCol[k]; !ok {
				b.fields = append(b.fields, k)
			}
		}
		for k := range meta.AutoUpdateTimeColumns {
			if _, ok := specifyAutoUpdateTimeCol[k]; !ok {
				b.fields = append(b.fields, k)
			}
		}
	}

	ts := time.Now()
	if b.ts != nil {
		ts = *b.ts
	}
	elem := reflect.ValueOf(s).Elem()
	for k, v := range meta.AutoCreateTimeColumns {
		if _, ok := specifyAutoCreateTimeCol[k]; !ok {
			f := elem.FieldByName(v)
			if f.IsValid() {
				if f.CanSet() {
					if f.Kind() == reflect.Struct {
						f.Set(reflect.ValueOf(ts))
					}
				}
			}
		}
	}
	for k, v := range meta.AutoUpdateTimeColumns {
		if _, ok := specifyAutoUpdateTimeCol[k]; !ok {
			f := elem.FieldByName(v)
			if f.IsValid() {
				if f.CanSet() {
					if f.Kind() == reflect.Struct {
						f.Set(reflect.ValueOf(ts))
					}
				}
			}
		}
	}

	names := make([]string, 0, len(b.fields))
	for _, n := range b.fields {
		names = append(names, ":"+n)
	}
	syntax := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, meta.TableName, strings.Join(b.fields, ","), strings.Join(names, ","))

	log.Infof("SQL: %s value: %#v", syntax, s)
	return &SQL{
		Query: syntax,
		Args:  []interface{}{s},
	}, nil
}

func (b *insertBuilder) Exec(s schema) (sql.Result, error) {
	sql, err := b.ToSQL(s)
	if err != nil {
		return nil, err
	}
	return b.h.NamedExec(sql.Query, s)
}

type updateBuilder struct {
	h      handler
	fields []string
	ts     *time.Time
}

func newUpdate(h handler, ts *time.Time, f ...string) *updateBuilder {
	return &updateBuilder{
		h:      h,
		fields: f,
		ts:     ts,
	}
}

func (b *updateBuilder) Where(clause string) *execUpdateBuilder {
	return &execUpdateBuilder{
		h:      b.h,
		fields: b.fields,
		clause: clause,
		ts:     b.ts,
	}
}

type execUpdateBuilder struct {
	h      handler
	fields []string
	clause string
	ts     *time.Time
}

func (b *execUpdateBuilder) ToSQL(s schema) (*SQL, error) {
	meta := metas[s.TableName()]

	specifyAutoUpdateTimeCol := map[string]bool{}
	if len(b.fields) <= 0 {
		fs := []string{}
		for _, f := range meta.Fields {
			if meta.IsAutoIncrement(f) {
				continue
			}
			if meta.IsAutoCreateTime(f) {
				continue
			}
			fs = append(fs, f)
		}
		b.fields = fs
	} else {
		for _, fn := range b.fields {
			if meta.IsAutoUpdateTime(fn) {
				specifyAutoUpdateTimeCol[fn] = true
			}
		}
		for k := range meta.AutoUpdateTimeColumns {
			if _, ok := specifyAutoUpdateTimeCol[k]; !ok {
				b.fields = append(b.fields, k)
			}
		}
	}

	ts := time.Now()
	if b.ts != nil {
		ts = *b.ts
	}
	elem := reflect.ValueOf(s).Elem()
	for k, v := range meta.AutoUpdateTimeColumns {
		if _, ok := specifyAutoUpdateTimeCol[k]; !ok {
			f := elem.FieldByName(v)
			if f.IsValid() {
				if f.CanSet() {
					if f.Kind() == reflect.Struct {
						f.Set(reflect.ValueOf(ts))
					}
				}
			}
		}
	}

	fields := make([]string, 0, len(b.fields))
	for _, n := range b.fields {
		fields = append(fields, fmt.Sprintf("%s=:%s", n, n))
	}
	syntax := []string{fmt.Sprintf(`UPDATE %s SET %s`, meta.TableName, strings.Join(fields, ","))}
	if b.clause != "" {
		syntax = append(syntax, fmt.Sprintf("WHERE %s", b.clause))
	}
	syntax = []string{strings.Join(syntax, " ")}

	log.Infof("SQL: %s value: %#v", syntax[0], s)
	return &SQL{
		Query: syntax[0],
		Args:  []interface{}{s},
	}, nil
}

func (b *execUpdateBuilder) Exec(s schema) (sql.Result, error) {
	sql, err := b.ToSQL(s)
	if err != nil {
		return nil, err
	}
	return b.h.NamedExec(sql.Query, s)
}

type deleteBuilder struct {
	h handler
}

func newDelete(h handler) *deleteBuilder {
	return &deleteBuilder{
		h: h,
	}
}

func (b *deleteBuilder) Where(clause string) *execDeleteBuilder {
	return &execDeleteBuilder{
		h:      b.h,
		clause: clause,
	}
}

type execDeleteBuilder struct {
	h      handler
	clause string
}

func (b *execDeleteBuilder) ToSQL(s schema) (*SQL, error) {
	meta := metas[s.TableName()]
	query, args, err := sqlx.Named(fmt.Sprintf("DELETE FROM %s WHERE %s", meta.TableName, b.clause), s)
	if err != nil {
		return nil, err
	}

	log.Infof("SQL: %s value: %#v", query, args)
	return &SQL{
		Query: query,
		Args:  args,
	}, nil
}

func (b *execDeleteBuilder) Exec(s schema) (sql.Result, error) {
	sql, err := b.ToSQL(s)
	if err != nil {
		return nil, err
	}
	return b.h.Exec(sql.Query, sql.Args...)
}
