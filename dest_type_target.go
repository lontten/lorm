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
func (c *OrmContext) initTargetDest(dest interface{}) {
	v := reflect.ValueOf(dest)
	arr := make([]reflect.Value, 0)
	isPtr, baseValue := basePtrValue(v)
	c.dest=dest
	c.destValue=baseValue

	is, base := baseStructValue(baseValue)
	if isPtr && is {
		arr = append(arr, base)
		c.isSlice = false
		c.destValueArr = arr
		c.destBaseValue = base
		return
	}
	//slice
	is, has, base := baseSliceValue(baseValue)
	if !has {
		c.err = ErrContainEmpty
		return
	}
	if !is {
		c.err = errors.New("need ptr or slice")
		return
	}
	if base.Kind() != reflect.Struct {
		c.err = errors.New("need a slice struct type")
		return
	}
	c.isSlice = true
	c.destValueArr = utils.ToSliceValue(baseValue)
	c.destBaseValue = base
}

func (c *OrmContext) checkTargetDestField() {
	v := c.destBaseValue
	err := checkValidFieldTypStruct(v)
	if err != nil {
		c.err = err
	}
	return
}
