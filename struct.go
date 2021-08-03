package lorm

import (
	"database/sql/driver"
	"errors"
	"github.com/lontten/lorm/types"
	"github.com/lontten/lorm/utils"
	"reflect"
	"strconv"
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

func checkStructValidField(v reflect.Value) bool {
	code, _ := basePtrStructBaseValue(v)
	if code > 0 {
		return true
	}
	if v.Kind() == reflect.Struct {
		_, ok := v.Interface().(driver.Valuer)
		if !ok {
			return false
		}
		_, ok = v.Interface().(types.NullEr)
		if !ok {
			return false
		}
		return true
	}
	if v.Kind() == reflect.Slice {
		return true
	}
	return false
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

//用于检查，单一值的合法性，base 或 valuer struct
// bool true 代表有效 false:无效-nil
// err 不合法
func checkValidFieldTypOne(v reflect.Value) (bool, error) {

	is, base := basePtrValue(v)
	if is {
		isNil := v.IsNil()
		if isNil { //数值无效，直接返回false，不再进行合法性检查
			return false, nil
		}
	}

	is = baseBaseValue(base)
	if is {
		return true, nil
	}

	is, base = baseStructValue(base)
	if is {
		_, ok := base.Interface().(driver.Valuer)
		if !ok {
			return false, errors.New("struct field " + base.String() + " need imp sql Value")
		}
		_, ok = base.Interface().(types.NullEr)
		if !ok {
			return false, errors.New("struct field " + base.String() + " need imp core NullEr ")
		}

		return true, nil
	}

	return false, errors.New("need a struct or base type")
}

// string[] - ptr				1
// map[string]intface - ptr		2
// struct - ptr					3
// base struct-base- ptr		4
func checkArgTyp(v reflect.Value) (int, error) {
	is, base := basePtrValue(v)
	if is {
		isNil := v.IsNil()
		if isNil { //数值无效，直接返回false，不再进行合法性检查
			return 0, errors.New("  is nil")
		}
	}

	//slice string
	is, base, err := baseSliceValue(base)
	if err != nil {
		return 0, err
	}
	if is {
		if base.Kind() != reflect.String {
			return 0, errors.New("  type err slice " + base.Kind().String())
		}
		return 1, nil
	}

	//map string
	is, key, value, err := baseMapValue(base)
	if err != nil {
		return 0, err
	}
	if is {
		if key.Kind() != reflect.String {
			return 0, errors.New(" map type key err no string  ")
		}
		if value.Kind() != reflect.Interface {
			return 0, errors.New(" map type value err no interface  ")
		}
		return 2, nil
	}

	// base
	is = baseBaseValue(base)
	if is {
		return 4, nil
	}

	is, base = baseStructValue(base)
	if !is {
		return 0, errors.New("  type err   " + base.Kind().String())
	}
	//struct-base
	_, ok := base.Interface().(driver.Valuer)
	if ok {
		return 4, nil
	}

	numField := base.NumField()
	for i := 0; i < numField; i++ {
		field := base.Field(i)
		ok := checkStructValidField(field)
		if !ok {
			return 0,errors.New("  type err struct filed  " )
		}
	}
	return 3,nil
}

//用于检查，id  base 或  struct field nuller
// bool true 代表有效 false:无效-nil
// err 不合法
func checkValidPrimaryKey(v []interface{}, ids []string) ([]interface{}, error) {
	singlePk := len(ids) == 1
	arr := make([]interface{}, 0)
	for i, e := range v {
		value := reflect.ValueOf(e)
		is, base := basePtrValue(value)
		if is {
			isNil := value.IsNil()
			if isNil { //数值无效，直接返回false，不再进行合法性检查
				return nil, errors.New("PrimaryKey " + strconv.Itoa(i) + ": is nil")
			}
		}
		is, err := checkValidFieldTypOne(base)
		if err != nil {
			return nil, err
		}
		if is {
			if !singlePk {
				return nil, errors.New("PrimaryKey arg is err")
			}
			arr = append(arr, e)
			continue
		}
		is, base = baseStructValue(base)
		if !is {
			return nil, errors.New("need a struct or base type")
		}
		err = checkValidFieldTypStruct(base)
		if err != nil {
			return nil, err
		}

		if singlePk {
			return nil, errors.New("PrimaryKey arg is err")
		}
		for _, id := range ids {
			field := base.FieldByName(id)
			arr = append(arr, field.Interface())
		}
	}
	return arr, nil
}

var structValidCache = make(map[reflect.Type]error)
var mutexStructValidCache sync.Mutex

func checkValidFieldTypStruct(va reflect.Value) error {
	mutexStructValidCache.Lock()
	defer mutexStructValidCache.Unlock()

	_, base := basePtrValue(va)

	typ := base.Type()
	b, ok := structValidCache[typ]
	if ok {
		return b
	}

	is, base := baseStructValue(base)
	if !is {
		return errors.New("need a struct")
	}

	numField := base.NumField()
	for i := 0; i < numField; i++ {
		field := base.Field(i)

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
				return errors.New("struct field " + field.String() + " need imp core NullEr ")
			}

		}
	}
	structValidCache[typ] = nil
	return nil
}
