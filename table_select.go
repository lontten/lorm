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
	"database/sql"
	"reflect"

	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm/sqltype"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
)

// First 根据条件获取第一个
func First[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (t *T, err error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...) // 表名，whenUpdateSet，select配置
	ctx.sqlType = sqltype.Select
	ctx.limit = types.NewInt64(1)

	dest := new(T)
	ctx.initScanDestOneT(dest)
	ctx.initConf() //初始化表名，主键，自增id

	if ctx.err != nil {
		return nil, ctx.err
	}
	if ctx.lastSql == "" {
		if ctx.autoPrimaryKeyColumnName != "" {
			ctx.lastSql = " ORDER BY " + ctx.autoPrimaryKeyColumnName + " DESC"
		}
	}

	ctx.initColumns()
	ctx.initColumnsValueSoftDel()

	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return nil, nil
	}

	rows, err := db.query(dialectSql, ctx.args...)
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

func List[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (list []T, err error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...) // 表名，whenUpdateSet，select配置
	ctx.sqlType = sqltype.Select

	var dest = &[]T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, false)
	ctx.initConf() //初始化表名，主键，自增id
	ctx.initColumns()
	ctx.initColumnsValueSoftDel()

	if ctx.err != nil {
		return nil, ctx.err
	}
	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return nil, nil
	}

	rows, err := db.query(dialectSql, ctx.args...)
	if err != nil {
		return nil, err
	}
	num, err := ctx.Scan(rows)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return make([]T, 0), nil
	}
	return *dest, nil
}

func ListP[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (list []*T, err error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Select

	var dest = &[]*T{}
	v := reflect.ValueOf(dest).Elem()
	baseV := reflect.ValueOf(new(T)).Elem()
	t := baseV.Type()

	ctx.initScanDestListT(dest, v, baseV, t, false)
	ctx.initConf() //初始化表名，主键，自增id
	ctx.initColumns()
	ctx.initColumnsValueSoftDel()

	if ctx.err != nil {
		return nil, ctx.err
	}
	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return nil, nil
	}

	rows, err := db.query(dialectSql, ctx.args...)
	if err != nil {
		return nil, err
	}
	num, err := ctx.Scan(rows)
	if err != nil {
		return nil, err
	}
	if num == 0 {
		return make([]*T, 0), nil
	}
	return *dest, nil
}

// Has
func Has[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (t bool, err error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Select
	ctx.modelSelectFieldNames = []string{"1"}

	dest := new(T)
	ctx.initScanDestOneT(dest)

	ctx.initConf() //初始化表名，主键，自增id
	ctx.initColumnsValueSoftDel()

	if ctx.err != nil {
		return false, ctx.err
	}
	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return false, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return false, nil
	}

	rows, err := db.query(dialectSql, ctx.args...)
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)
	return rows.Next(), nil
}

// Count
func Count[T any](db Engine, wb *WhereBuilder, extra ...*ExtraContext) (t int64, err error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...)
	ctx.sqlType = sqltype.Select
	ctx.modelSelectFieldNames = []string{"COUNT(*)"}

	dest := new(T)
	ctx.initScanDestOneT(dest)
	ctx.initConf() //初始化表名，主键，自增id
	ctx.initColumnsValueSoftDel()

	if ctx.err != nil {
		return 0, ctx.err
	}
	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return 0, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()
	if ctx.noRun {
		return 0, nil
	}

	var total int64
	rows, err := db.query(dialectSql, ctx.args...)
	if err != nil {
		return 0, err
	}
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)
	for rows.Next() {
		box := reflect.ValueOf(&total).Interface()
		err = rows.Scan(box)
		if err != nil {
			return total, err
		}
		return total, nil
	}
	return total, errors.New("rows no data")
}

