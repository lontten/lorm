//  Copyright 2025 lontten lontten@163.com
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package lorm

import (
	"reflect"
)

// 创建 row 返回数据，字段 对应的 struct 字段的box
// 返回值 box, vp, v
// box	struct 的 字段box列表
// vp	struct 的 引用
// v	struct 的 值
func (ctx *ormContext) createColBoxNew(base reflect.Type, cfLink map[string]compC, rowColumnTypeMap map[int]rowColumnType) (box []any, vp, v reflect.Value, fun func() error) {
	vp = reflect.New(base)
	v = reflect.Indirect(vp)
	tP := vp.Interface()

	colBox, fun := ctx.createColBox(v, tP, cfLink, rowColumnTypeMap)
	return colBox, vp, v, fun
}

// 创建 row 返回数据，字段 对应的 struct 字段的box
// 返回值 box, vp, v
// box	struct 的 字段 引用列表
// vp	struct 的 引用 Value
// v	struct 的 值   Value
func (ctx *ormContext) createColBox(v reflect.Value, tP any, cfLink map[string]compC, rowColumnTypeMap map[int]rowColumnType) (box []any, fun func() error) {
	fun = func() error { return nil }
	length := len(cfLink)
	if length == 0 {
		box = make([]any, 1)
		box[0] = tP
		return
	}

	box = make([]any, length)
	var converters []func() error
	fun = func() error {
		for _, f := range converters {
			e := f()
			if e != nil {
				return e
			}
		}
		return nil
	}

	for _, f := range cfLink {
		if f.columnName == "" { // "" 表示此列不接收
			box[f.columnIndex] = new(any)
			continue
		}

		field := v.FieldByName(f.fieldName)

		valBox, convertFunc := ctx.convertCtx.Get(f.columnName)
		if valBox != nil && convertFunc != nil {
			box[f.columnIndex] = valBox

			converters = append(converters, func(f reflect.Value, fieldName string, val any, c compC) func() error {
				return func() error {
					val = convertFunc(val)
					f.Set(reflect.ValueOf(val))
					return nil
				}
			}(field, f.fieldName, valBox, f))
			continue
		}

		columnType := rowColumnTypeMap[f.columnIndex]
		// 字段可以接收null 或者 返回值不为null;可以直接用 结构体字段接收
		if f.canNull || columnType.noNull {
			box[f.columnIndex] = field.Addr().Interface()
			continue
		}

		var tmpVal any
		if f.isScanner {
			tmpVal = new(any)
		} else {
			tmpVal = allocDatabaseType(columnType.databaseTypeName)
			if tmpVal == nil {
				tmpVal = allocType(field.Type())
				if tmpVal == nil {
					panic("field not support")
				}
			}
		}
		box[f.columnIndex] = tmpVal

		converters = append(converters, func(f reflect.Value, fieldName string, val any, c compC) func() error {
			return func() error {
				return FieldSetValNil(f, fieldName, val)
			}
		}(field, f.fieldName, tmpVal, f))
	}
	return
}

// sql返回 row 字段下标 对应的  struct 字段名（""表示不接收该列数据）
type ColIndex2FieldNameMap []string
