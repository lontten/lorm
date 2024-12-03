// Copyright (c) 2024 lontten
// lorm is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lorm

import "reflect"

type baseTokenType int

const (
	tInsert baseTokenType = iota
	tUpdate
	tDelete
	tSelect
	tCount
	tExist
	tWhereModel
	tInsertOrUpdate
	tInsertIgnore

	tScanOne
	tScanFirst
	tScanList

	//	--------------------
	// 对应数据 t reflect.Type
	tTableName
	// 对应数据 pk 主键值列表
	tPrimaryKey
	// 对应数据 wb
	tWhereBuilder

	// 对应数据 v dest
	tTableNameDestValue
)

type baseToken struct {
	typ  baseTokenType
	dest any
	v    reflect.Value
	t    reflect.Type

	pk []any

	//主键名-列表,这里考虑到多主键
	primaryKeyNames []string
	//主键值-列表
	primaryKeyValues [][]any

	//字段列表-not nil
	columns []string
	//值列表-多个-not nil
	columnValues []any

	wb *WhereBuilder
}
