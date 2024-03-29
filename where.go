package lorm

import (
	"bytes"
	"reflect"
	"strings"
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
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" = ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

func (w *WhereBuilder) In(query string, args ArgArray, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	if len(args) == 0 {
		return w
	}
	var queryArr []string
	for i := 0; i < len(args); i++ {
		queryArr = append(queryArr, "?")
	}

	var str = query + " IN (" + strings.Join(queryArr, ",") + ") "

	w.context.wheres = append(w.context.wheres, str)
	w.context.args = append(w.context.args, args...)
	return w
}

func (w *WhereBuilder) NotIn(query string, args ArgArray, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	if len(args) == 0 {
		return w
	}
	var queryArr []string
	for i := 0; i < len(args); i++ {
		queryArr = append(queryArr, "?")
	}

	var str = query + " NOT IN (" + strings.Join(queryArr, ",") + ") "

	w.context.wheres = append(w.context.wheres, str)
	w.context.args = append(w.context.args, args...)
	return w
}

func (w *WhereBuilder) NotEq(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" != ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

// 小于
func (w *WhereBuilder) Less(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" < ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

// 小于等于
func (w *WhereBuilder) LessEq(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" <= ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

// 大于
func (w *WhereBuilder) Greater(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" > ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

// 大于等于
func (w *WhereBuilder) GreaterEq(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" >= ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

func (w *WhereBuilder) Between(query string, arg1, arg2 interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.context.wheres = append(w.context.wheres, query+" BETWEEN ? AND ? ")
	w.context.args = append(w.context.args, arg1)
	w.context.args = append(w.context.args, arg2)
	return w
}

func (w *WhereBuilder) Arg(arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w.context.args = append(w.context.args, arg)
	return w
}

func (w *WhereBuilder) Args(args ...interface{}) *WhereBuilder {
	w.context.args = append(w.context.args, args...)
	return w
}

func (w *WhereBuilder) IsNull(query string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.context.wheres = append(w.context.wheres, query+" IS NULL ")
	return w
}

func (w *WhereBuilder) IsNotNull(query string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.context.wheres = append(w.context.wheres, query+" IS NOT NULL ")
	return w
}

func (w *WhereBuilder) NotBetween(query string, arg1, arg2 interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.context.wheres = append(w.context.wheres, query+" NOT BETWEEN ? AND ? ")
	w.context.args = append(w.context.args, arg1)
	w.context.args = append(w.context.args, arg2)
	return w
}

func (w *WhereBuilder) And(query string, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.context.wheres = append(w.context.wheres, query)
	return w
}

func (w *WhereBuilder) Ne(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
	}
	w.context.wheres = append(w.context.wheres, query+" <> ? ")
	w.context.args = append(w.context.args, arg)
	return w
}

func (w *WhereBuilder) Like(query string, arg interface{}, condition ...bool) *WhereBuilder {
	for _, b := range condition {
		if !b {
			return w
		}
	}
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
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
	arg = getTargetInter(reflect.ValueOf(arg))
	if arg == nil {
		return w
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
