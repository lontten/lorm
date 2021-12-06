package lorm

import (
	"github.com/pkg/errors"
	"reflect"
)

// ptr
// 检查数据类型 comp-struct
func (ctx *OrmContext) initTargetDest(dest interface{}) {
	if ctx.err != nil {
		return
	}
	v := reflect.ValueOf(dest)
	isPtr, base, err := basePtrDeepValue(v)
	if err != nil {
		ctx.err = err
		return
	}
	if !isPtr {
		ctx.err = errors.New("need a ptr")
		return
	}

	err = checkCompField(base)
	if err != nil {
		ctx.err = errors.New("need a struct")
		return
	}

	ctx.scanDest = dest
	ctx.scanDest = dest

	ctx.destValue = base
	ctx.destBaseValue = base

	ctx.destBaseType = base.Type()
	ctx.scanDestBaseType = base.Type()
}

// * struct
//  struct
// comp-struct 获取 destBaseValue
func (ctx *OrmContext) initTargetDest2TableName(dest interface{}) {
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

//检查sturct的filed是否合法，valuer，nuller
func (ctx *OrmContext) checkTargetDestField() {
	if ctx.err != nil {
		return
	}
	ctx.err = checkCompField(ctx.destBaseValue)
}
