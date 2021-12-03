package lorm

import (
	"errors"
	"github.com/lontten/lorm/utils"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

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

//---------------struct-new-----------------
// v0.6
func newStruct(t reflect.Type) reflect.Value {
	tPtr := reflect.New(t)
	if isSingleType(t) {
		return tPtr
	}
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Ptr {
			f := reflect.New(field.Type.Elem())
			tPtr.Elem().Field(i).Set(f)
		}
	}
	return tPtr
}

//--------------------comp-field-valuer---------
// v0.6 检查一个 struct 是否合法
var compValidCache = make(map[reflect.Type]struct{})
var mutexCompValidCache sync.Mutex

func checkCompFieldScan(typ reflect.Type) error {
	mutexCompValidCache.Lock()
	defer mutexCompValidCache.Unlock()

	_, ok := compValidCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		is := checkStructFieldType(typ)
		if is {
			compValidCache[typ] = struct{}{}
			return nil
		}
	}
	//map
	if kind == reflect.Map {
		is := checkMapFieldType(typ)
		if is {
			compValidCache[typ] = struct{}{}
			return nil
		}
	}
	return errors.New("need a struct or map-scan")
}

//--------------------comp-field-valuer-nuller---------
// v0.7 检查一个 comp 是否合法
var compValidNullerCache = make(map[reflect.Type]reflect.Value)
var mutexCompValidNullerCache sync.Mutex

func checkCompField(va reflect.Value) error {
	mutexCompValidNullerCache.Lock()
	defer mutexCompValidNullerCache.Unlock()

	typ := va.Type()
	_, ok := compValidNullerCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		is := checkStructFieldValue(va)
		if is {
			compValidNullerCache[typ] = va
			return nil
		}
	}
	//map
	if kind == reflect.Map {
		is := checkMapFieldValue(va)
		if is {
			compValidNullerCache[typ] = va
			return nil
		}
	}
	return errors.New("need a struct or map")
}
