package lorm

import (
	"bytes"
	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/utils"
	"reflect"
	"strings"
	"sync"
)

type OrmConf struct {
	//po生成文件目录
	PoDir string
	//是否覆盖，默认true
	IsFileOverride bool

	//作者
	Author string
	//是否开启ActiveRecord模式,默认false
	IsActiveRecord bool

	IdType int

	//表名
	//TableNameFun >  tag > TableNamePrefix
	TableNamePrefix string
	TableNameFun    func(t reflect.Value, dest any) string

	//主键 默认为id
	PrimaryKeyNames   []string
	PrimaryKeyNameFun func(v reflect.Value, dest any) []string

	//多租户
	TenantIdFieldName    string                      //多租户的  租户字段名 空字符串极为不启用多租户
	TenantIdValueFun     func() any                  //租户的id值，获取函数
	TenantIgnoreTableFun func(tableName string) bool //该表是否忽略多租户，true忽略该表，即没有多租户
}

var typeTableNameCache = map[reflect.Type]string{}
var typeTableNameMu sync.Mutex

func getTypeTableName(t reflect.Type, tableNamePrefix string) string {
	s, ok := typeTableNameCache[t]
	if ok {
		return s
	}
	typeTableNameMu.Lock()
	defer typeTableNameMu.Unlock()
	s, ok = typeTableNameCache[t]
	if ok {
		return s
	}

	name := t.String()
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[index+1:]
	}
	name = utils.Camel2Case(name)
	if tableNamePrefix != "" {
		name = tableNamePrefix + name
	}
	typeTableNameCache[t] = name
	return name
}

// 不可缓存
// 获取表名
func (c OrmConf) tableName(v reflect.Value, dest any) string {
	// fun
	tableNameFun := c.TableNameFun
	if tableNameFun != nil {
		return tableNameFun(v, dest)
	}

	// tableName
	n := GetTableName(v)
	if n != nil {
		return *n
	}

	// structName
	t := v.Type()
	name := getTypeTableName(t, c.TableNamePrefix)
	return name
}

// 不可缓存
// 1.默认主键为id，
// 2.可以PrimaryKeyNames设置主键字段名
// 3.通过表名动态设置主键字段名-fn
func (c OrmConf) primaryKeys(v reflect.Value, dest any) []string {
	//fun
	primaryKeyNameFun := c.PrimaryKeyNameFun
	if primaryKeyNameFun != nil {
		return primaryKeyNameFun(v, dest)
	}

	list := GetPrimaryKeyNames(v)
	if len(list) > 0 {
		return list
	}

	// id
	return []string{"id"}
}

// 可缓存
func (c OrmConf) autoIncrements(v reflect.Value) []string {
	return GetAutoIncrements(v)
}

// 可以缓存
// 获取model字段对应的 db name，没有软删除
func (c OrmConf) getStructField(t reflect.Type) (columns []string, err error) {
	fiMap := getStructColName2fieldNameMap(t)
	arr := make([]string, len(fiMap))
	for f := range fiMap {
		arr = append(arr, f)
	}
	return arr, nil
}

// 可以缓存
// 获取model字段对应的 db name，包含软删除
func (c OrmConf) getStructAllField(t reflect.Type) (columns []string, err error) {
	fiMap := getStructColName2fieldNameAllMap(t)
	arr := make([]string, 0)
	for f := range fiMap {
		arr = append(arr, f)
	}
	return arr, nil
}

type colName2fieldNameMap map[string]string

type compCV struct {
	//有效字段列表
	columns []string
	//有效值列表
	columnValues []field.Value

	//零值字段列表
	modelZeroFieldNames []string

	//所有字段列表
	modelAllFieldNames []string

	//所有字段 dbName:fieldName
	modelAllFieldNameMap colName2fieldNameMap
}

// 忽略 软删除 字段
// 获取 struct 对应的字段名 和 其值
// 返回值类型有 Val,一种
func getStructCV(v reflect.Value) (compCV, error) {
	t := v.Type()
	cv := compCV{
		columns:             make([]string, 0),
		columnValues:        make([]field.Value, 0),
		modelZeroFieldNames: make([]string, 0),
		modelAllFieldNames:  make([]string, 0),
	}
	structFieldIndexMap := getStructColName2fieldNameMap(t)
	cv.modelAllFieldNameMap = structFieldIndexMap

	for column, i := range structFieldIndexMap {
		fieldV := v.FieldByName(i)
		inter := getFieldInterZero(fieldV)
		cv.modelAllFieldNames = append(cv.modelAllFieldNames, column)
		if inter != nil {
			cv.columns = append(cv.columns, column)
			cv.columnValues = append(cv.columnValues, field.Value{
				Type:  field.Val,
				Value: inter,
			})
		} else {
			cv.modelZeroFieldNames = append(cv.modelZeroFieldNames, column)
		}
	}

	return cv, nil
}

