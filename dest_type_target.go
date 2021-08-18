package lorm

import (
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

// * struct
// * slice
// slice

// 检查数据类型，
//获取基础 value
//bool 为 是否为 slice类型
func targetDestBaseValue2Slice(dest interface{}) ([]reflect.Value, reflect.Value, error) {
	v := reflect.ValueOf(dest)
	arr := make([]reflect.Value, 0)
	is, baseValue := basePtrValue(v)
	if !is { //不是 ptr ，必须是 slice
		is, base, err := baseSliceValue(baseValue)
		if err != nil {
			return arr, v, err
		}
		if !is {
			return arr, base, errors.New("need ptr or slice")
		}
		e := base.Index(0)
		if e.Kind() != reflect.Struct {
			return arr, base, errors.New("need a slice struct type")
		}
		return utils.ToSliceValue(baseValue), e, nil //true 为 slice struct
	}
	// ptr
	is, base := baseStructValue(baseValue)
	if is {
		arr = append(arr, base)
		return arr, base, nil //arr 为空，false 为 struct
	}
	return arr, base, errors.New("need ptr struct type")
}
