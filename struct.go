package lorm

import (
	"errors"
	"github.com/lontten/lorm/utils"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// set
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

// ---------------struct-new-----------------
// v0.6
func newStruct(t reflect.Type) reflect.Value {
	tPtr := reflect.New(t)
	if isValuerType(t) {
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

// --------------------comp-field-valuer---------
// v03 检查一个 struct/map 是否合法,valuer
var compFieldVCache = make(map[reflect.Type]struct{})
var mutexCompFieldVCache sync.Mutex

func checkCompFieldV(typ reflect.Type) error {
	mutexCompFieldVCache.Lock()
	defer mutexCompFieldVCache.Unlock()

	_, ok := compFieldVCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		err := checkStructFieldV(typ)
		if err != nil {
			return err
		} else {
			compFieldVCache[typ] = struct{}{}
			return nil
		}
	}
	//map
	if kind == reflect.Map {
		err := checkMapFieldV(typ)
		if err != nil {
			return err
		} else {
			compFieldVCache[typ] = struct{}{}
			return nil
		}
	}
	return errors.New("need a struct or map-scan")
}

// --------------------comp-field-valuer-nuller---------
// v0.7 检查一个 comp 是否合法
var compFieldVNCache = make(map[reflect.Type]reflect.Value)
var mutexCompFieldVNCache sync.Mutex

func checkCompFieldVN(va reflect.Value) error {
	mutexCompFieldVNCache.Lock()
	defer mutexCompFieldVNCache.Unlock()

	typ := va.Type()
	_, ok := compFieldVNCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		err := checkStructFieldVN(typ)
		if err == nil {
			compFieldVNCache[typ] = va
			return nil
		}
		return err
	}
	//map
	if kind == reflect.Map {
		is := checkMapFieldValue(va)
		if is {
			compFieldVNCache[typ] = va
			return nil
		}
	}
	return errors.New("checkCompFieldVN err;need a struct or map")
}

// todo 下面未重构--------------
