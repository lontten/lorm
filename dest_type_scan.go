package lsql

import (
	"github.com/pkg/errors"
	"reflect"
)

// ptr slice
func (ctx *OrmContext) initScanDestList(dest interface{}) {
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

	ctyp := checkCompType(base)
	if ctyp == Invade {
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
func (ctx *OrmContext) initScanDestOne(dest interface{}) {
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
		if !isSingleType(base) {
			ctx.err = errors.New("scan can't be a slice ")
			return
		}

	}

	ctyp := Single

	if !is {
		ctyp = checkCompType(base)
		if ctyp == Invade {
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
