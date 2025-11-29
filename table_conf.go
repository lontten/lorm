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

	"github.com/lontten/lcore/v2/lcutils"
)

type TableConfContext struct {
	tableName                string   // 表名
	primaryKeyColumnNames    []string // 主键字段名列表
	autoPrimaryKeyColumnName string   // 自增主键字段名
	otherAutoColumnName      []string // 其他自动生成字段名列表
	allAutoColumnName        []string // 全部自动生成字段名列表
	indexs                   []Index  // 数据库索引列表
}

func TableConf(name string) *TableConfContext {
	return &TableConfContext{
		tableName: name,
	}
}

type Index struct {
	Name      string   // 索引名称
	Unique    bool     // 是否唯一
	Columns   []string // 索引列
	IndexType string   // 索引类型
	Comment   string   // 索引注释
}

// Table 设置表名
func (c *TableConfContext) Table(name string) *TableConfContext {
	c.tableName = name
	return c
}

// PrimaryKeys 设置主键字段，多个字段为复合主键
func (c *TableConfContext) PrimaryKeys(name ...string) *TableConfContext {
	c.primaryKeyColumnNames = name
	c.initAutoPk()
	return c
}

// AutoColumn 会在数据库自动生成的字段
// 例如：
// 自增字段、虚拟列、计算列、默认值，等
// 在insert时，可以设置返回这些字段
func (c *TableConfContext) AutoColumn(name ...string) *TableConfContext {
	c.allAutoColumnName = name
	c.initAutoPk()
	return c
}

// AutoColumn 会在数据库自动生成的字段
// 例如：
// 自增字段、虚拟列、计算列、默认值，等
// 在insert时，可以设置返回这些字段
func (c *TableConfContext) initAutoPk() {
	if len(c.primaryKeyColumnNames) == 0 {
		return
	}
	if len(c.allAutoColumnName) == 0 {
		return
	}
	list := lcutils.BoolIntersection(c.allAutoColumnName, c.primaryKeyColumnNames)
	if len(list) == 0 {
		c.otherAutoColumnName = c.allAutoColumnName
		return
	}
	// 一般数据库，auto的主键字段，一般都是第一个；但是pg可以用 默认uuid+复合主键，实现多auto字段的复合主键（但是这种特殊情况不考虑）
	// 直接取第一个字段作为 autoPrimaryKeyColumnName
	// 对于 pg的特殊情况，也只取第一个字段，其他字段放在 otherAutoColumnName，不影响 lorm正常使用
	c.autoPrimaryKeyColumnName = list[0]

	c.otherAutoColumnName = lcutils.BoolDiff(c.allAutoColumnName, list[:1])
	return
}

var TableConfCache = map[reflect.Type]TableConfContext{}

func getTableConf(v reflect.Value) *TableConfContext {
	n, has := TableConfCache[v.Type()]
	if has {
		return &n
	}
	method := v.MethodByName("TableConf")
	if !method.IsValid() || method.IsZero() {
		return nil
	}

	values := method.Call(nil)

	if len(values) != 1 {
		return nil
	}
	value := values[0]
	if value.IsNil() {
		return nil
	}
	tc, ok := value.Interface().(*TableConfContext)
	if !ok {
		return nil
	}
	return tc
}

func getTableName(v reflect.Value) string {
	tc := getTableConf(v)
	if tc == nil {
		return ""
	}
	return tc.tableName
}

func getPrimaryKeyColumnNames(v reflect.Value) []string {
	tc := getTableConf(v)
	if tc == nil {
		return nil
	}
	return tc.primaryKeyColumnNames
}
