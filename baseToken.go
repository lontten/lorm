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

	//主键名-列表,这里考虑到多主键
	primaryKeyNames []string
	//主键值-列表
	primaryKeyValues [][]interface{}

	//字段列表-not nil
	columns []string
	//值列表-多个-not nil
	columnValues []interface{}

	wb *WhereBuilder
}

// todo 下面未重构--------------
