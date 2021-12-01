package lorm

import (
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"reflect"
)

// v0.7
//校验struct 的 field 是否合法
//1. check single
func checkField(t reflect.Type) error {
	_, base := checkPackType(t)

	typ := checkCompType(base)
	if typ != Single {
		return errors.New("field没有实现valuer " + t.String())
	}
	return nil
}

// v0.7
//校验struct 的 field 是否合法
//1. check single
//3. nuller
func checkFieldNuller(t reflect.Type) error {
	isNuller := false
	typ, base := checkPackType(t)
	if typ != None {
		isNuller = true
	} else {
		isNuller = isNullerType(base)
	}

	ctyp := checkCompType(base)
	if ctyp != Single {
		return errors.New("field  no imp valuer:: " + t.String())
	}
	//nuller
	if isNuller {
		return nil
	}
	return errors.New("field  no imp nuller:: " + t.String())
}

// v0.7
//获取field的值 interface类型
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
