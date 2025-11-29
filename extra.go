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
	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/insert-type"
	"github.com/lontten/lorm/return-type"
)

// ExtraContext 扩展参数
type ExtraContext struct {
	insertType       insert_type.InsertType
	returnType       return_type.ReturnType
	allowFullTableOp bool // 是否允许全表操作，默认false不允许全表操作
	showSql          bool
	noRun            bool
	skipSoftDelete   bool
	tableName        string
	selectColumns    []string // select 的 字段名，空为返回所有字段
	limit            *int64
	offset           *int64
	orderByTokens    []string // 排序

	convertCtx ConvertCtx // 查询结果转换函数

	columns      []string
	columnValues []field.Value

	// 唯一索引字段名列表
	duplicateKeyNames []string

	whenUpdateSet *SetContext

	err error
}

func E() *ExtraContext {
	return &ExtraContext{
		whenUpdateSet: Set(),
		convertCtx:    ConvertCtx{}.Init(),
	}
}

// whenUpdateSet 中的错误已经被上抛到 ExtraContext，所以只用判断 ExtraContext 的 err
func (e *ExtraContext) GetErr() error {
	if e.err != nil {
		return e.err
	}
	return nil
}

func (e *ExtraContext) ShowSql() *ExtraContext {
	e.showSql = true
	return e
}

// AllowFullTableOp
// 允许全表操作，默认false不允许全表操作
func (e *ExtraContext) AllowFullTableOp() *ExtraContext {
	e.allowFullTableOp = true
	return e
}

func (e *ExtraContext) NoRun() *ExtraContext {
	e.noRun = true
	return e
}

// Convert
// 查询结果转换函数
func (e *ExtraContext) Convert(c Convert) *ExtraContext {
	e.convertCtx.Add(c)
	return e
}

// SkipSoftDelete
// 跳过软删除
func (e *ExtraContext) SkipSoftDelete() *ExtraContext {
	e.skipSoftDelete = true
	return e
}

func (e *ExtraContext) TableName(name string) *ExtraContext {
	e.tableName = name
	return e
}

// Select
// select 的 字段名，覆盖根据结构体生成的字段列表
func (e *ExtraContext) Select(name ...string) *ExtraContext {
	e.selectColumns = name
	return e
}

func (e *ExtraContext) OrderBy(name string, condition ...bool) *ExtraContext {
	for _, b := range condition {
		if !b {
			return e
		}
	}
	e.orderByTokens = append(e.orderByTokens, name+" ASC")
	return e
}

func (e *ExtraContext) OrderDescBy(name string, condition ...bool) *ExtraContext {
	for _, b := range condition {
		if !b {
			return e
		}
	}
	e.orderByTokens = append(e.orderByTokens, name+" DESC")
	return e
}

func (e *ExtraContext) Limit(num int64, condition ...bool) *ExtraContext {
	for _, b := range condition {
		if !b {
			return e
		}
	}
	e.limit = types.NewInt64(num)
	return e
}

func (e *ExtraContext) Offset(num int64, condition ...bool) *ExtraContext {
	for _, b := range condition {
		if !b {
			return e
		}
	}
	e.offset = types.NewInt64(num)
	return e
}

func (e *ExtraContext) ReturnType(typ return_type.ReturnType) *ExtraContext {
	e.returnType = typ
	return e
}

func (e *ExtraContext) SetNull(name string) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type: field.Null,
	})
	return e
}

func (e *ExtraContext) SetNow(name string) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type: field.Now,
	})
	return e
}

func (e *ExtraContext) Set(name string, value any) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Val,
		Value: value,
	})
	return e
}

// 自增，自减
func (e *ExtraContext) SetIncrement(name string, num any) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Increment,
		Value: num,
	})
	return e
}

// 自定义表达式
// SetExpression("name", "substr(time('now'), 12)") // sqlite 设置时分秒
func (e *ExtraContext) SetExpression(name string, expression string) *ExtraContext {
	e.columns = append(e.columns, name)
	e.columnValues = append(e.columnValues, field.Value{
		Type:  field.Expression,
		Value: expression,
	})
	return e
}

type DuplicateKey struct {
	e *ExtraContext
}

//.whenDuplicateKey(name ...string, )
//.do(nothing, nil)
//.do(update, all, .set(), select ("name", "age"))
//.do(replace, all, .set(), select ("name", "age"))

// WhenDuplicateKey
// 唯一索引冲突,设置索引字段列表；
// Mysql 不用设置，如果设置了也会直接忽略
// Postgresql 必须设置,如果没有设置，则为主键字段
func (e *ExtraContext) WhenDuplicateKey(name ...string) *DuplicateKey {
	e.duplicateKeyNames = name
	return &DuplicateKey{
		e: e,
	}
}

func (dk *DuplicateKey) DoNothing() *ExtraContext {
	dk.e.insertType = insert_type.Ignore
	return dk.e
}

func (dk *DuplicateKey) update(insertType insert_type.InsertType, set ...*SetContext) *ExtraContext {
	dk.e.insertType = insertType
	if len(set) == 0 {
		return dk.e
	}
	if sc := set[0]; sc != nil {
		dk.e.whenUpdateSet = sc
		// err上抛到 ExtraContext
		if sc.err != nil {
			dk.e.err = sc.err
		}
	}
	return dk.e
}

// DoUpdate
// 更新字段未设置时，默认更新所有 有值字段
func (dk *DuplicateKey) DoUpdate(set ...*SetContext) *ExtraContext {
	return dk.update(insert_type.Update, set...)
}

// DoReplace
// 更新字段未设置时，默认更新所有 有值字段
func (dk *DuplicateKey) DoReplace(set ...*SetContext) *ExtraContext {
	return dk.update(insert_type.Replace, set...)
}
