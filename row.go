package lsql

import (
	"reflect"
)

//创建用来存放row中值得 引用
func createColBox(base reflect.Type, cfLink ColFieldIndexLinkMap) (box []interface{}, vp, v reflect.Value) {
	vp = newStruct(base)
	v = reflect.Indirect(vp)
	length := len(cfLink)
	box = make([]interface{}, 1)
	if length == 0 {
		box[0] = v.Addr().Interface()
		return
	}
	box = make([]interface{}, length)
	for c, f := range cfLink {
		if f < 0 { // -1 表示此列不接收
			box[c] = new([]uint8)
		} else {
			box[c] = v.Field(f).Addr().Interface()
		}

	}
	return
}

type ColFieldIndexLinkMap []int
