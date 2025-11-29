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

func (b *SqlBuilder[T]) ScanOne(dest any) (rowsNum int64, err error) {
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	ctx.initScanDestOne(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}

	b.initSelectSql()
	dialect.getSql(b.query)
	ctx.originalArgs = b.args
	ctx.printSql()

	if ctx.noRun {
		return 0, nil
	}

	rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLn(rows)
}

func (b *SqlBuilder[T]) ScanList(dest *[]T) (rowsNum int64, err error) {
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	ctx.initScanDestList(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}
	b.initSelectSql()

	dialect.getSql(b.query)
	ctx.originalArgs = b.args
	ctx.printSql()

	if ctx.noRun {
		return 0, nil
	}
	rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
	if err != nil {
		return 0, err
	}
	return ctx.Scan(rows)
}

func (b *SqlBuilder[T]) ScanListP(dest *[]*T) (rowsNum int64, err error) {
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	db := b.db
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	ctx.initScanDestList(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}
	b.initSelectSql()

	dialect.getSql(b.query)
	ctx.originalArgs = b.args
	ctx.printSql()

	if ctx.noRun {
		return 0, nil
	}
	rows, err := db.query(ctx.dialectSql, ctx.originalArgs...)
	if err != nil {
		return 0, err
	}
	return ctx.Scan(rows)
}
