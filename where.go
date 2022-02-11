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
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return &WhereBuilder{context: w.context}
	}
	w.context.wheres = append(w.context.wheres, query+" = ? ")
	w.context.args = append(w.context.args, arg)
	return &WhereBuilder{context: w.context}
}

func (w *WhereBuilder) Ne(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return &WhereBuilder{context: w.context}
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return &WhereBuilder{context: w.context}
	}
	w.context.wheres = append(w.context.wheres, query+" <> ? ")
	w.context.args = append(w.context.args, arg)
	return &WhereBuilder{context: w.context}
}

func (w *WhereBuilder) Like(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}

	var key = ""
	switch arg.(type) {
	case string:
		key = "%" + arg.(string) + "%"
	case []byte:
		key = "%" + string(arg.([]byte)) + "%"
	case *string:
		key = "%" + *arg.(*string) + "%"
	case *[]byte:
		key = "%" + string(*arg.(*[]byte)) + "%"
	default:
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" LIKE ? ")
	w.context.args = append(w.context.args, key)
	return &WhereBuilder{context: w.context}
}

func (w *WhereBuilder) NoLike(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return &WhereBuilder{context: w.context}
		}
	}
	var key = ""
	switch arg.(type) {
	case string:
		key = "%" + arg.(string) + "%"
	case []byte:
		key = "%" + string(arg.([]byte)) + "%"
	case *string:
		key = "%" + *arg.(*string) + "%"
	case *[]byte:
		key = "%" + string(*arg.(*[]byte)) + "%"
	default:
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" NO  LIKE ? ")
	w.context.args = append(w.context.args, key)
	return &WhereBuilder{context: w.context}
}
