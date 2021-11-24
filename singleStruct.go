package lorm

import (
	"database/sql/driver"
	"github.com/pkg/errors"
	"reflect"
)

// v0.5
// 判断是否 single or struct
func checkDestSingle(value reflect.Value) (bool,reflect.Value, error) {
	_, base := basePtrValue(value)
	is := baseStructValue(base)
	if is { //single or struct
		_, ok := base.Interface().(driver.Valuer)
		return ok, base,nil
	}

	//必定 single
	is = baseBaseValue(base)
	if is {
		return true,base, nil
	}

	is, _ = baseSliceDeepValue(base)
	if is {
		return true, base,nil
	}

	return false,base, errors.New("type err")

}
