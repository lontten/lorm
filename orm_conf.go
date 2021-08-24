package lorm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/lontten/lorm/types"
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

func (c OrmConf) ScanLn(rows *sql.Rows, v interface{}) (num int64, err error) {
	defer rows.Close()
	value := reflect.ValueOf(v)
	code, base := basePtrStructBaseValue(value)
	if code == -1 {
		return 0, errors.New("dest need a  ptr")
	}
	if code == -2 {
		return 0, errors.New("need a ptr struct or base type")
	}

	num = 1
	t := base.Type()

	columns, err := rows.Columns()
	if err != nil {
		return
	}
	cfm, err := c.getColFieldIndexLinkMap(columns, t)
	if err != nil {
		return
	}
	if rows.Next() {
		box, _, v := createColBox(t, cfm)
		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return
		}
		base.Set(v)
	}

	if rows.Next() {
		return 0, errors.New("result to many for one")
	}
	return
}

func (c OrmConf) Scan(rows *sql.Rows, v interface{}) (int64, error) {
	defer rows.Close()
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return 0, errors.New("need a ptr type")
	}
	arr := value.Elem()
	if arr.Kind() != reflect.Slice {
		return 0, errors.New("need a slice type")
	}

	slice := arr.Type()

	base := slice.Elem()
	isPtr := base.Kind() == reflect.Ptr
	code, base := baseStructBaseType(base)
	if code == -2 {
		return 0, errors.New("need a struct or base type in  slice")
	}

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	cfm, err := c.getColFieldIndexLinkMap(columns, base)
	fmt.Println(len(cfm))
	fmt.Println("------")
	if err != nil {
		return 0, err
	}
	var num int64 = 0
	for rows.Next() {
		box, vp, v := createColBox(base, cfm)

		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
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

func (c OrmConf) tableName(v reflect.Value) (string, error) {
	base := v.Type()

	// fun
	name := base.String()
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[index+1:]
	}
	name = utils.Camel2Case(name)

	tableNameFun := c.TableNameFun
	if tableNameFun != nil {
		return tableNameFun(name, base), nil
	}

	// tag

	numField := base.NumField()
	tagTableName := ""
	for i := 0; i < numField; i++ {
		if tag := base.Field(i).Tag.Get("tableName"); tag != "" {
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

func (c OrmConf) primaryKeys(tableName string, v reflect.Value) []string {
	primaryKeyNameFun := c.PrimaryKeyNameFun
	if primaryKeyNameFun != nil {
		return primaryKeyNameFun(tableName, v)
	}

	primaryKeyName := c.PrimaryKeyNames
	if len(primaryKeyName) != 0 {
		return primaryKeyName
	}

	//todo 获取 struct 中 tag为id 的 filed ，为 primaryKeyNames 可多个

	return []string{"id"}
}
func (c OrmConf) initColumns(v reflect.Value) (columns []string, err error) {
	typ := v.Type()
	base, err := baseStructTypePtr(typ)
	if err != nil {
		return
	}

	cMap := make(map[string]int)

	numField := base.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := base.Field(i)
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

func (c OrmConf) initColumnsValue(v reflect.Value) (columns []string, values []interface{}, err error) {
	columns = make([]string, 0)
	values = make([]interface{}, 0)

	t := v.Type()

	mappingColumns, err := c.getStructMappingColumns(t)
	if err != nil {
		return
	}

	for column, i := range mappingColumns {
		field := v.Field(i)

		typ, validField, ok := baseStructValidField(field)
		if !ok {
			return nil, nil, errors.New("struct field " + field.String() + " need field is ptr slice struct")
		}

		if typ == 0 {
			columns = append(columns, column)
			values = append(values, validField.Interface())
		}

		if typ == 1 || typ == 2 {
			if !field.IsNil() {
				columns = append(columns, column)
				values = append(values, validField.Interface())
			}
		}

		if typ == 3 {
			vv := validField.Interface().(types.NullEr)
			if !vv.IsNull() {
				value, _ := validField.Interface().(driver.Valuer).Value()
				columns = append(columns, column)
				values = append(values, value)
			}
		}
	}
	return

}

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

//获取struct对应的字段名 和 其值   有效部分
func (c OrmConf) getStructMappingColumnsValueNotNil(v reflect.Value) (columns []string, values []interface{}, err error) {
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

//获取struct对应的字段名 和 其值   包含 value nil 的部分
func (c OrmConf) getStructMappingColumnsValueList(v []reflect.Value) ([]string, [][]interface{}, error) {
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
			vas = append(vas, inter)
		}
		values = append(values, vas)
	}
	return columns, values, nil
}

func (c OrmConf) getColFieldIndexLinkMap(columns []string, typ reflect.Type) (ColFieldIndexLinkMap, error) {
	is := baseBaseType(typ)
	if is {
		return ColFieldIndexLinkMap{}, nil
	}

	colNum := len(columns)
	cfm := make([]int, colNum)
	fm, err := getFieldMap(typ, c.FieldNamePrefix)
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
