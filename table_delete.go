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
	"github.com/lontten/lorm/sqltype"
)

func Delete[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (int64, error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Delete

	dest := new(T)
	ctx.initScanDestOneT(dest)
	ctx.initConf() //初始化表名，主键，自增id

	ctx.initColumnsValueSoftDel()

	ctx.wb.And(wb)

	dialect.tableDelGen()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return 0, nil
	}

	exec, err := db.exec(dialectSql, ctx.args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}
