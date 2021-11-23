package lorm

import (
	"database/sql/driver"
	"github.com/lontten/lorm/types"
	"reflect"
)

//------------------base------------------
//	v0.6
//是否基本类型
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

//----------------------struct------------------------
//	v0.6
//是否struct类型
func baseStructType(t reflect.Type) (bool, reflect.Type) {
	if t.Kind() == reflect.Struct {
		return true, t
	}
	return false, t
}
func baseStructValue(v reflect.Value) (bool, reflect.Value) {
	if v.Kind() == reflect.Struct {
		return true, v
	}
	return false, v
}

//----------------------struct-valuer------------------------
//	v0.6
//是否struct类型 valuer ,
//是否struct，是否有valuer
func baseStructValueValuer(v reflect.Value) (bool, bool, reflect.Value) {
	is, base := baseStructValue(v)
	if !is {
		return false, false, base
	}
	_, ok := base.Interface().(driver.Valuer)
	return is, ok, base
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
//检查是 ptr 还是 slice
type packType int

const (
	None  packType = iota
	Ptr   packType = iota
	Slice packType = iota
)

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

func checkPackTypeValue(v reflect.Value) (packType, reflect.Value) {
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

// is 是否 slice has 是否有内容
func baseMapType(t reflect.Type) (is, has bool) {
	if t.Kind() == reflect.Map {
		if t.Len() == 0 {
			return true, false
		}
		return true, true
	}
	return false, false
}

//--------------------------------------------------

//--------------------
//检查是 single 还是 composite
type compType int

const (
	Invade    compType = iota
	Single    compType = iota
	Composite compType = iota
)

func checkCompTypeValue(v reflect.Value, canEmpty bool) (compType, reflect.Value) {
	t := v.Type()
	is := baseBaseType(t)
	if is {
		return Single, v
	}
	is, has := baseMapType(t)
	if is {
		if canEmpty || has {
			return Composite, v
		}
		return Invade, v
	}

	isStruct, isValuer, base := baseStructValueValuer(v)
	if isStruct {
		if isValuer {
			return Single, base
		}
		return Composite, base
	}
	return Invade, v
}
//-----------------------nuller---------------------------------
//检查是否nuller
func checkBaseNuller(v reflect.Value) bool {
	_, ok := v.Interface().(types.NullEr)
	return ok
}
//-----------------------map--------------------------------
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
		_, base := checkPackTypeValue(value)
		//base
		is := baseBaseType(base.Type())
		if !is {
			//struct
			_, is, _ = baseStructValueValuer(base)
			if !is {
				return false
			}
		}
	}
	return true
}
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
		_, base := checkPackTypeValue(value)
		//base
		is := baseBaseType(base.Type())
		if !is {
			//struct valuer
			_, is, _ = baseStructValueValuer(base)
			if !is {
				return false
			}
			//nuller
			is = checkBaseNuller(base)
			if !is {
				return false
			}
		}
	}
	return true
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
