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

func (q *StmtQueryContext[T]) ScanOne(dest *T) (int64, error) {
	db := q.db
	ctx := db.getCtx()

	ctx.initScanDestOne(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}

	rows, err := q.db.query(q.args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLn(rows)
}

func (q *StmtQueryContext[T]) ScanList(dest *[]T) (int64, error) {
	db := q.db
	ctx := db.getCtx()

	ctx.initScanDestList(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}
	rows, err := q.db.query(q.args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLn(rows)
}

func (q *StmtQueryContext[T]) ScanListP(dest *[]T) (int64, error) {
	db := q.db
	ctx := db.getCtx()

	ctx.initScanDestList(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}
	rows, err := q.db.query(q.args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLn(rows)
}
