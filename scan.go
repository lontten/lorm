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

	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
)

type rowColumnType struct {
	index            int    // 字段再row中的位置
	noNull           bool   // true 字段必定不为null
	databaseTypeName string // 字段-数据库数据类型
}

// ScanLn
// 接收一行结果
// 1.ptr single/comp
// 2.slice- single
func (ctx ormContext) ScanLn(rows *sql.Rows) (num int64, err error) {
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)

	num = 0
	t := ctx.destBaseType
	v := ctx.destBaseValue
	tP := ctx.scanDest

	columns, err := rows.Columns()
	if err != nil {
		return
	}

	cfm := make(map[string]compC)
	if ctx.destBaseTypeIsComp {
		cfm = getColIndex2FieldNameMap(columns, t)
	}

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return 0, err
	}
	var rowColumnTypeMap = make(map[int]rowColumnType)
	for i, columnType := range columnTypes {
		nullable, ok := columnType.Nullable()
		rowColumnTypeMap[i] = rowColumnType{
			index:            i,
			databaseTypeName: columnType.DatabaseTypeName(),
			noNull:           ok && !nullable,
		}
	}

	if rows.Next() {
		box, convert := ctx.createColBox(v, tP, cfm, rowColumnTypeMap)
		err = rows.Scan(box...)
		if err != nil {
			return
		}
		err = convert()
		if err != nil {
			return
		}
		num++
	}

	if rows.Next() {
		return 0, errors.New("result to many for one")
	}
	return
}

// Scan
// 接收多行结果
func (ctx ormContext) Scan(rows *sql.Rows) (int64, error) {
	defer func(rows *sql.Rows) {
		utils.PanicErr(rows.Close())
	}(rows)

	var num int64 = 0
	t := ctx.destBaseType
	arr := ctx.scanV
	isPtr := ctx.destSliceItemIsPtr

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	cfm := getColIndex2FieldNameMap(columns, t)

	columnTypes, err := rows.ColumnTypes()
	if err != nil {
		return 0, err
	}

	var rowColumnTypeMap = make(map[int]rowColumnType)
	for i, columnType := range columnTypes {
		nullable, ok := columnType.Nullable()
		rowColumnTypeMap[i] = rowColumnType{
			index:            i,
			databaseTypeName: columnType.DatabaseTypeName(),
			noNull:           ok && !nullable,
		}
	}

	for rows.Next() {
		box, vp, v, convert := ctx.createColBoxNew(t, cfm, rowColumnTypeMap)

		err = rows.Scan(box...)
		if err != nil {
			return 0, err
		}
		err = convert()
		if err != nil {
			return 0, err
		}

		if isPtr {
			arr.Set(reflect.Append(arr, vp))
		} else {
			arr.Set(reflect.Append(arr, v))
		}
		num++
	}
	return num, nil
}
