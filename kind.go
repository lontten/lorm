package lorm

import (
	"errors"
	"reflect"
)

//todo 下面未重构--------------

// atom 原子类型			#作为字段使用。
// composite 非原子类型		#多个atom组成的实体类

// ------------------base------------------
// 是否基本类型
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

// ------------------struct------------------
// 是否struct类型
func _isStructType(t reflect.Type) bool {
	return t.Kind() == reflect.Struct
}

// -----------------map-------
// is 是否 slice has 是否有内容
func baseMapValue(v reflect.Value) (is bool, key reflect.Value) {
	if v.Kind() != reflect.Map {
		return false, v
	}
	if v.Len() == 0 {
		return false, v
	}
	key = v.MapKeys()[0]
	return true, key
}

// -----------------atom-------
func isAtomType(t reflect.Type) bool {
	return checkAtomType(t) == Atom
}

// -----------------composite-------
func isCompType(t reflect.Type) bool {
	return checkAtomType(t) == Composite
}

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

// -----------------------nuller---------------------------------
// v03
// 检查是否nuller
func isNullerType(t reflect.Type) bool {
	return t.Implements(ImpNuller)
}

// ----------------------valuer------------------------
//
//	v03
//
// 是否valuer
func isValuerType(t reflect.Type) bool {
	if _isBaseType(t) {
		return true
	}
	return t.Implements(ImpValuer)
}

//	v03
//
// 是否 valuer nuller
func isVNType(t reflect.Type) bool {
	return isValuerType(t) && isNullerType(t)
}

// -------------ptr-----------------
// v03
// 是指针类型，返回指针的基类型,如果是ptr，但是是nil，则返回error
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

// -------------ptr-deep-------------------
// v0.7
// 是指针类型，返回指针的最基类型
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

// -----------------slice---------------
// v0.7
// 是数组类型，返回数组的基类型
func baseSliceType(t reflect.Type) (bool, reflect.Type) {
	if isValuerType(t) {
		return false, t
	}

	if t.Kind() == reflect.Slice {
		return true, t.Elem()
	}
	return false, t
}

// v0.7
// is 是否 slice
func baseSliceValue(v reflect.Value) (bool, reflect.Value) {
	if isValuerType(v.Type()) {
		return false, v
	}

	if v.Kind() == reflect.Slice {
		if v.Len() == 0 {
			return false, v
		}
		return true, v.Index(0)
	}
	return false, v
}

// ---------------------slice-deep------------------
// v0.7
// 是数组类型，返回数组的最基类型
func baseSliceDeepType(t reflect.Type) (ok bool, structType reflect.Type) {
	return _baseSliceDeepType(t)
	//isSlice := false
	//tmp := t
	//
	//for true {
	//	isDeepFlag := true //base
	//
	//	is, base = _baseSliceDeepType(base)
	//	if is {
	//		isDeepFlag = false
	//		isSlice = true
	//	}
	//	if isDeepFlag {
	//		if isSlice {
	//			return true, base
	//		}
	//		return false, t
	//	}
	//	tmp = base
	//}
	//return false, t
}

// v0.6
func baseSliceDeepValue(v reflect.Value) (bool, reflect.Value, error) {
	isSlice := false
	tmp := v
	for true {
		flag := false //base
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

// 是数组类型，返回数组的最基类型
func _baseSliceDeepType(t reflect.Type) (bool, reflect.Type) {
	isSlice := false
base:
	is, t := baseSliceType(t)
	if is {
		isSlice = true
		goto base
	}
	return isSlice, t
}

// slice 用到value的一定是非空
func _baseSliceDeepValue(v reflect.Value) (bool, reflect.Value, error) {
	isSlice := false
base:
	kind := v.Kind()
	if kind == reflect.Ptr || kind == reflect.Map {
		if v.IsNil() {
			return false, v, ErrNil
		}
	}
	if kind == reflect.Slice {
		if v.IsNil() {
			return true, v, ErrNil
		}
	}

	is, v := baseSliceValue(v)
	if is {
		isSlice = true
		goto base
	}
	return isSlice, v, nil
}

// --------------------
// v03
// 包装类型：检查是 ptr 还是 slice
type packType int

const (
	None packType = iota
	Ptr
	Slice
)

// v03
// 检查是否是ptr，slice类型
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

// --------------------
// v03
// 数据原子性
// atom composite
// 检查是 atom 还是 composite
// base和实现valuer的是 atom 即可以作为实体类 字段的数据类型，可以对应一个数据库字段
// 否则，并且类型是struct、map的是composite 即拥有多个 字段属性的 实体类，对应一条查询记录
// 其他为 invalid
type atomType int

const (
	Invalid atomType = iota
	Atom
	Composite
)

// v03
func checkAtomValue(v reflect.Value) atomType {
	is := isValuerType(v.Type())
	if is {
		return Atom
	}
	is = _isStructType(v.Type())
	if is {
		return Composite
	}

	is, _ = baseMapValue(v)
	if is {
		return Composite
	}

	return Invalid
}

// v03
func checkAtomType(t reflect.Type) atomType {
	is := isValuerType(t)
	if is {
		return Atom
	}

	is = _isStructType(t)
	if is {
		return Composite
	}

	is, _ = baseMapType(t)
	if is {
		return Composite
	}
	return Invalid
}

// -----------------------map--------------------------------
// v03
// scan不需要必须nuller
// 检查map key是否string，value是否valuer
func checkMapFieldType(t reflect.Type) bool {
	if t.Key().Kind() != reflect.String {
		return false
	}

	if !isAtomType(t.Elem()) {
		return false
	}
	return true
}

// v03
// scan不需要必须nuller
// 检查map key是否string，value是否valuer
func checkMapFieldV(t reflect.Type) error {
	if t.Key().Kind() != reflect.String {
		return errors.New("map key need string")
	}
	return checkFieldVError(t.Elem())
}

// v03
// scan不需要必须nuller
// 检查map key是否string，value是否 valuer/nuller
func checkMapFieldVN(t reflect.Type) error {
	if t.Key().Kind() != reflect.String {
		return errors.New("map key need string")
	}
	return checkFieldVNError(t.Elem())
}

// v0
// 检查map key是否string，value是否valuer/nuller
func checkMapFieldValue(v reflect.Value) bool {
	key := v.MapKeys()[0]
	if key.Kind() != reflect.String {
		return false
	}

	//valuer
	t := v.MapIndex(key).Type()

	if !isAtomType(t) {
		return false
	}

	//nuller
	is := isNullerType(t)
	if !is {
		return false
	}
	return true
}

// -----------------------struct--------------------------------
// v03
// scan不需要必须nuller
// 检查 struct field，value是否valuer
func checkStructFieldV(t reflect.Type) error {
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		err := checkFieldVError(t.Field(i).Type)
		if err != nil {
			return err
		}
	}
	return nil
}

// v03
// 检查 struct field，value是否valuer/nuller
func checkStructFieldVN(t reflect.Type) error {
	numField := t.NumField()
	for i := 0; i < numField; i++ {
		err := checkFieldVNError(t.Field(i).Type)
		if err != nil {
			return err
		}
	}
	return nil
}
