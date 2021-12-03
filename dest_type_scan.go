package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

// ptr slice
// 检查数据类型 valuer
func (ctx *OrmContext) initScanDestSlice(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	_, v, err := basePtrDeepValue(v)

	if err != nil {
		ctx.err = err
		return
	}

	typ, base := checkPackType(v.Type())

	ctyp := checkCompType(base)
	if ctyp == Invade {
		ctx.err = errors.New("need a struct or base type -scan dest slice")
		return
	}

	if typ == Slice {
		ctx.isSlice = true
		ctx.sliceItemIsPtr = base.Kind() == reflect.Ptr
	}

	ctx.destTypeIsComp = ctyp == Composite
	ctx.dest = dest
	ctx.destValue = v
	ctx.destBaseType = base

}
