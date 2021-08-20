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

func baseSliceDeepType(t reflect.Type) (ok bool, structType reflect.Type) {
base:
	is, base := baseSliceType(t)
	if is {
		t = base
		goto base
	}

	code, base := baseStructBaseType(base)
	if code < 0 {
		return false, t
	}
	return true, base
}

func baseSliceDeepValue(v reflect.Value) (is, has bool, base reflect.Value) {
	_, base = basePtrValue(v)
base:
	is, has, base = baseSliceValue(base)
	if is {
		if !has { //是slice 内容空
			return
		}
		goto base
	}

	code, base := baseStructBaseValue(base)
	if code < 0 {
		return false, true, base
	}
	return true, true, base
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
// slice 3
// no type -2
func baseStructBaseSliceValue(v reflect.Value) (int, reflect.Value) {
	is, base := baseStructValue(v)
	if is {
		return 1, base
	}
	is = baseBaseValue(base)
	if is {
		return 2, base
	}
	is,_, _ = baseSliceDeepValue(base)
	if is {
		return 3, base
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

//struct
// struct 1
// base 2
// slice 3
// no type -2
func baseStructBaseSliceType(t reflect.Type) (int, reflect.Type) {
	is, base := baseStructType(t)
	if is {
		return 1, base
	}
	is = baseBaseType(base)
	if is {
		return 2, base
	}
	is, base = baseSliceDeepType(base)
	if is {
		return 3, base
	}
	return -2, t
}
