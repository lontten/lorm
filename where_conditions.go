//  Copyright 2025 lontten lontten@163.com
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package lorm

import (
	"fmt"
	"reflect"

	"github.com/lontten/lorm/utils"
)

//------------------model/map/id------------------

// 过滤 软删除
func (w *WhereBuilder) Model(v any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	list := getStructCV(reflect.ValueOf(v))
	for _, cv := range list {
		if cv.isSoftDel || cv.isZero {
			continue
		}
		w.fieldValue(cv.columnName, cv.value)
	}
	return w
}

func (w *WhereBuilder) Map(v any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	list := getMapCV(reflect.ValueOf(v))
	for _, cv := range list {
		w.fieldValue(cv.columnName, cv.value)
	}
	return w
}
func (w *WhereBuilder) PrimaryKey(args ...any) *WhereBuilder {
	if w.err != nil {
		return w
	}
	argsLen := len(args)
	if argsLen == 0 {
		return w
	}
	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type: PrimaryKeys,
			args: args,
		},
	})
	return w
}

func (w *WhereBuilder) FilterPrimaryKey(args ...any) *WhereBuilder {
	if w.err != nil {
		return w
	}
	argsLen := len(args)
	if argsLen == 0 {
		return w
	}
	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type: FilterPrimaryKeys,
			args: args,
		},
	})
	return w
}

//------------------model/map/id-end------------------
//------------------eq------------------

// Eq
// x = ?
func (w *WhereBuilder) Eq(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	isNil := utils.IsNil(arg)
	if isNil {
		w.err = fmt.Errorf("invalid use of Eq: argument for field '%s' is nil", field)
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  Eq,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// NotEq
// x <> ?
func (w *WhereBuilder) NotEq(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	isNil := utils.IsNil(arg)
	if isNil {
		w.err = fmt.Errorf("invalid use of NotEq: argument for field '%s' is nil", field)
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  Neq,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// BoolIn
// IN (?)
func (w *WhereBuilder) BoolIn(condition bool, field string, args ...any) *WhereBuilder {
	if w.err != nil {
		return w
	}
	if !condition {
		return w
	}
	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  In,
			query: field,
			args:  args,
		},
	})
	return w
}

// In
// IN (?)
func (w *WhereBuilder) In(field string, args ...any) *WhereBuilder {
	if w.err != nil {
		return w
	}
	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  In,
			query: field,
			args:  args,
		},
	})
	return w
}

// BoolNotIn
// NOT IN (?)
func (w *WhereBuilder) BoolNotIn(condition bool, field string, args ...any) *WhereBuilder {
	if w.err != nil {
		return w
	}
	if !condition {
		return w
	}
	argsLen := len(args)
	if argsLen == 0 {
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  NotIn,
			query: field,
			args:  args,
		},
	})
	return w
}

// NotIn
// NOT IN (?)
func (w *WhereBuilder) NotIn(field string, args ...any) *WhereBuilder {
	if w.err != nil {
		return w
	}
	argsLen := len(args)
	if argsLen == 0 {
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  NotIn,
			query: field,
			args:  args,
		},
	})
	return w
}

// Contains
// pg 独有
// [1] @< [1,2]
func (w *WhereBuilder) Contains(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	isNil := utils.IsNil(arg)
	if isNil {
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  Contains,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// Lt
// x < a
func (w *WhereBuilder) Lt(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	isNil := utils.IsNil(arg)
	if isNil {
		w.err = fmt.Errorf("invalid use of Lt: argument for field '%s' is nil", field)
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  Less,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// Lte
// x <= a
func (w *WhereBuilder) Lte(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	isNil := utils.IsNil(arg)
	if isNil {
		w.err = fmt.Errorf("invalid use of Lte: argument for field '%s' is nil", field)
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  LessEq,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// Gt
// x > a
func (w *WhereBuilder) Gt(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	isNil := utils.IsNil(arg)
	if isNil {
		w.err = fmt.Errorf("invalid use of Gt: argument for field '%s' is nil", field)
		return w
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  Greater,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// Gte
// x >= a
func (w *WhereBuilder) Gte(field string, arg any, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}
	isNil := utils.IsNil(arg)
	if isNil {
		w.err = fmt.Errorf("invalid use of Gte: argument for field '%s' is nil", field)
		return w
	}
	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  GreaterEq,
			query: field,
			args:  []any{arg},
		},
	})
	return w
}

// IsNull
// x IS NULL
func (w *WhereBuilder) IsNull(field string, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  IsNull,
			query: field,
		},
	})
	return w
}

// IsNotNull
// x IS NOT NULL
func (w *WhereBuilder) IsNotNull(field string, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}

	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  IsNotNull,
			query: field,
		},
	})
	return w
}

// IsFalse
// x IS FALSE
func (w *WhereBuilder) IsFalse(field string, condition ...bool) *WhereBuilder {
	if w.err != nil {
		return w
	}
	for _, b := range condition {
		if !b {
			return w
		}
	}
	w.andWheres = append(w.andWheres, WhereBuilder{
		clause: &Clause{
			Type:  IsFalse,
			query: field,
		},
	})
	return w
}
