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
	tWhereBuilder
	tWhereModel
	tInsertOrUpdate
	tInsertIgnore

	tScanOne
	tScanFirst
	tScanList

	//	--------------------
	// 对应数据 t reflect.Type
	tableName
	// 对应数据 pk 主键值列表
	tPrimaryKey

	// 对应数据 v dest
	tableNameDestValue
)

type baseToken struct {
	typ  baseTokenType
	dest interface{}
	v    reflect.Value
	t    reflect.Type

	pk []interface{}

	where *WhereBuilder
}

// todo 下面未重构--------------
