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
	"strings"
	"time"
)

// allocKindType 根据reflect.Kind分配对应的SQL数据类型
// 返回值为any类型，实际为sql包中的 nullable 类型指针
func allocKindType(kind reflect.Kind) any {
	switch kind {
	// 字符串类型
	case reflect.String:
		return &sql.NullString{}

	// 整数类型
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return &sql.NullInt64{}

	// 浮点数类型
	case reflect.Float32, reflect.Float64:
		return &sql.NullFloat64{}

	// 布尔类型
	case reflect.Bool:
		return &sql.NullBool{}

	// 结构体类型（特殊处理时间类型）
	case reflect.Struct:
		// 注意：通过Kind无法直接判断是否为time.Time，因为它本质也是struct
		// 这里需要额外判断具体类型
		// 实际使用时可能需要传入reflect.Type而非Kind才能准确判断
		return nil

	// 切片类型（二进制数据）
	case reflect.Slice:
		// 字节切片对应二进制类型
		return &sql.RawBytes{}

	// 指针类型
	case reflect.Ptr:
		// 指针类型需要根据指向的元素类型处理
		// 这里简单返回nil，实际使用时可能需要解引用
		return nil

	// 其他类型
	default:
		return nil
	}
}

// 辅助函数：根据具体类型分配SQL类型（处理time.Time等特殊情况）
func allocType(t reflect.Type) any {
	// 特殊处理时间类型
	if t == reflect.TypeOf(time.Time{}) {
		return &sql.NullTime{}
	}
	// 处理字节切片
	if t == reflect.TypeOf([]byte{}) {
		return &sql.RawBytes{}
	}
	// 其他类型通过Kind处理
	return allocKindType(t.Kind())
}

func allocDatabaseType(databaseType string) any {
	// 统一转换为大写并去除前后空格
	dbType := strings.ToUpper(strings.TrimSpace(databaseType))

	switch dbType {
	// 字符串类型
	case "VARCHAR", "TEXT", "CHAR", "CLOB", "NCLOB", "NVARCHAR", "NVARCHAR2",
		"LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "NCHAR", "CHARACTER VARYING",
		"CHARACTER", "VARCHAR2", "STRING":
		return &sql.NullString{}

	// 整数类型
	case "INT", "INTEGER", "BIGINT", "SMALLINT", "TINYINT", "MEDIUMINT",
		"SERIAL", "BIGSERIAL", "SMALLSERIAL":
		return &sql.NullInt64{}

	// 精确数值类型（金额等）- 使用字符串避免精度损失
	case "NUMERIC", "DECIMAL", "NUMBER", "MONEY":
		return &sql.NullString{}

	// 浮点数类型（适合科学计算，不适合金额）
	case "FLOAT", "DOUBLE", "REAL", "BINARY_FLOAT",
		"BINARY_DOUBLE", "DOUBLE PRECISION", "FLOAT4", "FLOAT8":
		return &sql.NullFloat64{}

	// 布尔类型
	case "BOOLEAN", "BOOL", "BIT":
		return &sql.NullBool{}

	// 日期时间类型
	case "DATE", "DATETIME", "TIMESTAMP", "TIME", "TIMESTAMPTZ", "TIMESTAMP WITH TIME ZONE",
		"TIMESTAMP WITHOUT TIME ZONE", "TIME WITH TIME ZONE", "TIME WITHOUT TIME ZONE",
		"DATETIME2", "SMALLDATETIME", "YEAR", "INTERVAL":
		return &sql.NullTime{}

	// 二进制类型
	case "BLOB", "BYTEA", "BINARY", "VARBINARY", "LONGBLOB", "MEDIUMBLOB",
		"TINYBLOB", "RAW", "LONG RAW", "BFILE":
		return &sql.RawBytes{}

	// JSON类型
	case "JSON", "JSONB":
		return &sql.NullString{}

	// UUID类型
	case "UUID":
		return &sql.NullString{}

	// 枚举类型（通常作为字符串处理）
	case "ENUM":
		return &sql.NullString{}

	default:
		return nil
	}
}
