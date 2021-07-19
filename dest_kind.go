package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

// * struct
// * slice
// slice
func destBaseValueCheckSlice(v reflect.Value) (bool, reflect.Value, error) {
	is, base := basePtrValue(v)
	if !is {
		is, base = baseSliceValue(base)
		if !is {
			return false, base, errors.New("need ptr or slice")
		}
		e := base.Index(0)

		if e.Kind() != reflect.Struct {
			return false, base, errors.New("need a slice struct type")
		}

		return true, e, nil //true 为 slice struct
	}
	is, base = baseStructValue(base)
	if is {
		return false, base, nil //false 为 struct
	}
	return false, base, errors.New("need ptr struct type")
}
