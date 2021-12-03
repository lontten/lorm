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
	PrimaryKeyNameFun func(tableName string, base reflect.Value) []string

	//逻辑删除 logicDeleteFieldName不为零值，即开启
	// LogicDeleteYesSql   lg.deleted_at is null
	// LogicDeleteNoSql   lg.deleted_at is not null
	// LogicDeleteSetSql   lg.deleted_at = now()
	LogicDeleteYesSql string
	LogicDeleteNoSql  string
	LogicDeleteSetSql string

	//多租户 tenantIdFieldName不为零值，即开启
	TenantIdFieldName    string
	TenantIdValueFun     func() interface{}
	TenantIgnoreTableFun func(tableName string, base reflect.Value) bool
}

// v0.7
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

// v0.6
func (c OrmConf) primaryKeys(tableName string, v reflect.Value) []string {
	//fun
	primaryKeyNameFun := c.PrimaryKeyNameFun
	if primaryKeyNameFun != nil {
		return primaryKeyNameFun(tableName, v)
	}

	//conifg id name
	primaryKeyName := c.PrimaryKeyNames
	if len(primaryKeyName) != 0 {
		return primaryKeyName
	}

	// id
	return []string{"id"}
}

//v0.7
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

//v0.6
//获取struct对应的字段名 有效部分
func (c OrmConf) getStructMappingColumns(t reflect.Type) (map[string]int, error) {
	cMap := make(map[string]int)

	numField := t.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		name := field.Name

		if name == "ID" {
			cMap["Id"] = i
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

//0.6
//获取comp 对应的字段名 和 其值   排除 nil部分
func (c OrmConf) getCompColumnsValueNoNil(v reflect.Value) (columns []string, values []interface{}, err error) {
	columns = make([]string, 0)
	values = make([]interface{}, 0)

	t := v.Type()

	mappingColumns, err := c.getStructMappingColumns(t)
	if err != nil {
		return
	}

	for column, i := range mappingColumns {
		inter := getFieldInter(v.Field(i))

		if inter != nil {
			columns = append(columns, column)
			values = append(values, inter)
		}

	}
	return
}

//0.6
//获取comp 对应的字段名 和 其值   不排除 nil部分
func (c OrmConf) getCompAllColumnsValue(v reflect.Value) (columns []string, values []interface{}, err error) {
	columns = make([]string, 0)
	values = make([]interface{}, 0)

	t := v.Type()

	mappingColumns, err := c.getStructMappingColumns(t)
	if err != nil {
		return
	}

	for column, i := range mappingColumns {
		inter := getFieldInter(v.Field(i))
		columns = append(columns, column)
		values = append(values, inter)
	}
	return
}

//0.6
//获取comp 对应的字段名 和 其值   不排除 nil部分
func (c OrmConf) getCompAllColumnsValueList(v []reflect.Value) ([]string, [][]interface{}, error) {
	columns := make([]string, 0)
	values := make([][]interface{}, 0)

	mappingColumns, err := c.getStructMappingColumns(v[0].Type())
	if err != nil {
		return nil, nil, err
	}

	for column := range mappingColumns {
		columns = append(columns, column)
	}

	for _, value := range v {
		mappingColumns, err = c.getStructMappingColumns(value.Type())
		if err != nil {
			return nil, nil, err
		}

		vas := make([]interface{}, 0)
		for _, column := range columns {
			j := mappingColumns[column]
			inter := getFieldInter(value.Field(j))
			if inter == nil {
				inter = "default"
			}
			vas = append(vas, inter)
		}
		values = append(values, vas)
	}
	return columns, values, nil
}

func (c OrmConf) getColFieldIndexLinkMap(columns []string, t reflect.Type) (ColFieldIndexLinkMap, error) {
	if isSingleType(t) {
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

//tableName表名
//keys
//hasTen true开启多租户
func (c OrmConf) genDelSqlCommon(tableName string, keys []string, hasTen bool) []byte {
	var bb bytes.Buffer
	whereSql := c.GenWhere(keys, hasTen)

	logicDeleteSetSql := ormConfig.LogicDeleteSetSql
	logicDeleteYesSql := ormConfig.LogicDeleteYesSql
	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	logicDeleteYesSql = strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	if logicDeleteSetSql == lgSql {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.WriteString(string(whereSql))
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(lgSql)
		bb.WriteString(string(whereSql))
		bb.WriteString(" and ")
		bb.WriteString(logicDeleteYesSql)
	}
	return bb.Bytes()
}

//有tenantid功能
func (c OrmConf) GenWhere(keys []string, hasTen bool) []byte {
	if hasTen {
		keys = append(keys, ormConfig.TenantIdFieldName)
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

//tableName表名
//columns
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
