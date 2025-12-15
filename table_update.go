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

	"github.com/lontten/lorm/sqltype"
	"github.com/lontten/lorm/utils"
)

func Update(db Engine, dest any, wb *WhereBuilder, extra ...*ExtraContext) (int64, error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Update

	ctx.initModelDest(dest)
	ctx.initConf() //初始化表名，主键，自增id

	ctx.initColumnsValue() //初始化cv
	ctx.initColumnsValueExtra()

	ctx.initColumnsValueSoftDel() // 软删除

	ctx.wb.And(wb)

	dialect.tableUpdateGen()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return 0, nil
	}

	exec, err := db.exec(dialectSql, ctx.originalArgs...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func UpdateByPrimaryKey(db Engine, dest any, extra ...*ExtraContext) (int64, error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Update

	ctx.initModelDest(dest)
	ctx.initConf() //初始化表名，主键，自增id

	ctx.initColumnsValue() //初始化cv
	ctx.initColumnsValueExtra()

	ctx.initColumnsValueSoftDel() // 软删除

	wb := W()
	for _, name := range ctx.primaryKeyColumnNames {
		find := utils.Find(ctx.columns, name)
		if find == -1 {
			return 0, fmt.Errorf("primaryKey column %s not set value", name)
		} else {
			wb.fieldValue(name, ctx.columnValues[find])
			ctx.columnValues = append(ctx.columnValues[:find], ctx.columnValues[find+1:]...)
			ctx.columns = append(ctx.columns[:find], ctx.columns[find+1:]...)
		}
	}

	ctx.wb.And(wb)

	dialect.tableUpdateGen()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return 0, nil
	}

	exec, err := db.exec(dialectSql, ctx.originalArgs...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}
