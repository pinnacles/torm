package torm

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type KV map[string]interface{}

type selectBuilder struct {
	h      handler
	fields []string
}

func newSelect(h handler, f ...string) *selectBuilder {
	return &selectBuilder{
		h:      h,
		fields: f,
	}
}

func (s *selectBuilder) Where(clause string, kv KV) *querySelectBuilder {
	return &querySelectBuilder{
		h:      s.h,
		fields: s.fields,
		clause: clause,
		kv:     kv,
	}
}

func (s *selectBuilder) Query(res interface{}) error {
	q := &querySelectBuilder{
		h:      s.h,
		fields: s.fields,
	}
	return q.Query(res)
}

type querySelectBuilder struct {
	h      handler
	fields []string
	clause string
	kv     KV
}

func (q *querySelectBuilder) ToSQL(res interface{}) (*SQL, error) {

	if reflect.TypeOf(res).Kind() != reflect.Ptr {
		return nil, fmt.Errorf("Query must be specified Ptr type")
	}

	var meta *tableMeta
	switch reflect.TypeOf(res).Elem().Kind() {
	case reflect.Slice:
		s, ok := interface{}(reflect.New(reflect.TypeOf(res).Elem().Elem()).Interface()).(schema)
		if !ok {
			return nil, fmt.Errorf("res is expected to pass schema type or slice of schema")
		}
		meta = metas[s.TableName()]
	default:
		s, ok := res.(schema)
		if !ok {
			return nil, fmt.Errorf("res is expected to pass schema type or slice of schema")
		}
		meta = metas[s.TableName()]
	}

	selectColumns := []string{"*"}
	if len(q.fields) > 0 {
		if q.fields[0] != "*" {
			selectColumns = q.fields
		}
	} else {
		selectColumns = meta.Fields
	}
	quoted := make([]string, 0, len(selectColumns))
	for _, col := range selectColumns {
		quoted = append(quoted, fmt.Sprintf("`%s`", col))
	}

	syntax := []string{fmt.Sprintf("SELECT %s FROM `%s`", strings.Join(quoted, ","), meta.TableName)}
	args := []interface{}{}
	if q.clause != "" {
		syntax = append(syntax, fmt.Sprintf("WHERE %s", q.clause))
		query, params, err := sqlx.Named(strings.Join(syntax, " "), q.kv)
		if err != nil {
			return nil, err
		}
		syntax = []string{query}
		args = params
	}

	asSliceForIn := false
	for _, arg := range args {
		if reflect.TypeOf(arg).Kind() == reflect.Slice {
			asSliceForIn = true
			break
		}
	}

	var query string
	var params []interface{}
	if asSliceForIn {
		var err error
		query, params, err = sqlx.In(syntax[0], args...)
		if err != nil {
			return nil, err
		}
		query = q.h.Rebind(query)
	} else {
		query = syntax[0]
		params = args
	}

	log.Infof("SQL: %s values: %#v", query, params)
	return &SQL{
		Query: query,
		Args:  params,
	}, nil
}

func (q *querySelectBuilder) Query(res interface{}) error {
	sql, err := q.ToSQL(res)
	if err != nil {
		return err
	}

	var resIsSlice bool
	switch reflect.TypeOf(res).Elem().Kind() {
	case reflect.Slice:
		resIsSlice = true
	default:
		resIsSlice = false
	}

	if resIsSlice {
		return q.h.Select(res, sql.Query, sql.Args...)
	} else {
		return q.h.Get(res, sql.Query, sql.Args...)
	}
}
