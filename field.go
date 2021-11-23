package lorm

import (
	"database/sql/driver"
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"reflect"
)
// v0.6
//只校验struct 的 field 是否合法
//1. check single
//2. valuer
// v0.6
func checkField(v reflect.Value) error {
	_, base := checkPackTypeValue(v)
	typ, base := checkCompTypeValue(base, true)
	if typ!=Single {
		return errors.New("check field err :: "+	v.String())
	}
	return nil
}
// v0.6
//只校验struct 的 field 是否合法
//1. check single
//2. valuer
//3. nuller
func checkFieldNuller(v reflect.Value) error {
	isNuller:=false
	typ, base := checkPackTypeValue(v)
	if typ!=None {
		isNuller=true
	}else {
		isNuller = checkBaseNuller(base)
	}

	ctyp, base := checkCompTypeValue(base, false)
	if ctyp!=Single {
		return errors.New("check field err :: "+	v.String())
	}
	//nuller
	if isNuller {
		return nil
	}
	return errors.New("check field err :: "+	v.String())
}


// v0.5
func getFieldInter(v reflect.Value) interface{} {

	isPtr, base := basePtrValue(v)
	if !isPtr {
		//必须 nuller struct
		is, base := baseStructValue(base)
		if is {
			vv := base.Interface().(types.NullEr)
			if !vv.IsNull() {
				value, _ := base.Interface().(driver.Valuer).Value()
				return value
			}
		}

		_, has, _ := baseSliceDeepValue(base)
		if has {
			return base.Interface()
		}

	}

	code, base := baseStructBaseValue(base)
	if code > 0 {
		return base.Interface()
	}

	_, has, _ := baseSliceDeepValue(base)
	if has {
		return base.Interface()
	}

	return nil
}
