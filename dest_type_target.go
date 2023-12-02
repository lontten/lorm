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

// todo 下面未重构--------------

//   - struct
//     struct
//
// comp-struct 获取 destBaseValue
func (ctx *ormContext) initParamDest2TableName(dest interface{}) {
	if ctx.err != nil {
		return
	}
	t := reflect.TypeOf(dest)
	if t.Kind() == reflect.String {
		ctx.tableName = dest.(string)
		return
	}
	_, base := basePtrType(t)
	ctx.destBaseType = base

}

// 检查sturct的filed是否合法，valuer，nuller
func (ctx *ormContext) checkParamDestField() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompFieldVN(ctx.destBaseValue)
}
