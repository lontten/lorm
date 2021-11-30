package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

// v0.5
// 判断是否 single or struct
func checkDestSingle(value reflect.Value) (bool,reflect.Value, error) {
	_, base,err := basePtrValue(value)
	if err != nil {
		return false, reflect.Value{}, err
	}
	is := isValuerType(base.Type())
	if is { //single or struct
		return true, base,nil
	}

	//必定 single
	is = _isBaseType(base.Type())
	if is {
		return true,base, nil
	}

	is, _,err = baseSliceDeepValue(base)
	if err != nil {
		return false, reflect.Value{}, err
	}
	if is {
		return true, base,nil
	}

	return false,base, errors.New("type err")

}
