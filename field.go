package lorm

import (
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"reflect"
)

// v03
// 校验struct 的 field 是否合法
// 1. check valuer，不是 valuer 则返回error
func checkFieldVError(t reflect.Type) error {
	_, base := checkPackType(t)

	is := isValuerType(base)
	if !is {
		return errors.New("field没有实现valuer " + t.String())
	}
	return nil
}

// v03
// 校验struct 的 field 是否合法
// 1. check valuer，不是 valuer 则返回error
func checkFieldV(t reflect.Type) bool {
	_, base := checkPackType(t)
	return isValuerType(base)
}

// v03
// 校验struct 的 field 是否合法 ：没有同时 valuer nuller 则报错
// 1. check single
// 3. nuller
func checkFieldVNError(t reflect.Type) error {
	isNuller := false
	typ, base := checkPackType(t)
	if typ != None {
		//如果是 ptr、slice类型，肯定是有 nuller
		isNuller = true
	} else {
		//直接判断是否 nuller
		isNuller = isNullerType(base)
	}

	is := isValuerType(base)
	if !is {
		return errors.New("field  no imp valuer:: " + t.String())
	}
	//nuller
	if isNuller {
		return nil
	}
	return errors.New("field  no imp nuller:: " + t.String())
}

// v03
// 校验struct 的 field 是否合法 ：没有同时 valuer nuller 则报错
// 1. check single
// 3. nuller
func checkFieldVN(t reflect.Type) bool {
	isNuller := false
	typ, base := checkPackType(t)
	if typ != None {
		//如果是 ptr、slice类型，肯定是有 nuller
		isNuller = true
	} else {
		//直接判断是否 nuller
		isNuller = isNullerType(base)
	}

	isValuer := isValuerType(base)
	return isNuller && isValuer
}

// 获取field的值 interface类型
// 1. 先去ptr，2.再去slice，最后取值
// 这里默认 v是 atom类型，有 valuer
func getFieldInter(v reflect.Value) interface{} {
	_, v, err := basePtrDeepValue(v)
	if err != nil {
		return nil
	}

	is, _, err := baseSliceDeepValue(v)
	if err != nil {
		return nil
	}

	if is {
		value, err := types.ArrayOf(v.Interface()).Value()
		if err != nil {
			return nil
		}
		return value
	}
	return v.Interface()
}

// 获取field的值 interface类型
// 1. 先去ptr，2.再去slice，最后取值
// 这里没有默认 v是 atom 类型，参数不一定是 有 valuer ，所以有isValuerType 判断
func getTargetInter(v reflect.Value) interface{} {
	_, v, err := basePtrDeepValue(v)
	if err != nil {
		return nil
	}
	if isValuerType(v.Type()) {
		return v.Interface()
	}

	is, _, err := baseSliceDeepValue(v)
	if err != nil {
		return nil
	}

	if is {
		value, err := types.ArrayOf(v.Interface()).Value()
		if err != nil {
			return nil
		}
		return value
	}
	return v.Interface()
}

// todo 下面未重构--------------
