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

func (b *SqlBuilder[T]) One() (t *T, err error) {
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	dest := new(T)
	ctx.initScanDestOneT(dest)
	if ctx.err != nil {
		return nil, ctx.err
	}

	b.initSelectSql()

	dialect.getSql(b.query)
	ctx.originalArgs = b.args
	ctx.printSql()

	if ctx.noRun {
		return nil, nil
	}

	rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
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

func (b *SqlBuilder[T]) List() (list []T, err error) {
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	var dest = &[]T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, false)
	if ctx.err != nil {
		return nil, ctx.err
	}

	b.initSelectSql()

	dialect.getSql(b.query)
	ctx.originalArgs = b.args
	ctx.printSql()

	if ctx.noRun {
		return nil, nil
	}
	rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
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

func (b *SqlBuilder[T]) ListP() (list []*T, err error) {
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	var dest = &[]*T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, true)
	if ctx.err != nil {
		return nil, ctx.err
	}

	b.initSelectSql()

	dialect.getSql(b.query)
	ctx.originalArgs = b.args
	ctx.printSql()

	if ctx.noRun {
		return nil, nil
	}
	rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
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
