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
	"time"
)

func FieldSetValNil(f reflect.Value, val any) error {
	if val == nil {
		f.SetZero()
		return nil
	}

	switch nullValPtr := val.(type) {
	case *any:
		if nullValPtr == nil {
			f.SetZero()
			return nil
		} else {
			if *nullValPtr == nil {
				f.SetZero()
				return nil
			}
			// 如果 字段实现了 sql.Scanner，但是因为 不是 指针，就创建一个临时指针类型接收，然后再赋值
			newVal := reflect.New(f.Type())
			scanner := newVal.Interface().(sql.Scanner)
			if err := scanner.Scan(*nullValPtr); err != nil {
				return err
			}
			f.Set(newVal.Elem())
			return nil
		}
	case *sql.NullString:
		if nullValPtr.Valid {
			f.SetString(nullValPtr.String)
		} else {
			f.SetString("")
		}
	case *sql.NullInt64:
		if nullValPtr.Valid {
			f.SetInt(nullValPtr.Int64)
		} else {
			f.SetInt(0)
		}
	case *sql.NullFloat64:
		if nullValPtr.Valid {
			f.SetFloat(nullValPtr.Float64)
		} else {
			f.SetFloat(0)
		}
	case *sql.NullBool:
		if nullValPtr.Valid {
			f.SetBool(nullValPtr.Bool)
		} else {
			f.SetBool(false)
		}
	case *sql.NullTime:
		if nullValPtr.Valid {
			f.Set(reflect.ValueOf(nullValPtr.Time))
		} else {
			f.Set(reflect.ValueOf(time.Time{}))
		}
	case *sql.RawBytes:
		f.SetBytes(*nullValPtr)
	}
	return nil
}
