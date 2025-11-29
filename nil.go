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
	"database/sql"
	"reflect"
)

// checkHandleNull 测试类型的 Scan 方法是否能处理 NULL（即 value 为 nil 时不报错）
func checkHandleNull(typ reflect.Type) (canNull bool, isScanner bool) {
	if !typ.Implements(ImpScanner) {
		typ = reflect.PointerTo(typ)
		if !typ.Implements(ImpScanner) {
			return
		}
	}

	isScanner = true

	// 创建该类型的实例，尝试用 nil 调用其 Scan 方法，看是否报错
	instance := reflect.New(typ.Elem()).Interface().(sql.Scanner)
	err := instance.Scan(nil)
	canNull = err == nil // 不报错说明能处理 NULL
	return
}
