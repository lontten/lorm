package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

func (ctx *ormContext) initScanDestList(dest any) {
	if ctx.hasErr() {
		return
	}

	v := reflect.ValueOf(dest)
	isPtr, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if !isPtr {
		ctx.err = errors.New("scan must be a ptr ")
		return
	}

	isSlice, t := baseSliceType(v.Type())
	if !isSlice {
		ctx.err = errors.New("scanList must be a slice ")
		return
	}
	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanIsPtr = isPtr

	ctx.destV = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = true
	ctx.destSliceItemIsPtr = t.Kind() == reflect.Ptr
}

func (ctx *ormContext) initScanDestListT(dest any, v reflect.Value, t reflect.Type, destSliceItemIsPtr bool) {
	if ctx.hasErr() {
		return
	}

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanIsPtr = true

	ctx.destV = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = true
	ctx.destSliceItemIsPtr = destSliceItemIsPtr
}

func (ctx *ormContext) initScanDestOne(dest any) {
	if ctx.hasErr() {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if !isPtr {
		ctx.err = errors.New("scan must be a ptr ")
		return
	}

	t := v.Type()

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanIsPtr = isPtr
	ctx.destV = v

	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}

// dest 类型 struct 、所有 valuer 类型
func (ctx *ormContext) initScanDestOneT(dest any) {
	if ctx.hasErr() {
		return
	}

	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanIsPtr = true

	ctx.destV = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}

// 从dest中获取filed 的名字，dest必须为struct或者*struct
func (ctx *ormContext) initDestScanField(dest any) {
	if ctx.hasErr() {
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

	err = checkFieldV(v.Type())
	if err != nil {
		ctx.err = err
		return
	}
	//ctx.scanDest = dest
	//
	//ctx.scanIsSlice = false
	//ctx.destSliceItemIsPtr = false
	//
	//ctx.scanDestBaseType = base
	//ctx.destBaseTypeIsComp = ctyp == Composite
	//
	//ctx.destValue = v

	//todo 把filed 获取到存入 ctx

}
