package lorm

import (
	"github.com/pkg/errors"
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

func baseSliceValue(v reflect.Value) (ok bool, structType reflect.Value) {
	if v.Kind() == reflect.Slice {
		return true, v.Elem()
	}
	return false, v
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


//struct
func checkScanTypeLn(t reflect.Type) (reflect.Type, error) {
	is, base := basePtrStructBaseType(t)
	if !is {
		return t, errors.New("need a ptr struct or base type")
	}
	return base, nil
}

func checkScanType(t reflect.Type) (reflect.Type, error) {
	_, base := basePtrType(t)
	is, base := baseSliceType(base)
	if !is {
		return t, errors.New("need a slice type")
	}

	is, base = basePtrStructBaseType(base)
	if !is {
		return t, errors.New("need a ptr struct or base type")
	}
	return base, nil
}
