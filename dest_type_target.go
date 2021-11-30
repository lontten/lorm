package lorm

import (
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

// ptr slice
// 检查数据类型 comp-struct
func (ctx *OrmContext) initTargetDestSlice(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)

	arr := make([]reflect.Value, 0)

	typ, base,err := checkPackValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	ctyp := checkCompValue(base, false)

	if ctyp == SliceEmpty {
		ctx.err = errors.New("slice no's empty")
		return
	}
	if ctyp != Composite {
		ctx.err = errors.New("need a struct")
		return
	}

	ctx.dest = dest
	ctx.destValue = v
	ctx.destBaseValue = base

	if typ == Ptr {
		arr = append(arr, base)
		ctx.isSlice = false
		ctx.destValueArr = arr
		return
	}

	if typ == Slice {
		ctx.isSlice = true
		ctx.destValueArr = utils.ToSliceValue(v)
		return
	}

}

// ptr
//不能是slice
// 检查数据类型 comp-struct
func (ctx *OrmContext) initTargetDest(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, base, err := basePtrDeepValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if !isPtr {
		ctx.err = errors.New("need a ptr")
		return
	}

	err = checkCompField(base)
	if err != nil {
		ctx.err = errors.New("need a struct")
		return
	}

	arr := make([]reflect.Value, 0)

	ctx.dest = dest
	ctx.destValue = base
	ctx.destBaseValue = base

	arr = append(arr, base)
	ctx.isSlice = false
	ctx.destValueArr = arr
	return

}

// * struct
//  struct
// comp-struct 获取 destBaseValue
func (ctx *OrmContext) initTargetDestOnlyBaseValue(dest interface{}) {
	if ctx.err != nil {
		return
	}
	value := reflect.ValueOf(dest)
	_, base, err := basePtrDeepValue(value)
	if err != nil {
		ctx.err = err
		return
	}
	err = checkCompField(base)
	if err != nil {
		ctx.err = errors.New("need a struct")
		return
	}
	ctx.destBaseValue = base
}

//检查sturct的filed是否合法，valuer，nuller
func (ctx *OrmContext) checkTargetDestField() {
	if ctx.err != nil {
		return
	}
	v := ctx.destBaseValue
	err := checkCompField(v)
	ctx.err = err
}

//检查sturct的filed是否合法，valuer，nuller
func (ctx *OrmContext) checkScanDestField() {
	if ctx.err != nil {
		return
	}
	v := ctx.destBaseValue
	err := checkCompFieldScan(v)
	ctx.err = err
}
