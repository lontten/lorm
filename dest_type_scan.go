package lorm

import (
	"fmt"
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
	fmt.Println(base.String())

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

//
////struct 只检查 struct是否合格，不检查 filed
//func checkScanTypeLn(t reflect.Type) (reflect.Type, error) {
//	is, base := basePtrType(t)
//	if !is {
//		return t, errors.New("need a ptr")
//	}
//	code, base := baseStructBaseType(t)
//	if code<0 {
//		return t, errors.New("need a ptr struct or base type")
//	}
//	return base, nil
//}
//// slice 只检查 struct是否合格，不检查 filed
//func checkScanType(t reflect.Type) (reflect.Type, error) {
//	_, base := basePtrType(t)
//	is, base := baseSliceType(base)
//	if !is {
//		return t, errors.New("need a slice type")
//	}
//
//	baseType, _ := baseStructBaseSliceType(base)
//
//	if baseType < 0 {
//		return t, errors.New("need a slice struct or base type")
//	}
//	return base, nil
//}
//
//
