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

import "reflect"

type StmtQueryContext[T any] struct {
	db   Stmter
	args []any
}

func StmtQuery[T any](db Stmter, args ...any) *StmtQueryContext[T] {
	return &StmtQueryContext[T]{
		db:   db,
		args: args,
	}
}

func (q *StmtQueryContext[T]) Convert(c Convert) *StmtQueryContext[T] {
	q.db.getCtx().convertCtx.Add(c)
	return q
}

func (q *StmtQueryContext[T]) One() (*T, error) {
	db := q.db
	args := q.args
	ctx := db.getCtx()

	dest := new(T)

	ctx.initScanDestOneT(dest)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(args...)
	if err != nil {
		return nil, err
	}
	num, err := ctx.ScanLn(rows)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, nil
	}
	return dest, nil
}

func (q *StmtQueryContext[T]) List() ([]T, error) {
	db := q.db
	args := q.args
	ctx := db.getCtx()

	var dest = &[]T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, false)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(args...)
	if err != nil {
		return nil, err
	}
	num, err := ctx.Scan(rows)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, nil
	}
	return *dest, nil
}

func (q *StmtQueryContext[T]) ListP() ([]*T, error) {
	db := q.db
	args := q.args
	ctx := db.getCtx()

	var dest = &[]*T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, true)
	if ctx.err != nil {
		return nil, ctx.err
	}

	rows, err := db.query(args...)
	if err != nil {
		return nil, err
	}
	num, err := ctx.Scan(rows)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return nil, nil
	}
	return *dest, nil
}
