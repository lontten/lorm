package lorm

import (
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

// ptr slice
// 检查数据类型 comp-struct
func (c *OrmContext) initTargetDestSlice(dest interface{}) {
	if c.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	_, v = basePtrDeepValue(v)

	arr := make([]reflect.Value, 0)

	typ, base := checkPackTypeValue(v)
	ctyp := checkCompTypeValue(base, false)

	if ctyp == SliceEmpty {
		c.err = errors.New("slice no's empty")
		return
	}
	if ctyp != Composite {
		c.err = errors.New("need a struct")
		return
	}

	c.dest = dest
	c.destValue = v
	c.destBaseValue = base

	if typ == Ptr {
		arr = append(arr, base)
		c.isSlice = false
		c.destValueArr = arr
		return
	}

	if typ == Slice {
		c.isSlice = true
		c.destValueArr = utils.ToSliceValue(v)
		return
	}

}

// ptr
//不能是slice
// 检查数据类型 comp-struct
func (c *OrmContext) initTargetDest(dest interface{}) {
	if c.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	_, v = basePtrDeepValue(v)

	arr := make([]reflect.Value, 0)

	is, base := basePtrDeepValue(v)
	if !is {
		c.err = errors.New("need a ptr")
		return
	}
	is = checkCompStructValue(base)
	if !is {
		c.err = errors.New("need a struct")
		return
	}

	c.dest = dest
	c.destValue = v
	c.destBaseValue = base

	arr = append(arr, base)
	c.isSlice = false
	c.destValueArr = arr
	return

}

// * struct
//  struct
// comp-struct 获取 destBaseValue
func (c *OrmContext) initTargetDestOnlyBaseValue(dest interface{}) {
	if c.err != nil {
		return
	}
	value := reflect.ValueOf(dest)
	_, base := basePtrDeepValue(value)
	is := checkCompStructValue(base)
	if !is {
		c.err = errors.New("need a struct")
		return
	}
	c.destBaseValue = base
}

//检查sturct的filed是否合法，valuer，nuller
func (c *OrmContext) checkTargetDestField() {
	if c.err != nil {
		return
	}
	v := c.destBaseValue
	err := checkStructValidFieldNuller(v)
	c.err = err
}
