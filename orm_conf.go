package lorm

import (
	"bytes"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
	"strings"
	"unicode"
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
	TableNameFun    func(structName string, dest interface{}) string

	//字段名
	FieldNamePrefix string

	//主键 默认为id
	PrimaryKeyNames   []string
	PrimaryKeyNameFun func(tableName string) []string

	//逻辑删除 logicDeleteFieldName不为零值，即开启
	// LogicDeleteYesSql   deleted_at is null
	// LogicDeleteNoSql   deleted_at is not null
	// LogicDeleteSetSql   deleted_at = now()
	LogicDeleteYesSql string
	LogicDeleteNoSql  string
	LogicDeleteSetSql string

	//多租户
	TenantIdFieldName    string                      //多租户的  租户字段名 空字符串极为不启用多租户
	TenantIdValueFun     func() interface{}          //租户的id值，获取函数
	TenantIgnoreTableFun func(tableName string) bool //该表是否忽略多租户，true忽略该表，即没有多租户
}

// v03 从 type中获取表名
func (c OrmConf) tableName(t reflect.Type) (string, error) {

	// fun
	name := t.String()
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[index+1:]
	}
	name = utils.Camel2Case(name)

	tableNameFun := c.TableNameFun
	if tableNameFun != nil {
		return tableNameFun(name, t), nil
	}

	// tag

	numField := t.NumField()
	tagTableName := ""
	for i := 0; i < numField; i++ {
		if tag := t.Field(i).Tag.Get("tableName"); tag != "" {
			if tagTableName == "" {
				tagTableName = tag
			} else {
				return "", errors.New("has to many tableName tag")
			}
		}
	}
	if tagTableName != "" {
		return tagTableName, nil
	}

	// structName
	tableNamePrefix := c.TableNamePrefix
	if tableNamePrefix != "" {
		return tableNamePrefix + name, nil
	}

	return name, nil
}

// v03 不可以缓存，因为fn
// 1.默认主键为id，
// 2.可以PrimaryKeyNames设置主键字段名
// 3.通过表名动态设置主键字段名-fn
func (c OrmConf) primaryKeys(tableName string) []string {
	//fun
	primaryKeyNameFun := c.PrimaryKeyNameFun
	if primaryKeyNameFun != nil {
		return primaryKeyNameFun(tableName)
	}

	//conifg id name
	primaryKeyName := c.PrimaryKeyNames
	if len(primaryKeyName) != 0 {
		return primaryKeyName
	}

	// id
	return []string{"id"}
}

// v03 可以缓存
//
//	主键Id、ID，都转化为id
//
// tag== db:name  可以自定义名字
// tag== core:-  跳过
// 过滤掉首字母小写的字段
// 只获取字段对应的 数据库 字段名
func (c OrmConf) initColumns(t reflect.Type) (columns []string, err error) {

	cMap := make(map[string]int)

	numField := t.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		name := field.Name
		if name == "ID" {
			cMap["id"] = i
			num++
			if len(cMap) < num {
				return columns, errors.New("字段:: id  error")
			}
			continue
		}

		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}
		name = utils.Camel2Case(name)

		if tag := field.Tag.Get("core"); tag == "-" {
			continue
		}

		if tag := field.Tag.Get("db"); tag != "" {
			name = tag
			cMap[name] = i
			num++
			if len(cMap) < num {
				return columns, errors.New("字段::" + "error")
			}
			continue
		}

		fieldNamePrefix := c.FieldNamePrefix
		if fieldNamePrefix != "" {
			cMap[fieldNamePrefix+name] = i
			num++
			if len(cMap) < num {
				return columns, errors.New("字段::" + "error")
			}
			continue
		}

		cMap[name] = i
		num++
		if len(cMap) < num {
			return columns, errors.New("字段::" + "error")
		}
	}
	arr := make([]string, len(cMap))

	var i = 0
	for s := range cMap {
		arr[i] = s
		i++
	}
	return arr, nil
}

// v03 可以缓存
//
//	主键Id、ID，都转化为id
//
// tag== db:name  可以自定义名字
// tag== core:-  跳过
// 过滤掉首字母小写的字段
// 获取struct对应的数据字段名：和其在struct中的index下标
func (c OrmConf) getStructMappingColumns(t reflect.Type) (map[string]int, error) {
	cMap := make(map[string]int)

	numField := t.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		name := field.Name

		if name == "ID" {
			cMap["id"] = i
			num++
			if len(cMap) < num {
				return cMap, errors.New("字段::id" + "error")
			}
			continue
		}

		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}
		name = utils.Camel2Case(name)

		if tag := field.Tag.Get("core"); tag == "-" {
			continue
		}

		if tag := field.Tag.Get("db"); tag != "" {
			name = tag
			cMap[name] = i
			num++
			if len(cMap) < num {
				return cMap, errors.New("字段::" + "error")
			}
			continue
		}

		fieldNamePrefix := c.FieldNamePrefix
		if fieldNamePrefix != "" {
			cMap[fieldNamePrefix+name] = i
			num++
			if len(cMap) < num {
				return cMap, errors.New("字段::" + "error")
			}
			continue
		}

		cMap[name] = i
		num++
		if len(cMap) < num {
			return cMap, errors.New("字段::" + "error")
		}
	}

	return cMap, nil
}

