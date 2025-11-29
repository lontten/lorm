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

	"github.com/lontten/lorm/field"
	"github.com/pkg/errors"
)

// 校验struct 的 field 是否合法
// 1. check valuer，不是 valuer 则返回error
func checkFieldV(t reflect.Type) error {
	_, base := basePtrType(t)
	is := isValuerType(base)
	if !is {
		return errors.New("field没有实现valuer " + base.String())
	}
	return nil
}

// 校验struct 的 field 是否合法
// 1. check valuer，不是 valuer 则返回error
func isFieldV(t reflect.Type) bool {
	_, base := basePtrType(t)
	return isValuerType(base)
}

// 校验struct 的 field 是否合法 ：没有同时 valuer scanner 则报错
func checkFieldVS(t reflect.Type) error {
	_, base := basePtrType(t)
	is := isValuerType(base)
	if !is {
		return errors.New("field  no imp valuer:: " + t.String())
	}
	is = isScannerType(base)
	if !is {
		return errors.New("field  no imp scanner:: " + t.String())
	}
	return nil
}

// 校验struct 的 field 是否合法 ：没有同时 valuer scanner
func isFieldVS(t reflect.Type) bool {
	_, base := basePtrType(t)
	return isValuerType(base) && isScannerType(base)
}

// 零值为 nil
func getFieldInterZero(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return v.Elem().Interface()
	}
	if v.IsZero() {
		return nil
	}
	return v.Interface()
}

// 返回值类型有 None,Null,Val,三种
func getFieldInter(v reflect.Value) field.Value {
	if !v.IsValid() {
		return field.Value{
			Type: field.None,
		}
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return field.Value{
				Type: field.Null,
			}
		}
		return field.Value{
			Type:  field.Val,
			Value: v.Elem().Interface(),
		}
	}
	return field.Value{
		Type:  field.Val,
		Value: v.Interface(),
	}
}

// ptr nil 为 null
// 非ptr，零值为 null
func isFieldNull(v reflect.Value) bool {
	if !v.IsValid() {
		return true
	}
	if v.Kind() != reflect.Ptr {
		return v.IsZero()
	}
	return v.IsNil()
}
