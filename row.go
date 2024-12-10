package lorm

import (
	"reflect"
)

// 创建 row 返回数据，字段 对应的 struct 字段的box
// 返回值 box, vp, v
// box	struct 的 字段box列表
// vp	struct 的 引用
// v	struct 的 值
func createColBox(base reflect.Type, cfLink ColIndex2FieldNameMap) (box []any, vp, v reflect.Value) {
	vp = newStruct(base)
	v = reflect.Indirect(vp)
	length := len(cfLink)
	box = make([]any, 1)
	if length == 0 {
		box[0] = v.Addr().Interface()
		return
	}
	box = make([]any, length)
	for c, f := range cfLink {
		if f == "" { // "" 表示此列不接收
			box[c] = new([]uint8)
		} else {
			box[c] = v.FieldByName(f).Addr().Interface()
		}
	}
	return
}

// 创建 row 返回数据，字段 对应的 struct 字段的box
// 返回值 box, vp, v
// box	struct 的 字段box列表
// vp	struct 的 引用
// v	struct 的 值
func createColBoxTNew[T any](cfLink ColIndex2FieldNameMap) (box []any, vp, v reflect.Value) {
	var tP = new(T)
	vp = reflect.ValueOf(tP)
	v = reflect.Indirect(vp)
	length := len(cfLink)
	box = make([]any, 1)
	if length == 0 {
		box[0] = v.Addr().Interface()
		return
	}
	box = make([]any, length)
	for c, f := range cfLink {
		if f == "" { // "" 表示此列不接收
			box[c] = new([]uint8)
		} else {
			box[c] = v.FieldByName(f).Addr().Interface()
		}
	}
	return
}

// 创建 row 返回数据，字段 对应的 struct 字段的box
// 返回值 box, vp, v
// box	struct 的 字段 引用列表
// vp	struct 的 引用 Value
// v	struct 的 值   Value
func createColBoxT[T any](v reflect.Value, tP T, cfLink ColIndex2FieldNameMap) (box []any) {
	length := len(cfLink)
	if length == 0 {
		box = make([]any, 1)
		box[0] = tP
		return
	}
	box = make([]any, length)
	for c, f := range cfLink {
		if f == "" { // "" 表示此列不接收
			box[c] = new([]uint8)
		} else {
			box[c] = v.FieldByName(f).Addr().Interface()
		}
	}
	return
}

// sql返回 row 字段下标 对应的  struct 字段名（""表示不接收该列数据）
type ColIndex2FieldNameMap []string
