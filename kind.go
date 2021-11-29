package lorm

import (
	"database/sql/driver"
	"fmt"
	"github.com/lontten/lorm/types"
	"reflect"
)

//------------------base------------------
//	v0.6
//是否基本类型
func isBaseType(t reflect.Type) bool {
	kind := t.Kind()
	if reflect.Invalid < kind && kind < reflect.Array {
		return true
	}
	if reflect.String == kind {
		return true
	}
	return false
}
func isBaseValue(v reflect.Value) bool {
	kind := v.Kind()
	if reflect.Invalid < kind && kind < reflect.Array {
		return true
	}
	if reflect.String == kind {
		return true
	}
	return false
}

//-----------------------nuller---------------------------------
// v0.6
//检查是否nuller
func isNullerValue(v reflect.Value) bool {
	_, ok := v.Interface().(types.NullEr)
	return ok
}

func isNullerType(t reflect.Type) bool {
	return t.Implements(ImpNuller)
}

//----------------------valuer------------------------
//	v0.6
//是否valuer
func isValuerType(t reflect.Type) bool {
	if isBaseType(t) {
		return true
	}
	return t.Implements(ImpValuer)
}
func isValuerValue(v reflect.Value) bool {
	if isBaseValue(v) {
		return true
	}
	_, ok := v.Interface().(driver.Valuer)
	return ok
}

//----------------------struct-valuer------------------------
//	v0.6
//是否struct类型 valuer ,
//是否struct，是否有valuer
func baseStructValueValuer(v reflect.Value) (bool, bool) {
	is := isValuerValue(v)
	if !is {
		return false, false
	}
	_, ok := v.Interface().(driver.Valuer)
	return is, ok
}
func baseStructTypeValuer(t reflect.Type) (bool, bool) {
	is := isValuerType(t)
	if !is {
		return false, false
	}
	return is, t.Implements(ImpValuer)
}

//-------------ptr-----------------
// v0.6
//是指针类型，返回指针的基类型
func basePtrType(t reflect.Type) (bool, reflect.Type) {
	if t.Kind() == reflect.Ptr {
		return true, t.Elem()
	}
	return false, t
}
func basePtrValue(v reflect.Value) (is bool, structType reflect.Value) {
	if v.Kind() == reflect.Ptr {
		return true, v.Elem()
	}
	return false, v
}

//-------------ptr-deep-------------------
// v0.6
//是指针类型，返回指针的最基类型
func basePtrDeepType(t reflect.Type) (bool, reflect.Type) {
	isPtr := false
base:
	fmt.Println(t.Kind())
	fmt.Println(t.String())
	if t.Kind() == reflect.Ptr {
		isPtr = true
		t = t.Elem()
		goto base
	}
	return isPtr, t
}
func basePtrDeepValue(v reflect.Value) (bool, reflect.Value) {
	isPtr := false
base:
	if v.Kind() == reflect.Ptr {
		isPtr = true
		v = v.Elem()
		goto base
	}
	return isPtr, v
}

//-----------------slice---------------
// v0.6
//是数组类型，返回数组的基类型
func baseSliceType(t reflect.Type) (bool, reflect.Type) {
	if t.Kind() == reflect.Slice {
		return true, t.Elem()
	}
	return false, t
}

// v0.6
// is 是否 slice ; has 是否有内容
func baseSliceValue(v reflect.Value) (is, has bool, structType reflect.Value) {
	if v.Kind() == reflect.Slice {
		return true, v.Len() != 0, v.Index(0)
	}
	return false, false, v
}

//---------------------slice-deep------------------
// v0.6
//是数组类型，返回数组的最基类型
func baseSliceDeepType(t reflect.Type) (ok bool, structType reflect.Type) {
	isSlice := false
	flag := true //base
	for flag {
		is, t := basePtrDeepType(t)
		if is {
			flag = false
		}

		is, t = _baseSliceDeepType(t)
		if is {
			flag = false
		}
		if flag {
			return isSlice, t
		}
		isSlice = true
	}
	return false, t
}

// v0.6
func baseSliceDeepValue(v reflect.Value) (bool, reflect.Value) {
	isSlice := false
	flag := true //base
	for flag {
		is, v := basePtrDeepValue(v)
		if is {
			flag = false
		}

		is, v = _baseSliceDeepValue(v)
		if is {
			flag = false
		}
		if flag {
			return isSlice, v
		}
		isSlice = true
	}
	return false, v
}