type compCV struct {
	//字段列表-not nil
	columns []string
	//值列表-多个-not nil
	columnValues []interface{}

	//字段列表-nil
	nilColumns []string
	//值列表-多个-nil
	nilColumnValues []interface{}

	//字段列表-all
	allColumns []string
	//值列表-多个-all
	allColumnValues []interface{}
}

// 03
// 获取comp :struct/map 对应的字段名 和 其值
func (c OrmConf) getCompCV(v reflect.Value) (compCV, error) {
	t := v.Type()
	cv := compCV{
		columns:         make([]string, 0),
		columnValues:    make([]interface{}, 0),
		nilColumns:      make([]string, 0),
		nilColumnValues: make([]interface{}, 0),
		allColumns:      make([]string, 0),
		allColumnValues: make([]interface{}, 0),
	}
	if _isStructType(t) {
		mappingColumns, err := c.getStructMappingColumns(t)
		if err != nil {
			return cv, err
		}

		for column, i := range mappingColumns {
			inter := getFieldInter(v.Field(i))
			cv.allColumns = append(cv.allColumns, column)
			cv.allColumnValues = append(cv.allColumnValues, inter)

			if inter != nil {
				cv.columns = append(cv.columns, column)
				cv.columnValues = append(cv.columnValues, inter)
			} else {
				cv.nilColumns = append(cv.nilColumns, column)
				cv.nilColumnValues = append(cv.nilColumnValues, inter)
			}
		}
	} else {
		for _, k := range v.MapKeys() {
			inter := getFieldInter(v.MapIndex(k))

			cv.allColumns = append(cv.allColumns, k.String())
			cv.allColumnValues = append(cv.allColumnValues, inter)

			cv.columns = append(cv.columns, k.String())
			cv.columnValues = append(cv.columnValues, inter)
		}
	}
	return cv, nil
}

// 03
// 获取comp :struct/map 对应的字段名
func (c OrmConf) getCompC(t reflect.Type) (compCV, error) {
	cv := compCV{
		columns:         make([]string, 0),
		columnValues:    make([]interface{}, 0),
		nilColumns:      make([]string, 0),
		nilColumnValues: make([]interface{}, 0),
		allColumns:      make([]string, 0),
		allColumnValues: make([]interface{}, 0),
	}
	if _isStructType(t) {
		mappingColumns, err := c.getStructMappingColumns(t)
		if err != nil {
			return cv, err
		}
		for column := range mappingColumns {
			cv.allColumns = append(cv.allColumns, column)
		}
	} else {
		//map 无法获取 字段名
	}
	return cv, nil
}

// todo 下面未重构--------------

func (c OrmConf) getColFieldIndexLinkMap(columns []string, t reflect.Type) (ColFieldIndexLinkMap, error) {
	if isValuerType(t) {
		return ColFieldIndexLinkMap{}, nil
	}

	colNum := len(columns)
	cfm := make([]int, colNum)
	fm, err := getFieldMap(t, c.FieldNamePrefix)
	if err != nil {
		return nil, err
	}

	validNum := 0
	for i, column := range columns {
		index, ok := fm[column]
		if !ok {
			cfm[i] = -1
			continue
		}
		cfm[i] = index
		validNum++
	}

	if colNum == 1 && validNum == 0 {
		return ColFieldIndexLinkMap{}, nil
	}
	return cfm, nil
}

// tableName表名
// keys
// hasTen true开启多租户
func (c OrmConf) genDelSqlCommon(tableName string, keys []string) []byte {
	var bb bytes.Buffer

	hasTen := c.TenantIdFieldName != "" && !c.TenantIgnoreTableFun(tableName)
	whereSql := c.GenWhere(keys, hasTen)

	logicDeleteSetSql := c.LogicDeleteSetSql
	logicDeleteYesSql := c.LogicDeleteYesSql
	if logicDeleteSetSql == "" {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.WriteString(string(whereSql))
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(logicDeleteSetSql)
		bb.WriteString(string(whereSql))
		bb.WriteString(" and ")
		bb.WriteString(logicDeleteYesSql)
	}
	return bb.Bytes()
}

// tableName表名
// keys
// hasTen true开启多租户
func (c OrmConf) genDelSqlByWhere(tableName string, where []byte) []byte {
	hasTen := c.TenantIdFieldName != "" && !c.TenantIgnoreTableFun(tableName)

	var bb bytes.Buffer
	whereSql := c.whereExtra(where, hasTen)

	logicDeleteSetSql := c.LogicDeleteSetSql
	logicDeleteYesSql := c.LogicDeleteYesSql
	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	logicDeleteYesSql = strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	if logicDeleteSetSql == lgSql {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.Write(whereSql)
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(lgSql)
		bb.Write(whereSql)
		bb.WriteString(" and ")
		bb.WriteString(logicDeleteYesSql)
	}
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
	bb.Write(where)

	logicDeleteYesSql := c.LogicDeleteYesSql
	lg := strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	if lg != logicDeleteYesSql {
		bb.WriteString(" and ")
		bb.WriteString(lg)
	}

	if hasTen {
		bb.WriteString(" AND ")
		bb.WriteString(c.TenantIdFieldName)
		bb.WriteString(" = ? ")
	}

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
