package lorm

import (
	"errors"
	"reflect"
)

/**
从传入的struct实例，获取实例类对应的表名，解析字段是否合法
*/
// ptr
// 检查数据类型  struct
func (ctx *ormContext) initModelDest(dest any) {
	if ctx.hasErr() {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, v, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	ctx.scanIsPtr = isPtr
	t := v.Type()

	if ctx.checkParam {
		if _isStructType(t) {
			ctx.err = errors.New("dest need is struct")
			return
		}
		err = checkCompFieldVS(v)
		if err != nil {
			ctx.err = err
			return
		}
	}

	ctx.paramModelBaseV = v

	ctx.scanDest = dest
	ctx.scanIsPtr = isPtr

	ctx.destV = v
	ctx.destBaseType = t
	ctx.destBaseTypeIsComp = true

	ctx.destIsSlice = false
	ctx.destSliceItemIsPtr = false
}

// string 或者 struct
func (ctx *ormContext) initNameDest(dest any) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	_, base, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if base.Kind() == reflect.String {
		ctx.tableName = dest.(string)
		return
	}
	ctx.initModelDest(dest)
}

// 参数分为comp 的 vn
// 接受是comp、atom 的 v
// 即：检查 comp 的 vn
func (ctx *ormContext) checkDestParamScan() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompFieldVS(ctx.paramModelBaseV)
}

// 参数分为comp vn
// 即：检查 comp 的 vn
func (ctx *ormContext) checkDestParam() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompFieldVS(ctx.paramModelBaseV)
}

// 接受是comp、atom 的 v
func (ctx *ormContext) checkDestScan() {
	if ctx.hasErr() {
		return
	}
	if isValuerType(ctx.destBaseType) {
		return
	}
	ctx.err = checkCompFieldV(ctx.destBaseType)
}
