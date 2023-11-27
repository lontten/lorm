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
	dest interface{}
	v    reflect.Value
	t    reflect.Type

	pk []interface{}

	wb *WhereBuilder
}

// todo 下面未重构--------------
