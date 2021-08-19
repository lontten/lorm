package lorm

import (
	"reflect"
)

//type
func baseBaseType(t reflect.Type) bool {
	kind := t.Kind()
	if reflect.Invalid < kind && kind < reflect.Array {
		return true
	}
	if reflect.String == kind {
		return true
	}
	return false
}

func baseSliceType(t reflect.Type) (ok bool, structType reflect.Type) {
	if t.Kind() == reflect.Slice {
		return true, t.Elem()
	}
	return false, t
}

func basePtrType(t reflect.Type) (is bool, structType reflect.Type) {
	if t.Kind() == reflect.Ptr {
		return true, t.Elem()
	}
	return false, t
}

func baseStructType(t reflect.Type) (is bool, structType reflect.Type) {
	if t.Kind() == reflect.Struct {
		return true, t
	}
	return false, t
}

//value

func baseBaseValue(v reflect.Value) bool {
	kind := v.Kind()
	if reflect.Invalid < kind && kind < reflect.Array {
		return true
	}
	if reflect.String == kind {
		return true
	}
	return false
}

// is 是否 slice has 是否有内容
func baseMapValue(v reflect.Value) (is, has bool, key, value reflect.Value) {
	if v.Kind() == reflect.Map {
		if v.Len() == 0 {
			return true, false, v, v
		}
		key = v.MapKeys()[0]
		return true, true, key, v.MapIndex(key)
	}
	return false, false, v, v
}

// is 是否 slice has 是否有内容
func baseSliceValue(v reflect.Value) (is, has bool, structType reflect.Value) {
	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return true, false, v
		}
		return true, true, v.Index(0)
	}
	return false, false, v
}

func basePtrValue(v reflect.Value) (is bool, structType reflect.Value) {
	if v.Kind() == reflect.Ptr {
		return true, v.Elem()
	}
	return false, v
}

func baseStructValue(v reflect.Value) (is bool, structType reflect.Value) {
	if v.Kind() == reflect.Struct {
		return true, v
	}
	return false, v
}

//struct

//必须为 struct或基础类型  的 指针类型
func basePtrStructBaseType(t reflect.Type) (is bool, structType reflect.Type) {
	is, base := basePtrType(t)
	if !is {
		return false, t
	}
	is, base = baseStructType(base)
	if is {
		return true, base
	}
	is = baseBaseType(base)
	if is {
		return true, base
	}
	return false, t
}

//struct
// no ptr -1
// struct 1
// base 2
// no type -2
func basePtrStructBaseValue(v reflect.Value) (t int, structType reflect.Value) {
	is, base := basePtrValue(v)
	if !is {
		return -1, v
	}
	is, base = baseStructValue(base)
	if is {
		return 1, base
	}
	is = baseBaseValue(base)
	if is {
		return 2, base
	}
	return -2, v
}

//struct
// struct 1
// base 2
// no type -2
func baseStructBaseValue(v reflect.Value) (int, reflect.Value) {
	is, base := baseStructValue(v)
	if is {
		return 1, base
	}
	is = baseBaseValue(base)
	if is {
		return 2, base
	}
	return -2, v
}

//struct
// struct 1
// base 2
// no type -2
func baseStructBaseType(t reflect.Type) (int, reflect.Type) {
	is, base := baseStructType(t)
	if is {
		return 1, base
	}
	is = baseBaseType(base)
	if is {
		return 2, base
	}
	return -2, t
}
