package lorm

import (
	"errors"
	"github.com/lontten/lorm/utils"
	"reflect"
	"sync"
	"unicode"
)

type StructValidFieldValueMap map[string]any

// ---------------struct-new-----------------
/**
根据 反射type，创建一个 struct,并返回 引用
*/
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

// --------------------comp-field-valuer-scanner---------
// 检查一个 comp 是否合法
var compFieldVSCache = make(map[reflect.Type]reflect.Value)
var mutexCompFieldVSCache sync.Mutex

func checkCompFieldVS(va reflect.Value) error {
	mutexCompFieldVSCache.Lock()
	defer mutexCompFieldVSCache.Unlock()

	typ := va.Type()
	_, ok := compFieldVSCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()

	//struct
	if kind == reflect.Struct {
		err := checkStructFieldVS(typ)
		if err == nil {
			compFieldVSCache[typ] = va
			return nil
		}
		return err
	}
	//map
	if kind == reflect.Map {
		is := checkMapFieldValue(va)
		if is {
			compFieldVSCache[typ] = va
			return nil
		}
	}
	return errors.New("checkCompFieldVS err;need a struct or map")
}

var colName2fieldNameAllMapCache = make(map[reflect.Type]colName2fieldNameMap)

var colName2fieldNameAllMapMutex sync.Mutex

//	 可以缓存
//
//		主键ID，转化为id
//
// tag== lrom:-  跳过
// 过滤掉首字母小写的字段
// 获取model对应的数据字段名：和其在model中的字段名
func getStructColName2fieldNameAllMap(t reflect.Type) colName2fieldNameMap {
	fields, ok := colName2fieldNameAllMapCache[t]
	if ok {
		return fields
	}
	colName2fieldNameAllMapMutex.Lock()
	defer colName2fieldNameAllMapMutex.Unlock()
	fields, ok = colName2fieldNameAllMapCache[t]
	if ok {
		return fields
	}

	cfMap := colName2fieldNameMap{}

	numField := t.NumField()
	for i := 0; i < numField; i++ {
		structField := t.Field(i)
		// 跳过软删除字段
		if utils.IsSoftDelFieldType(structField.Type) {
			continue
		}

		name := structField.Name
		// 过滤掉首字母小写的字段
		if unicode.IsLower([]rune(name)[0]) {
			continue
		}

		tag := structField.Tag.Get("lorm")
		if tag == "-" {
			continue
		}

		if name == "ID" {
			cfMap["id"] = "ID"
			continue
		}

		if tag != "" {
			cfMap[tag] = name
			continue
		}

		cfMap[utils.Camel2Case(name)] = name
	}

	return cfMap
}
