package lorm

import (
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
	"reflect"
)

// ptr slice
// 检查数据类型 comp-struct
func (ctx *OrmContext) initScanDestSlice(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)

	arr := make([]reflect.Value, 0)

	packTyp, err := checkPackValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	typ:=packTyp.Typ
	v=packTyp.Base
	base:=packTyp.SliceBase

	ctyp := checkCompValue(base)

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
		ctx.scanValueArr = arr
		return
	}

	if typ == Slice {
		ctx.isSlice = true
		ctx.scanValueArr = utils.ToSliceValue(v)
		return
	}

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
