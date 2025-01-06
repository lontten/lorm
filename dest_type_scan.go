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
	_, base := basePtrType(t)
	ctyp := checkAtomType(base)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = isPtr

	ctx.destBaseType = base
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = true
	ctx.destSliceItemIsPtr = t.Kind() == reflect.Ptr
}

func (ctx *ormContext) initScanDestListT(dest any, v, baseV reflect.Value, t reflect.Type, destSliceItemIsPtr bool) {
	if ctx.hasErr() {
		return
	}

	ctyp := checkAtomType(t)
	if ctyp == Invalid {
		ctx.err = errors.New("scan type is not supported")
		return
	}

	ctx.scanDest = dest
	ctx.scanV = v
	ctx.scanIsPtr = true

	ctx.destBaseValue = baseV
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
	ctx.scanV = v
	ctx.scanIsPtr = isPtr

	ctx.destBaseValue = v
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
	ctx.scanV = v
	ctx.scanIsPtr = true

	ctx.destBaseValue = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = ctyp == Composite

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}
