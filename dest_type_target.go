package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

// * struct
// * slice
// slice

// 检查数据类型，
//获取基础 value
//bool 为 是否为 slice类型
func targetDestBaseValueCheckSlice(v reflect.Value) (bool, reflect.Value, error) {
	is, base := basePtrValue(v)
	if !is {
		is, base,err := baseSliceValue(base)
		if err != nil {
			return false, v, err
		}
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

