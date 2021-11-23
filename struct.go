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

//去除 所有 ptr slice 获取 struct ，不为struct 或 基础类型 为false
//1 ptr 2slice 3struct  0基础类型
//Deprecated
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

// *struct
//Deprecated
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

// v0.5
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

// v0.5 检查一个 struct 是否合法
var structValidCache = make(map[reflect.Type]reflect.Value)
var mutexStructValidCache sync.Mutex

func checkValidFieldTypStruct(va reflect.Value) error {
	mutexStructValidCache.Lock()
	defer mutexStructValidCache.Unlock()

	typ := va.Type()
	_, ok := structValidCache[typ]
	if ok {
		return nil
	}

	is, base := baseStructValue(va)
	if !is {
		return errors.New("need a struct")
	}

	err := _checkStructFieldValid(base)
	if err != nil {
		return err
	}

	structValidCache[typ] = base
	return nil
}

//v0.5
// 检查一个非 single struct 是否合法
func _checkStructFieldValid(v reflect.Value) error {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkFieldNuller(v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}
