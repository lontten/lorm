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










//--------------------struct-field-nuller---------
// v0.6 检查一个 struct 是否合法
var structValidCache = make(map[reflect.Type]reflect.Value)
var mutexStructValidCache sync.Mutex

//valuer
func checkStructValidField(va reflect.Value) error {
	mutexStructValidCache.Lock()
	defer mutexStructValidCache.Unlock()

	typ := va.Type()
	_, ok := structValidCache[typ]
	if ok {
		return nil
	}

	is := baseStructValue(va)
	if !is {
		return errors.New("need a struct")
	}

	err := _checkStructValidField(va)
	if err != nil {
		return err
	}

	structValidCache[typ] = va
	return nil
}

//v0.6
// 检查一个非 single struct 是否合法
func _checkStructValidField(v reflect.Value) error {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkField(v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

//--------------------struct-field-valuer-nuller---------
// v0.6 检查一个 struct 是否合法
var structValidNullerCache = make(map[reflect.Type]reflect.Value)
var mutexStructValidNullerCache sync.Mutex
//valuer
//nuller
func checkStructValidFieldNuller(va reflect.Value) error {
	mutexStructValidNullerCache.Lock()
	defer mutexStructValidNullerCache.Unlock()

	typ := va.Type()
	_, ok := structValidNullerCache[typ]
	if ok {
		return nil
	}

	is := baseStructValue(va)
	if !is {
		return errors.New("need a struct")
	}

	err := _checkStructValidFieldNuller(va)
	if err != nil {
		return err
	}

	structValidNullerCache[typ] = va
	return nil
}

//v0.6
// 检查一个非 single struct 是否合法
func _checkStructValidFieldNuller(v reflect.Value) error {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkFieldNuller(v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}



//--------------------comp-field-nuller---------
// v0.6 检查一个 struct 是否合法
var compValidCache = make(map[reflect.Type]reflect.Value)
var mutexCompValidCache sync.Mutex

func checkCompValidField(va reflect.Value) error {
	mutexCompValidCache.Lock()
	defer mutexCompValidCache.Unlock()

	typ := va.Type()
	_, ok := compValidCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()
	if kind==reflect.Struct {
		err := _checkStructValidField(va)
		if err != nil {
			return err
		}
		compValidCache[typ] = va
		return nil
	}
	if kind==reflect.Map {
		is := checkValidMap(va)
		if is {
			compValidCache[typ] = va
			return nil
		}
	}
	return errors.New("need a struct or map")
}

//v0.6
// 检查一个 comp 是否合法
func _checkCompValidField(v reflect.Value) error {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkField(v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

//--------------------comp-field-valuer-nuller---------
// v0.6 检查一个 comp 是否合法
var compValidNullerCache = make(map[reflect.Type]reflect.Value)
var mutexCompValidNullerCache sync.Mutex

func checkCompValidFieldNuller(va reflect.Value) error {
	mutexCompValidNullerCache.Lock()
	defer mutexCompValidNullerCache.Unlock()

	typ := va.Type()
	_, ok := compValidNullerCache[typ]
	if ok {
		return nil
	}

	kind := typ.Kind()
	if kind==reflect.Struct {
		err := _checkCompValidFieldNuller(va)
		if err != nil {
			return err
		}
		compValidNullerCache[typ] = va
		return nil
	}
	if kind==reflect.Map {
		is := checkValidMapValuer(va)
		if is {
			compValidNullerCache[typ] = va
			return nil
		}
	}
	return errors.New("need a struct or map")
}

//v0.6
// 检查一个 comp 是否合法
func _checkCompValidFieldNuller(v reflect.Value) error {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkFieldNuller(v.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}
