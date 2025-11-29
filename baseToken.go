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