// v0.6
//是数组类型，返回数组的最基类型
func _baseSliceDeepType(t reflect.Type) (ok bool, structType reflect.Type) {
	isSlice := false
base:
	is, t := baseSliceType(t)
	if is {
		isSlice = true
		goto base
	}
	return isSlice, t
}
func _baseSliceDeepValue(v reflect.Value) (bool, reflect.Value) {
	isSlice := false
base:
	is, has, v := baseSliceValue(v)
	if is && has {
		isSlice = true
		goto base
	}
	return isSlice, v
}

//--------------------
// v0.6
//检查是 ptr 还是 slice
type packType int

const (
	None  packType = iota
	Ptr   packType = iota
	Slice packType = iota
)

// v0.6
//检查是否是ptr，slice类型
func checkPackType(t reflect.Type) (packType, reflect.Type) {
	is, base := basePtrDeepType(t)
	if is {
		return Ptr, base
	}
	is, base = baseSliceDeepType(base)
	if is {
		return Slice, base
	}
	return None, t
}

func checkPackValue(v reflect.Value) (packType, reflect.Value) {
	is, base := basePtrDeepValue(v)
	if is {
		return Ptr, base
	}
	is, base = baseSliceDeepValue(base)
	if is {
		return Slice, base
	}
	return None, v
}

//-----------------map-------
// v0.6
// is 是否 slice has 是否有内容
func baseMapValue(v reflect.Value) (is, has bool, key reflect.Value) {
	if v.Kind() == reflect.Map {
		if v.Len() == 0 {
			return true, false, v
		}
		key = v.MapKeys()[0]
		return true, true, key
	}
	return false, false, v
}

// v0.6
// is 是否 slice has 是否有内容
func baseMapType(t reflect.Type) (is, has bool) {
	if t.Kind() != reflect.Map {
		return false, false
	}
	if t.Len() == 0 {
		return true, false
	}

	return true, true
}

//--------------------------------------------------

//--------------------
// v0.6
//检查是 single 还是 composite
type compType int

const (
	Invade     compType = iota
	SliceEmpty compType = iota
	Single     compType = iota
	Composite  compType = iota
)

func checkCompValue(v reflect.Value, canEmpty bool) compType {
	is := isBaseValue(v)
	if is {
		return Single
	}
	is, has, _ := baseMapValue(v)
	if is {
		if canEmpty || has {
			return Composite
		}
		return SliceEmpty
	}

	isStruct, isValuer := baseStructValueValuer(v)
	if isStruct {
		if isValuer {
			return Single
		}
		return Composite
	}
	return Invade
}

func checkCompType(t reflect.Type) compType {
	is := isBaseType(t)
	if is {
		return Single
	}

	is, _ = baseMapType(t)
	if is {
		return Composite
	}

	isStruct, isValuer := baseStructTypeValuer(t)
	if isStruct {
		if isValuer {
			return Single
		}
		return Composite
	}
	return Invade
}

//是 comp类型的struct，返回true
func checkCompStructValue(v reflect.Value) bool {
	isStruct, isValuer := baseStructValueValuer(v)
	return isStruct && !isValuer
}

//是 comp类型的struct，返回true
func checkCompStructType(t reflect.Type) bool {
	isStruct, isValuer := baseStructTypeValuer(t)
	return isStruct && !isValuer
}

//-----------------------map--------------------------------
// v0.6
//检查map key是否string，value是否valuer
func checkValidMap(v reflect.Value) bool {
	keys := v.MapKeys()
	for _, key := range keys {
		//key string
		if key.Kind() != reflect.String {
			return false
		}

		//nuller
		value := v.MapIndex(key)
		_, base := checkPackValue(value)
		//base
		is := isBaseType(base.Type())
		if !is {
			//struct
			_, is = baseStructTypeValuer(base.Type())
			if !is {
				return false
			}
		}
	}
	return true
}

// v0.6
//检查map key是否string，value是否valuer,nuller
func checkValidMapValuer(v reflect.Value) bool {
	keys := v.MapKeys()
	for _, key := range keys {
		//key string
		if key.Kind() != reflect.String {
			return false
		}

		//nuller
		value := v.MapIndex(key)
		_, base := checkPackValue(value)
		//base
		is := isBaseType(base.Type())
		if !is {
			//struct valuer
			_, is = baseStructTypeValuer(base.Type())
			if !is {
				return false
			}
			//nuller
			is = isNullerValue(base)
			if !is {
				return false
			}
		}
	}
	return true
}
