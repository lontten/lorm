package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

/**
从传入的struct实力，获取实例类对应的表名，解析字段是否合法
*/
// ptr
// 检查数据类型  struct
func (ctx *ormContext) initTargetDest(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, base, err := basePtrValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	ctx.destIsPtr = isPtr
	if _isStructType(base.Type()) {
		ctx.err = errors.New("dest need is struct")
		return
	}

	err = checkCompFieldVN(base)
	if err != nil {
		ctx.err = err
		return
	}

	ctx.scanDest = dest
	ctx.scanDest = dest

	ctx.destValue = base
	ctx.destBaseValue = base

	ctx.destBaseType = base.Type()
	ctx.scanDestBaseType = base.Type()
}

// string 或者 struct
func (ctx *ormContext) initNameDest(dest interface{}) {
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
	ctx.initTargetDest(dest)
}

// todo 下面未重构--------------

// 检查sturct的filed是否合法，valuer，nuller
func (ctx *ormContext) checkParamDestField() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompFieldVN(ctx.destBaseValue)
}

// 参数分为comp 的 vn
// 接受是comp、atom 的 v
// 即：检查 comp 的 vn
func (ctx *ormContext) checkDestParamScan() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompFieldVN(ctx.destBaseValue)
}

// 参数分为comp vn
// 即：检查 comp 的 vn
func (ctx *ormContext) checkDestParam() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompFieldVN(ctx.destBaseValue)
}

// 接受是comp、atom 的 v
func (ctx *ormContext) checkDestScan() {
	if ctx.err != nil {
		return
	}
	if isValuerType(ctx.destBaseType) {
		return
	}
	ctx.err = checkCompFieldV(ctx.destBaseType)
}
