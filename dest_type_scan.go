package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

/**
从传入的struct实体，获取需要scan的目标类型，
分为scanList，scanOne，
后面会改为泛型，则不再需要这些。
*/
// todo 下面未重构--------------
// ptr slice
func (ctx *ormContext) initScanDestList(dest interface{}) {
	if ctx.err != nil {
		return
	}
	if dest == nil {
		ctx.err = errors.New("scanList is nil")
		return
	}

	v := reflect.ValueOf(dest)
	is, v, err := basePtrValue(v)
	if !is {
		ctx.err = errors.New("scanList must be a ptr ")
		return
	}
	if err != nil {
		ctx.err = err
		return
	}
	is, base := baseSliceType(v.Type())
	if !is {
		ctx.err = errors.New("scanList must be a slice ")
		return
	}

	ctyp := checkAtomType(base)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest

	ctx.scanIsSlice = true
	ctx.scanSliceItemIsPtr = base.Kind() == reflect.Ptr

	ctx.scanDestBaseType = base
	ctx.scanDestBaseTypeIsComp = ctyp == Composite

	ctx.destValue = v

}

// ptr slice
func (ctx *ormContext) initScanDestOne(dest interface{}) {
	if ctx.err != nil {
		return
	}
	if dest == nil {
		ctx.err = errors.New("scan is nil")
		return
	}

	v := reflect.ValueOf(dest)
	is, v, err := basePtrValue(v)
	if !is {
		ctx.err = errors.New("scan must be a ptr ")
		return
	}
	if err != nil {
		ctx.err = err
		return
	}

	is, base := baseSliceType(v.Type())
	if is {
		if !isAtomType(base) {
			ctx.err = errors.New("scan can't be a slice ")
			return
		}

	}

	ctyp := Atom

	if !is {
		ctyp = checkAtomType(base)
		if ctyp == Invalid {
			ctx.err = errors.New("scan type is not supported")
			return
		}
	}

	ctx.scanDest = dest

	ctx.scanIsSlice = false
	ctx.scanSliceItemIsPtr = false

	ctx.scanDestBaseType = base
	ctx.scanDestBaseTypeIsComp = ctyp == Composite

	ctx.destValue = v

}

// v03
// 从dest中获取filed 的名字，dest必须为struct或者*struct
func (ctx *ormContext) initDestScanField(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	_, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	is := isCompType(v.Type())
	if is {
		ctx.err = errors.New("dest need is struct or map")
		return
	}
	is, _ = baseSliceType(v.Type())
	if is {
		ctx.err = errors.New("dest cannot slice")
		return
	}

	err = checkFieldVError(v.Type())
	if err != nil {
		ctx.err = err
		return
	}
	//ctx.scanDest = dest
	//
	//ctx.scanIsSlice = false
	//ctx.scanSliceItemIsPtr = false
	//
	//ctx.scanDestBaseType = base
	//ctx.scanDestBaseTypeIsComp = ctyp == Composite
	//
	//ctx.destValue = v

	//todo 把filed 获取到存入 ctx

}
