package lorm

import (
	"bytes"
	"reflect"
)

type WhereContext struct {
	wheres []string
	args   []interface{}
}

type WhereBuilder struct {
	context WhereContext
}

func (w *WhereBuilder) toWhereSqlOneself() ([]byte, []interface{}) {
	wheres := w.context.wheres
	var bb bytes.Buffer
	bb.WriteString("WHERE ")
	for i, where := range wheres {
		if i == 0 {
			bb.WriteString(" WHERE " + where)
			continue
		}
		bb.WriteString(" AND " + where)
	}
	return bb.Bytes(), w.context.args
}

func (w *WhereBuilder) Eq(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return &WhereBuilder{context: w.context}
		}
	}
	t := reflect.ValueOf(arg)
	w.context.wheres = append(w.context.wheres, query+" = ? ")
	if t.Kind() == reflect.Ptr {
		arg = t.Elem()
	}
	w.context.args = append(w.context.args, arg)
	return &WhereBuilder{context: w.context}
}

func (w *WhereBuilder) Ne(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return &WhereBuilder{context: w.context}
		}
	}
	t := reflect.ValueOf(arg)
	w.context.wheres = append(w.context.wheres, query+" <> ? ")
	if t.Kind() == reflect.Ptr {
		arg = t.Elem()
	}
	w.context.args = append(w.context.args, arg)
	return &WhereBuilder{context: w.context}
}

func (w *WhereBuilder) Like(query, arg string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return &WhereBuilder{context: w.context}
		}
	}
	w.context.wheres = append(w.context.wheres, query+" LIKE ? ")
	w.context.args = append(w.context.args, "%"+arg+"%")
	return &WhereBuilder{context: w.context}
}

func (w *WhereBuilder) NoLike(query, arg string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return &WhereBuilder{context: w.context}
		}
	}
	w.context.wheres = append(w.context.wheres, query+" NO  LIKE ? ")
	w.context.args = append(w.context.args, "%"+arg+"%")
	return &WhereBuilder{context: w.context}
}
