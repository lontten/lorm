package lorm

import (
	"fmt"
	"reflect"
)

//------------------base------------------
//	v0.7
//是否基本类型
func _isBaseType(t reflect.Type) bool {
	kind := t.Kind()
	if reflect.Invalid < kind && kind < reflect.Array {
		return true
	}
	if reflect.String == kind {
		return true
	}
	return false
}

//------------------struct------------------
//	v0.7
//是否基本类型
func _isStructType(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

//-----------------map-------
// v0.7
// is 是否 slice has 是否有内容
func baseMapValue(v reflect.Value) (is, has bool, key reflect.Value) {
	if v.Kind() != reflect.Map {
		return false, false, v
	}
	if v.Len() == 0 {
		return true, false, v
	}
	key = v.MapKeys()[0]
	return true, true, key
}

//-----------------single-------
// v0.7
func isSingleType(t reflect.Type) bool {
	return checkCompType(t) == Single
}

//-----------------composite-------
// v0.7
func isCompType(t reflect.Type) bool {
	return checkCompType(t) == Composite
}

// v0.7
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

//-----------------------nuller---------------------------------
// v0.7
//检查是否nuller
func isNullerType(t reflect.Type) bool {
	return t.Implements(ImpNuller)
}

//----------------------valuer------------------------
//	v0.7
//是否valuer
func isValuerType(t reflect.Type) bool {
	if _isBaseType(t) {
		return true
	}
	return t.Implements(ImpValuer)
}

//----------------------struct-comp------------------------
//	v0.7
//是否struct类型comp struct-comp
func isStructCompValue(v reflect.Value) bool {
	is := _isStructType(v.Type())
	if !is {
		return false
	}
	typ := checkCompValue(v, false)
	return typ == Composite
}

func isStructCompType(t reflect.Type) bool {
	is := _isStructType(t)
	if !is {
		return false
	}
	typ := checkCompType(t)
	return typ == Composite
}

//-------------ptr-----------------
// v0.7
//是指针类型，返回指针的基类型
func basePtrType(t reflect.Type) (bool, reflect.Type) {
	if t.Kind() == reflect.Ptr {
		return true, t.Elem()
	}
	return false, t
}
func basePtrValue(v reflect.Value) (bool, reflect.Value, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return false, v, ErrNil
		}
		return true, v.Elem(), nil
	}
	return false, v, nil
}

//-------------ptr-deep-------------------
// v0.7
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

func basePtrDeepValue(v reflect.Value) (bool, reflect.Value, error) {
	isPtr := false
base:
	if v.Kind() == reflect.Ptr {
		isPtr = true
		if v.IsNil() {
			return isPtr, v, ErrNil
		}
		v = v.Elem()
		goto base
	}
	return isPtr, v, nil
}

//-----------------slice---------------
// v0.7
//是数组类型，返回数组的基类型
func baseSliceType(t reflect.Type) (bool, reflect.Type) {
	typ := checkCompType(t)
	if typ != Invade {
		return false, t
	}
	if t.Kind() == reflect.Slice {
		return true, t.Elem()
	}
	return false, t
}

// v0.7
// is 是否 slice
func baseSliceValue(v reflect.Value, canEmpty bool) (is bool, structType reflect.Value) {
	typ := checkCompValue(v, canEmpty)
	if typ != Invade {
		return false, v
	}

	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return canEmpty, v
		}
		return true, v.Index(0)
	}
	return false, v
}

//---------------------slice-deep------------------
// v0.7
//是数组类型，返回数组的最基类型
func baseSliceDeepType(t reflect.Type) (ok bool, structType reflect.Type) {
	isSlice := false
	tmp := t
	flag := true //base
	for flag {
		is, base := basePtrDeepType(tmp)
		if is {
			flag = false
		}

		is, base = _baseSliceDeepType(base)
		if is {
			flag = false
		}
		if flag {
			if isSlice {
				return true, base
			}
			return false, t
		}
		isSlice = true
		tmp = base
	}
	return false, t
}

// v0.6
func baseSliceDeepValue(v reflect.Value) (bool, reflect.Value, error) {
	isSlice := false
	tmp := v
	for true {
		flag := false //base
		fmt.Println(tmp.String())
		isPtr, base, err := basePtrDeepValue(tmp)
		if err != nil {
			return false, v, err
		}
		if isPtr {
			flag = true
		}

		is, base, err := _baseSliceDeepValue(base)
		if err != nil {
			return false, v, err
		}
		if is {
			flag = true
		}
		if !flag {
			if isSlice {
				return true, base, nil
			}
			return false, v, nil
		}
		isSlice = true
		tmp = base
	}
	return false, v, nil
}

