package lorm

import (
	"database/sql/driver"
	"errors"
	"github.com/lontten/lorm/types"
	"github.com/lontten/lorm/utils"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

//获取struct对应数据库表名
func getStructTableName(v reflect.Value, config OrmConf) (string, error) {
	base := v.Type()

	// fun
	name := base.String()
	index := strings.LastIndex(name, ".")
	if index > 0 {
		name = name[index+1:]
	}
	name = utils.Camel2Case(name)

	tableNameFun := config.TableNameFun
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
	tableNamePrefix := config.TableNamePrefix
	if tableNamePrefix != "" {
		return tableNamePrefix + name, nil
	}

	return name, nil
}

//获取struct对应的字段名 有效部分
func getStructMappingColumns(t reflect.Type, config OrmConf) (map[string]int, error) {
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

		if tag := field.Tag.Get("lorm"); tag == "-" {
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

		fieldNamePrefix := config.FieldNamePrefix
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
func getStructMappingColumnsValueNotNull(v reflect.Value, config OrmConf) (columns []string, values []interface{}, err error) {
	columns = make([]string, 0)
	values = make([]interface{}, 0)

	t := v.Type()

	mappingColumns, err := getStructMappingColumns(t, config)
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

//set
var structFieldsMapCache = make(map[reflect.Type]fieldMap)

type fieldMap map[string]int

var mutex sync.Mutex

func getFieldMap(typ reflect.Type, fieldNamePrefix string) (fieldMap, error) {
	mutex.Lock()
	defer mutex.Unlock()
	fields, ok := structFieldsMapCache[typ]
	if ok {
		return fields, nil
	}
	numField := typ.NumField()
	arr := fieldMap{}
	var num = 0
	for i := 0; i < numField; i++ {
		field := typ.Field(i)
		name := field.Name

		if name == "ID" {
			arr["id"] = i
			num++
			if len(arr) < num {
				return arr, errors.New("字段:: id error")
			}
			continue
		}

		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}

		name = utils.Camel2Case(name)
		name = strings.TrimPrefix(name, fieldNamePrefix)
		if tag := field.Tag.Get("db"); tag != "" {
			name = tag
		}
		arr[name] = i
		num++
		if len(arr) < num {
			return arr, errors.New("字段::" + name + "error")
		}
	}

	structFieldsMapCache[typ] = arr
	return arr, nil
}

type StructValidFieldValueMap map[string]interface{}

var structValidCache = make(map[reflect.Type]error)
var mutexStructValidCache sync.Mutex

func checkValidStruct(va reflect.Value) error {
	mutexStructValidCache.Lock()
	defer mutexStructValidCache.Unlock()

	typ := va.Type()

	b, ok := structValidCache[typ]
	if ok {
		return b
	}
	value, err := baseStructValuePtr(va)
	if err != nil {
		return err
	}

	numField := value.NumField()
	for i := 0; i < numField; i++ {
		field := value.Field(i)
		typ, validField, ok := baseStructValidField(field)
		if !ok {
			return errors.New("struct field " + field.String() + " need field is ptr slice struct")
		}
		//为 struct类型
		if typ == 3 {
			_, ok = validField.Interface().(driver.Valuer)
			if !ok {
				return errors.New("struct field " + field.String() + " need imp sql Value")
			}
			_, ok = validField.Interface().(types.NullEr)
			if !ok {
				return errors.New("struct field " + field.String() + " need imp lorm NullEr ")
			}

		}
	}
	structValidCache[typ] = nil
	return nil
}

//去除 所有 ptr slice 获取 struct ，不为struct 或 基础类型 为false
//1 ptr 2slice 3struct  0基础类型
func baseStructValidField(v reflect.Value) (typ int, structValue reflect.Value, b bool) {
	structValue = v
	t := v.Type()

base:
	switch t.Kind() {
	case reflect.Ptr:
		if typ == 0 {
			typ = 1
		}
		t = t.Elem()
		goto base
	case reflect.Slice:
		if typ == 0 {
			typ = 2
		}
		t = t.Elem()
		goto base
	case reflect.Struct:
		if typ == 0 {
			typ = 3
		}
		return typ, structValue, true
	case reflect.Map:
		return
	case reflect.Interface:
		return
	case reflect.Func:
		return
	case reflect.Invalid:
		return
	case reflect.UnsafePointer:
		return
	case reflect.Uintptr:
		return
	default:
		return typ, structValue, true
	}
}

//1 ptr 2slice 3struct  0基础类型
func baseStructValidField2(v reflect.Value) (typ int, structValue reflect.Value, b bool) {
	structValue = v
	t := v.Type()

base:
	switch t.Kind() {
	case reflect.Ptr:
		if typ == 0 {
			typ = 1
		}
		t = t.Elem()
		goto base
	case reflect.Slice:
		if typ == 0 {
			typ = 2
		}
		t = t.Elem()
		goto base
	case reflect.Struct:
		if typ == 0 {
			typ = 3
		}
		return typ, structValue, true
	case reflect.Map:
		return
	case reflect.Interface:
		return
	case reflect.Func:
		return
	case reflect.Invalid:
		return
	case reflect.UnsafePointer:
		return
	case reflect.Uintptr:
		return
	default:
		return typ, structValue, true
	}
}

// *struct
func baseStructTypePtr(t reflect.Type) (structType reflect.Type, e error) {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
	case reflect.Struct:
		return t, nil
	default:
		return nil, errors.New("is not a struct or ptr struct type")
	}

	switch t.Kind() {
	case reflect.Struct:
		return t, nil
	default:
		return nil, errors.New("is not a struct or ptr struct type")
	}

}

// *struct
func baseStructValuePtr(v reflect.Value) (structValue reflect.Value, e error) {
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	case reflect.Struct:
		return v, nil
	default:
		return v, errors.New("is not a struct or ptr struct value")
	}

	switch v.Kind() {
	case reflect.Struct:
		return v, nil
	default:
		return v, errors.New("is not a struct or ptr struct value")
	}
}

//把 *slice  获取 slice
func baseSliceTypePtr(t reflect.Type) (structType reflect.Type, e error) {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
	case reflect.Slice:
	default:
		return nil, errors.New("is not a slice type")
	}

	switch t.Kind() {
	case reflect.Slice:
	default:
		return nil, errors.New("is not a slice type")
	}
	return t, nil
}

//	只能是 struct *struct []struct 三种类型

//把 *[]struct 类型 剔除 * [] 获取 struct 的基础类型
//typ 最表面的类型 1 ptr  ；   2 slice  ；  0 struct

func baseStructTypeSliceOrPtr(t reflect.Type) (typ int, structType reflect.Type, e error) {
	switch t.Kind() {
	case reflect.Ptr:
		typ = 1
		t = t.Elem()
	case reflect.Slice:
		typ = 2
		t = t.Elem()
	case reflect.Struct:
		return typ, t, nil
	default:
		return 0, nil, errors.New("is base not a ptr slice struct type")
	}
	switch t.Kind() {
	case reflect.Struct:
		return typ, t, nil
	default:
		return 0, nil, errors.New("is base not a ptr slice struct type")
	}
}

func newStruct(structTyp reflect.Type) reflect.Value {
	tPtr := reflect.New(structTyp)
	if baseBaseType(structTyp) {
		return tPtr
	}
	numField := structTyp.NumField()
	for i := 0; i < numField; i++ {
		field := structTyp.Field(i)
		if field.Type.Kind() == reflect.Ptr {
			f := reflect.New(field.Type.Elem())
			tPtr.Elem().Field(i).Set(f)
		}
	}
	return tPtr
}