// GetOrInsert
// d insert 的 对象，
// e 通用设置，select 自定义字段
func GetOrInsert[T any](db Engine, wb *WhereBuilder, d *T, extra ...*ExtraContext) (*T, error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...) // 表名，whenUpdateSet，select配置
	ctx.sqlType = sqltype.Select

	dest := new(T)
	ctx.initScanDestOneT(dest)
	ctx.initConf() //初始化表名，主键，自增id
	ctx.initColumns()
	ctx.initColumnsValueSoftDel()

	if ctx.err != nil {
		return nil, ctx.err
	}
	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()

	if ctx.noRun {
		return nil, nil
	}
	rows, err := db.query(dialectSql, ctx.args...)
	if err != nil {
		return nil, err
	}
	num, err := ctx.ScanLn(rows)
	if err != nil {
		return nil, err
	}
	if num != 0 {
		return dest, nil
	}

	//------------

	ctx.query.Reset()
	ctx.args = []any{}
	ctx.sqlType = sqltype.Insert

	ctx.initModelDest(d) //初始化参数

	ctx.initColumnsValue() //初始化cv
	ctx.initColumnsValueExtra()
	ctx.initColumnsValueSoftDel() // 软删除

	dialect.tableInsertGen()
	if ctx.hasErr() {
		return nil, ctx.err
	}

	dialect.getSql()
	dialectSql = ctx.dialectSql
	ctx.printSql()

	if ctx.noRun {
		return nil, nil
	}

	if ctx.returnAutoPrimaryKey == pkQueryReturn {
		rows, err = db.query(dialectSql, ctx.args...)
		if err != nil {
			return nil, err
		}
		num, err = ctx.ScanLn(rows)
		if err != nil {
			return nil, err
		}
		if num == 0 {
			return nil, errors.New("insert affected 0")
		}
		return d, nil
	}

	exec, err := db.exec(dialectSql, ctx.args...)
	if err != nil {
		return nil, err
	}
	if ctx.returnAutoPrimaryKey == pkFetchReturn {
		id, err := exec.LastInsertId()
		if err != nil {
			return nil, err
		}
		if id > 0 {
			ctx.setLastInsertId(id)
			if ctx.hasErr() {
				return nil, ctx.err
			}
		}
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errors.New("insert affected 0")
	}
	return d, nil
}

// HasOrInsert 根据条件查询是否已存在
// 如果存在，直接返回true，如果不存在返回false，并直接插入
// 应用场景：例如添加 用户 时，如果名字已存在，返回名字重复，否者正常添加。
// d insert 的 对象，
// e 通用设置，select 自定义字段
func HasOrInsert(db Engine, wb *WhereBuilder, d any, extra ...*ExtraContext) (bool, error) {
	db = db.init()
	dialect := db.getDialect()
	ctx := dialect.getCtx()
	ctx.initExtra(extra...) // 表名，whenUpdateSet，select配置
	ctx.modelSelectFieldNames = []string{"1"}
	ctx.sqlType = sqltype.Select

	ctx.initModelDest(d) //初始化参数
	ctx.initConf()       //初始化表名，主键，自增id

	ctx.initColumnsValue() //初始化cv
	ctx.initColumnsValueExtra()
	ctx.initColumnsValueSoftDel() // 软删除
	if ctx.err != nil {
		return false, ctx.err
	}

	ctx.wb.And(wb)

	dialect.tableSelectGen()
	if ctx.hasErr() {
		return false, ctx.err
	}

	dialect.getSql()
	dialectSql := ctx.dialectSql
	ctx.printSql()

	if ctx.noRun {
		return false, nil
	}
	rows, err := db.query(dialectSql, ctx.args...)
	if err != nil {
		return false, err
	}
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)
	if rows.Next() {
		return true, nil
	}

	//------------

	ctx.query.Reset()
	ctx.args = []any{}
	ctx.sqlType = sqltype.Insert
	ctx.initColumnsValueSoftDel() // 软删除

	dialect.tableInsertGen()
	if ctx.hasErr() {
		return false, ctx.err
	}

	dialect.getSql()
	dialectSql = ctx.dialectSql
	ctx.printSql()

	if ctx.noRun {
		return false, nil
	}

	if ctx.returnAutoPrimaryKey == pkQueryReturn {
		rows, err = db.query(dialectSql, ctx.args...)
		if err != nil {
			return false, err
		}
		_, err = ctx.ScanLn(rows)
		if err != nil {
			return false, err
		}
		return false, nil
	}

	exec, err := db.exec(dialectSql, ctx.args...)
	if err != nil {
		return false, err
	}
	if ctx.returnAutoPrimaryKey == pkFetchReturn {
		id, err := exec.LastInsertId()
		if err != nil {
			return false, err
		}
		if id > 0 {
			ctx.setLastInsertId(id)
			if ctx.hasErr() {
				return false, ctx.err
			}
		}
	}
	affected, err := exec.RowsAffected()
	if err != nil {
		return false, err
	}
	if affected == 0 {
		return false, errors.New("insert affected 0")
	}
	return false, nil
}
