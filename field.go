package lorm

import (
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"reflect"
)

// v0.6
//校验struct 的 field 是否合法
//1. check single
//2. valuer
// v0.6
func checkField(v reflect.Value) error {
	_, base := checkPackTypeValue(v)
	typ := checkCompTypeValue(base, true)
	if typ != Single {
		return errors.New("check field err :: " + v.String())
	}
	return nil
}

// v0.6
//校验struct 的 field 是否合法
//1. check single
//2. valuer
//3. nuller
func checkFieldNuller(v reflect.Value) error {
	isNuller := false
	typ, base := checkPackTypeValue(v)
	if typ != None {
		isNuller = true
	} else {
		isNuller = checkBaseNuller(base)
	}

	ctyp := checkCompTypeValue(base, false)
	if ctyp != Single {
		return errors.New("check field err :: " + v.String())
	}
	//nuller
	if isNuller {
		return nil
	}
	return errors.New("check field err :: " + v.String())
}

// v0.6
//获取field的值 interface类型
func getFieldInter(v reflect.Value) interface{} {
	typ, base := checkPackTypeValue(v)
	switch typ {
	case None:
		return nil
	case Ptr:
		return base.Interface()
	case Slice:
		if base.Len() == 0 {
			return nil
		}
		value, err := types.ArrayOf(v.Interface()).Value()
		if err != nil {
			return nil
		}
		return value
	}
	return nil
}
