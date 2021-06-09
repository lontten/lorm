package lorm

import (
	"database/sql/driver"
	"errors"
	"github.com/lontten/lorm/utils"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

//获取struct对应数据库表名
func getStructTableName(dest interface{}, config OrmConfig) (string, error) {
	typ := reflect.TypeOf(dest)
	base, err := baseSliceType(typ)
	if err != nil {
		return "", err
	}
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
	tableNamePrefix := config.TableNamePrefix
	if tableNamePrefix != "" {
		return tableNamePrefix + name, nil
	}

	return name, nil
}

//获取struct对应的字段名 有效部分
func getStructMappingColumns(t reflect.Type, config OrmConfig) (map[string]int, error) {
	cMap := make(map[string]int)

	numField := t.NumField()
	var num = 0
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		name := field.Name

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
func getStructMappingColumnsValue(dest interface{}, config OrmConfig) (columns []string, values []interface{}, err error) {
	columns = make([]string, 0)
	values = make([]interface{}, 0)

	t := reflect.TypeOf(dest)
	base, err := baseStructType(t)
	if err != nil {
		return
	}

	mappingColumns, err := getStructMappingColumns(base, config)
	if err != nil {
		return
	}

	v := reflect.ValueOf(dest)
	structValue, err := baseStructValue(v)
	if err != nil {
		return
	}

	for column, i := range mappingColumns {
		field := structValue.Field(i)
		indirect := reflect.Indirect(field)
		if field.Kind() == reflect.Struct {
			columns = append(columns, column)
			value, err := indirect.Interface().(driver.Valuer).Value()
			if err != nil {
				return nil, nil, err
			}

			values = append(values, value)
		} else {
			if !field.IsNil() {
				columns = append(columns, column)
				values = append(values, indirect.Interface())
			}
		}

	}
	return
}

//set
var structFieldsMapCache = make(map[reflect.Type]fieldMap)

type fieldMap map[string]int

var mutex sync.Mutex

func getFieldMap(typ reflect.Type) (fieldMap, error) {
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

		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}

		name = utils.Camel2Case(name)
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

	numField := va.NumField()
	for i := 0; i < numField; i++ {
		field := va.Field(i)
		validField, ok := baseStructValidField(field)
		if !ok {
			continue
		}
		_, ok = validField.Interface().(driver.Valuer)
		if !ok {
			return errors.New("struct field " + field.String() + " nedd imp sql Value")
		}
	}
	structValidCache[typ] = nil
	return nil
}

func baseStructType(t reflect.Type) (structType reflect.Type, e error) {
base:
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		goto base
	case reflect.Struct:
	default:
		return nil, errors.New("is not a struct type")
	}
	return t, nil
}

func baseStructValue(v reflect.Value) (structValue reflect.Value, e error) {
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		fallthrough
	case reflect.Struct:
	default:
		return v, errors.New("is not a struct value")
	}
	return v, nil
}

func baseStructValidField(v reflect.Value) (structValue reflect.Value, b bool) {
base:
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
		goto base
	case reflect.Slice:
		v = v.Elem()
		goto base
	case reflect.Struct:
	default:
		return v, false
	}
	return v, true
}

func baseSliceType(t reflect.Type) (structType reflect.Type, e error) {
	switch t.Kind() {
	case reflect.Ptr:
		t = t.Elem()
		fallthrough
	case reflect.Slice:
	default:
		return nil, errors.New("is not a slice type")
	}
	return t, nil
}

func baseSlicePtrType(t reflect.Type) (typ int, structType reflect.Type, e error) {
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
	default:
		return 0, nil, errors.New("is base not a  type")
	}
	return typ, t, nil
}

func newStruct(structTyp reflect.Type) reflect.Value {
	if structTyp.Kind() == reflect.Ptr {
		structTyp = structTyp.Elem()
	}
	tPtr := reflect.New(structTyp)
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