// v0.7
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

//slice 用到value的一定是非空
func _baseSliceDeepValue(v reflect.Value) (bool, reflect.Value, error) {
	isSlice := false
base:
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Slice || kind == reflect.Map {
		if v.IsNil() {
			return false, v, ErrNil
		}
	}
	is, v := baseSliceValue(v, false)
	if is {
		isSlice = true
		goto base
	}
	return isSlice, v, nil
}

//--------------------
// v0.7
//检查是 ptr 还是 slice
type packType int

const (
	None packType = iota
	Ptr
	Slice
)

// v0.7
//检查是否是ptr，slice类型
func checkPackType(t reflect.Type) (packType, reflect.Type) {
	isPtr, base := basePtrDeepType(t)
	is, base := baseSliceDeepType(base)
	if is {
		return Slice, base
	}

	//不是slice，才判断是否ptr
	if isPtr {
		return Ptr, base
	}

	return None, t
}

type PackTyp struct {
	Typ       packType
	Base      reflect.Value
	SliceBase reflect.Value
}

// ptr base
// slice base
func checkPackValue(v reflect.Value) (PackTyp, error) {
	isPtr, v, err := basePtrDeepValue(v)

	packTyp := PackTyp{
		Typ:       None,
		Base:      v,
		SliceBase: v,
	}

	if err != nil {
		return packTyp, err
	}

	is, base, err := baseSliceDeepValue(v)

	if err != nil {
		return packTyp, err
	}
	if is {
		packTyp.Typ = Slice
		packTyp.SliceBase = base
		return packTyp, nil
	}

	//不是slice，才判断是否ptr
	if isPtr {
		packTyp.Typ = Ptr
		return packTyp, nil
	}

	return packTyp, nil
}

//--------------------------------------------------

//--------------------
// v0.7
//检查是 single 还是 composite
//base和实现valuer的是single，
//否则，并且类型是struct、map的是composite
//其他为invade
type compType int

const (
	Invade compType = iota
	SliceEmpty
	Single
	Composite
)

func checkCompValue(v reflect.Value, canEmpty bool) compType {
	is := isValuerType(v.Type())
	if is {
		return Single
	}
	is = _isStructType(v.Type())
	if is {
		return Composite
	}

	is, has, _ := baseMapValue(v)
	if is {
		if canEmpty || has {
			return Composite
		}
		return SliceEmpty
	}

	return Invade
}

func checkCompType(t reflect.Type) compType {
	is := isValuerType(t)
	if is {
		return Single
	}

	is = _isStructType(t)
	if is {
		return Composite
	}

	is, _ = baseMapType(t)
	if is {
		return Composite
	}
	return Invade
}

//-----------------------map--------------------------------
// v0.7
//scan不需要必须nuller
//检查map key是否string，value是否valuer
func checkMapFieldScan(v reflect.Value) bool {
	keys := v.MapKeys()
	for _, key := range keys {
		//key string
		if key.Kind() != reflect.String {
			return false
		}

		//valuer
		value := v.MapIndex(key)
		packTyp, err := checkPackValue(value)
		if err != nil {
			continue
		}
		typ := checkCompValue(packTyp.SliceBase, true)
		if typ != Single {
			return false
		}
	}
	return true
}

// v0.7
//检查map key是否string，value是否valuer/nuller
func checkMapField(v reflect.Value) bool {
	keys := v.MapKeys()
	for _, key := range keys {
		//key string
		if key.Kind() != reflect.String {
			return false
		}

		//valuer
		value := v.MapIndex(key)
		packTyp, err := checkPackValue(value)
		if err != nil {
			continue
		}
		base := packTyp.Base
		typ := checkCompValue(base, true)
		if typ != Single {
			return false
		}

		//nuller
		is := isNullerType(base.Type())
		if !is {
			return false
		}
	}
	return true
}

//-----------------------struct--------------------------------
// v0.7
//scan不需要必须nuller
//检查 struct field，value是否valuer
func checkStructFieldScan(v reflect.Value) bool {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkField(v.Field(i).Type())
		if err != nil {
			return false
		}
	}
	return true
}

// v0.7
//检查 struct field，value是否valuer/nuller
func checkStructField(v reflect.Value) bool {
	numField := v.NumField()
	for i := 0; i < numField; i++ {
		err := checkFieldNuller(v.Field(i).Type())
		if err != nil {
			return false
		}
	}
	return true
}