// 获取map[string]any
// 返回值类型有 None,Null,Val,三种
func getMapCV(v reflect.Value) (compCV, error) {
	cv := compCV{
		columns:             make([]string, 0),
		columnValues:        make([]field.Value, 0),
		modelZeroFieldNames: make([]string, 0),
		modelAllFieldNames:  make([]string, 0),
	}

	for _, k := range v.MapKeys() {
		inter := getFieldInter(v.MapIndex(k))

		cv.columns = append(cv.columns, k.String())
		cv.columnValues = append(cv.columnValues, inter)
	}
	return cv, nil
}

// 获取 rows 返回数据，每个字段index 对应 struct 的字段 名字
func getColIndex2FieldNameMap(columns []string, t reflect.Type) (ColIndex2FieldNameMap, error) {
	if isValuerType(t) {
		return ColIndex2FieldNameMap{}, nil
	}

	colNum := len(columns)
	ci2fm := make([]string, colNum)
	cf := getStructColName2fieldNameMap(t)

	validNum := 0
	for i, column := range columns {
		fieldName, ok := cf[column]
		if !ok {
			ci2fm[i] = ""
			continue
		}
		ci2fm[i] = fieldName
		validNum++
	}

	if colNum == 1 && validNum == 0 {
		return ColIndex2FieldNameMap{}, nil
	}
	return ci2fm, nil
}

// tableName表名
// keys
// hasTen true开启多租户
func (c OrmConf) genDelSqlCommon(tableName string, keys []string) []byte {
	var bb bytes.Buffer

	//hasTen := c.TenantIdFieldName != "" && !c.TenantIgnoreTableFun(tableName)
	//whereSql := c.GenWhere(keys, hasTen)

	//logicDeleteSetSql := c.LogicDeleteSetSql
	//logicDeleteYesSql := c.LogicDeleteYesSql
	//if logicDeleteSetSql == "" {
	//	bb.WriteString("DELETE FROM ")
	//	bb.WriteString(tableName)
	//	bb.WriteString(string(whereSql))
	//} else {
	//	bb.WriteString("UPDATE ")
	//	bb.WriteString(tableName)
	//	bb.WriteString(" SET ")
	//	bb.WriteString(logicDeleteSetSql)
	//	bb.WriteString(string(whereSql))
	//	bb.WriteString(" and ")
	//	bb.WriteString(logicDeleteYesSql)
	//}
	return bb.Bytes()
}

// tableName表名
// keys
// hasTen true开启多租户
func (c OrmConf) genDelSqlByWhere(tableName string, where []byte) []byte {
	//hasTen := c.TenantIdFieldName != "" && !c.TenantIgnoreTableFun(tableName)

	var bb bytes.Buffer
	//whereSql := c.whereExtra(where, hasTen)
	//
	//logicDeleteSetSql := c.LogicDeleteSetSql
	//logicDeleteYesSql := c.LogicDeleteYesSql
	//lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	//logicDeleteYesSql = strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	//if logicDeleteSetSql == lgSql {
	//	bb.WriteString("DELETE FROM ")
	//	bb.WriteString(tableName)
	//	bb.Write(whereSql)
	//} else {
	//	bb.WriteString("UPDATE ")
	//	bb.WriteString(tableName)
	//	bb.WriteString(" SET ")
	//	bb.WriteString(lgSql)
	//	bb.Write(whereSql)
	//	bb.WriteString(" and ")
	//	bb.WriteString(logicDeleteYesSql)
	//}
	return bb.Bytes()
}

// GenWhere 有tenantId功能
func (c OrmConf) GenWhere(keys []string, hasTen bool) []byte {
	if hasTen {
		keys = append(keys, c.TenantIdFieldName)
	}
	if len(keys) == 0 {
		return []byte("")
	}

	var bb bytes.Buffer
	bb.WriteString(" WHERE ")
	bb.WriteString(keys[0])
	bb.WriteString(" = ? ")
	for i := 1; i < len(keys); i++ {
		bb.WriteString(" AND ")
		bb.WriteString(keys[i])
		bb.WriteString(" = ? ")
	}

	return bb.Bytes()
}

// 有tenantid功能
func (c OrmConf) whereExtra(where []byte, hasTen bool) []byte {
	var bb bytes.Buffer
	//bb.Write(where)
	//
	//logicDeleteYesSql := c.LogicDeleteYesSql
	//lg := strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	//if lg != logicDeleteYesSql {
	//	bb.WriteString(" and ")
	//	bb.WriteString(lg)
	//}
	//
	//if hasTen {
	//	bb.WriteString(" AND ")
	//	bb.WriteString(c.TenantIdFieldName)
	//	bb.WriteString(" = ? ")
	//}

	return bb.Bytes()
}

// tableName表名
// columns
func (c OrmConf) genSelectSqlCommon(tableName string, columns []string) []byte {

	var bb bytes.Buffer
	bb.WriteString(" SELECT ")
	for i, column := range columns {
		if i == 0 {
			bb.WriteString(column)
		} else {
			bb.WriteString(" , ")
			bb.WriteString(column)
		}
	}
	bb.WriteString(" FROM ")
	bb.WriteString(tableName)
	return bb.Bytes()
}
