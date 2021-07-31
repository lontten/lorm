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
func targetDestBaseValue2Slice(dest interface{}) ([]interface{}, reflect.Value, error) {
	v := reflect.ValueOf(dest)
	arr := make([]interface{}, 0)
	is, baseValue := basePtrValue(v)
	if !is { //不是 ptr ，必须是 slice
		is, base,err := baseSliceValue(baseValue)
		if err != nil {
			return arr, v, err
		}
		if !is {
			return arr, base, errors.New("need ptr or slice")
		}
		if baseValue.Len() == 0 {
			return arr, base, errors.New("slice can't len is 0")
		}
		e := base.Index(0)
		if e.Kind() != reflect.Struct {
			return arr, base, errors.New("need a slice struct type")
		}
		return utils.ToSlice(baseValue), e, nil //true 为 slice struct
	}
	// ptr
	is, base := baseStructValue(baseValue)
	if is {
		arr = append(arr, dest)
		return arr, base, nil //arr 为空，false 为 struct
	}
	return arr, base, errors.New("need ptr struct type")
}

