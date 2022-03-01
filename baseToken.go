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
	tPrimaryKey

	tScanOne
	tScanFirst
	tScanList
)

type baseToken struct {
	typ  baseTokenType
	dest interface{}
	v    reflect.Value

	pk []interface{}

	where *WhereBuilder
}
