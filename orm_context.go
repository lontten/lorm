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
	"strings"

	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/insert-type"
	"github.com/lontten/lorm/return-type"
	"github.com/lontten/lorm/softdelete"
	"github.com/lontten/lorm/sqltype"
	"github.com/lontten/lorm/utils"
	"github.com/pkg/errors"
)

type tableSqlType int

// d前缀是单表的意思，tableSqlType 只用于单表操作
const (
	dInsert tableSqlType = iota
	dInsertOrUpdate
	dUpdate
	dDelete
	dSelect
	dGetOrInsert
	dHas
	dCount
)

type returnAutoPrimaryKeyType int

// pk前缀是主键的意思
const (
	pkNoReturn    returnAutoPrimaryKeyType = iota
	pkQueryReturn                          // insert 时，可以直接 query 返回
	pkFetchReturn                          // insert 时，不能直接返回，需要手动LastInsertId获取
)

type ormContext struct {
	ormConf *OrmConf
	extra   *ExtraContext

	allowFullTableOp bool // 是否允许全表操作，默认false不允许全表操作

	// model 参数，用于校验字段类型是否合法
	paramModelBaseV reflect.Value

	// dest
	scanDest  any
	scanIsPtr bool
	scanV     reflect.Value // dest 去除ptr的value

	destBaseValue      reflect.Value // list第一个，去除ptr
	destBaseType       reflect.Type
	destBaseTypeIsComp bool
	// scan 为slice时，里面item是否是ptr
	destIsSlice        bool
	destSliceItemIsPtr bool

	log Logger
	err error

	tableSqlType tableSqlType //单表，sql类型crud

	isLgDel bool //是否启用了逻辑删除
	isTen   bool //是否启用了多租户

	// ------------------主键----------------------
	indexs []Index // 索引列表

	returnAutoPrimaryKey returnAutoPrimaryKeyType // 自增主键返回类型

	// 在不支持 insertCanReturn 的数据库中，使用 LastInsertId 返回 自增主键
	// First时，用来当默认排序字段
	autoPrimaryKeyColumnName string // 自增主键字段名
	autoPrimaryKeyFieldName  string // 自增主键字段名
	// 只能在 insert时，返回字段，只能支持 insertCanReturn 的数据库，可以返回
	otherAutoColumnNames []string // 其他自动生成字段名列表
	allAutoColumnNames   []string // 所有自动生成字段名列表

	autoPrimaryKeyFieldIsPtr    bool         // id 对应的model字段 是否是 ptr
	autoPrimaryKeyFieldBaseType reflect.Type // id 对应的model字段 type

	// id = 1
	//主键名-列表,这里考虑到多主键
	primaryKeyColumnNames []string

	// ------------------conf----------------------

	insertType     insert_type.InsertType
	returnType     return_type.ReturnType
	softDeleteType softdelete.SoftDelType
	skipSoftDelete bool       // 跳过软删除
	tableName      string     //当前表名
	checkParam     bool       // 是否检查参数
	showSql        bool       // 打印sql
	disableColor   bool       // 打印sql时，是否使用颜色
	noRun          bool       // 不实际执行
	convertCtx     ConvertCtx // 查询结果转换函数
	// ------------------conf-end----------------------

	// ------------------字段名：字段值----------------------

	columns      []string      // 有效字段 column
	columnValues []field.Value // 有效字段 value

	modelZeroColumnNames      []string // model 零值字段列表
	modelNoSoftDelColumnNames []string // model 所有字段列表- 忽略软删除字段
	modelAllColumnNames       []string // model 所有字段列表
	modelSelectFieldNames     []string // model select 字段列表
	// ------------------字段名：字段值-end----------------------

	//------------------scan----------------------
	sqlType sqltype.SqlType

	query       *strings.Builder // query sql
	originalSql string           // 原始sql
	dialectSql  string           // 方言 sql
	//参数
	originalArgs []any // 原始参数

	started bool

	wb       *WhereBuilder
	whereSql string // WhereBuilder 生成的 where sql
	lastSql  string // 最后拼接的sql
	limit    *int64
	offset   *int64
}

func (ctx *ormContext) setLastInsertId(lastInsertId int64) {
	var vp reflect.Value
	switch ctx.autoPrimaryKeyFieldBaseType.Kind() {
	case reflect.Int8:
		id := int8(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Int16:
		id := int16(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Int32:
		id := int32(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Int64:
		id := int64(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Int:
		id := int(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Uint:
		id := uint(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Uint8:
		id := uint8(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Uint16:
		id := uint16(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Uint32:
		id := uint32(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	case reflect.Uint64:
		id := uint64(lastInsertId)
		vp = reflect.ValueOf(&id)
		break
	default:
		ctx.err = errors.New("last_insert_id field type error")
		return
	}
	f := ctx.destBaseValue.FieldByName(ctx.autoPrimaryKeyFieldName)
	if ctx.autoPrimaryKeyFieldIsPtr {
		f.Set(vp)
	} else {
		f.Set(reflect.Indirect(vp))
	}
}
func (ctx *ormContext) initExtra(extra ...*ExtraContext) {
	var e *ExtraContext
	if len(extra) > 0 && extra[0] != nil {
		e = extra[0]
	} else {
		e = E()
	}
	// err 上抛到 ormContext
	if e.GetErr() != nil {
		ctx.err = e.GetErr()
		return
	}
	ctx.allowFullTableOp = e.allowFullTableOp
	ctx.convertCtx = e.convertCtx
	ctx.extra = e
	ctx.insertType = e.insertType
	ctx.returnType = e.returnType
	ctx.showSql = e.showSql
	ctx.noRun = e.noRun
	ctx.skipSoftDelete = e.skipSoftDelete
	ctx.tableName = e.tableName

	ctx.modelSelectFieldNames = e.selectColumns

	if len(e.orderByTokens) > 0 {
		ctx.lastSql += " ORDER BY " + strings.Join(e.orderByTokens, ",")
	}
	if ctx.limit == nil {
		ctx.limit = e.limit
	}
	ctx.offset = e.offset
}

// 初始化 表名,主键，自增id
func (ctx *ormContext) initConf() {
	if ctx.hasErr() {
		return
	}

	v := ctx.destBaseValue

	dest := ctx.scanDest
	t := ctx.destBaseType
	ctx.softDeleteType = utils.GetSoftDelType(t)

	if ctx.tableName == "" {
		ctx.tableName = ctx.ormConf.tableName(v, dest)
		if ctx.tableName == "" {
			ctx.err = ErrNoTableName
			return
		}
	}

	ctx.primaryKeyColumnNames = ctx.ormConf.primaryKeyColumnNames(v, dest)

	tc := getTableConf(v)
	if tc != nil {
		ctx.autoPrimaryKeyColumnName = tc.autoPrimaryKeyColumnName
		ctx.otherAutoColumnNames = tc.otherAutoColumnName
		ctx.allAutoColumnNames = tc.allAutoColumnName
	}
}

// 获取struct对应的字段名 和 其值，
// slice为全部，一个为非nil字段。
func (ctx *ormContext) initColumnsValue() {
	if ctx.hasErr() {
		return
	}

	cv := getStructCVMap(ctx.destBaseValue)
	ctx.columns = cv.columns
	ctx.columnValues = cv.columnValues

	ctx.modelZeroColumnNames = cv.modelZeroColumnNames
	ctx.modelNoSoftDelColumnNames = cv.modelAllColumnNames
	ctx.modelAllColumnNames = cv.modelAllColumnNames

	if ctx.scanIsPtr && ctx.returnType != return_type.None {
		if ctx.ormConf.insertCanReturn {
			ctx.returnAutoPrimaryKey = pkQueryReturn
		} else if ctx.autoPrimaryKeyColumnName != "" {
			ctx.returnAutoPrimaryKey = pkFetchReturn
		}
	}

	if ctx.returnAutoPrimaryKey == pkFetchReturn {
		fieldName, ok := cv.modelAllCFNameMap[ctx.autoPrimaryKeyColumnName]
		if !ok {
			ctx.err = errors.New("TableConfContext not set AutoPrimaryKey")
			return
		}
		ctx.autoPrimaryKeyFieldName = fieldName

		structField, _ := ctx.destBaseType.FieldByName(fieldName)
		isPtr, baseT := basePtrType(structField.Type)
		ctx.autoPrimaryKeyFieldIsPtr = isPtr
		ctx.autoPrimaryKeyFieldBaseType = baseT
	}
	return
}
func (ctx *ormContext) initColumns() {
	if ctx.hasErr() {
		return
	}
	if len(ctx.modelSelectFieldNames) == 0 {
		columns := getStructCList(ctx.destBaseType)
		ctx.modelSelectFieldNames = columns
	}
	return
}

func (ctx *ormContext) initColumnsValueExtra() {
	if ctx.hasErr() {
		return
	}
	e := ctx.extra
	whenUpdateSet := e.whenUpdateSet
	if whenUpdateSet.hasModel {
		oc := &ormContext{
			ormConf:        ctx.ormConf,
			skipSoftDelete: true,
		}
		oc.initTargetDestOne(whenUpdateSet.model) //初始化参数
		oc.initColumnsValue()                     //初始化cv

		whenUpdateSet.columns = append(whenUpdateSet.columns, oc.columns...)
		whenUpdateSet.columnValues = append(whenUpdateSet.columnValues, oc.columnValues...)
	}
	if ctx.hasErr() {
		return
	}
	for i, column := range e.columns {
		cv := e.columnValues[i]
		if cv.Type == field.Null || cv.Type == field.Now {
			ctx.modelZeroColumnNames = append(ctx.modelZeroColumnNames, column)
		}
		find := utils.Find(ctx.columns, column)
		if find == -1 {
			ctx.columns = append(ctx.columns, column)
			ctx.columnValues = append(ctx.columnValues, cv)
		} else {
			ctx.columnValues[find] = cv
		}
	}
	return
}
func (ctx *ormContext) initColumnsValueSoftDel() {
	if ctx.hasErr() {
		return
	}
	if ctx.skipSoftDelete {
		return
	}

	switch ctx.sqlType {
	case sqltype.Insert:
		value, has := softdelete.SoftDelTypeNoFVMap[ctx.softDeleteType]
		if has {
			ctx.columns = append(ctx.columns, value.Name)
			ctx.columnValues = append(ctx.columnValues, value.ToValue())
		}
		break
	case sqltype.Delete:
		if ctx.softDeleteType != softdelete.None {
			// set
			value, has := softdelete.SoftDelTypeYesFVMap[ctx.softDeleteType]
			if has {
				ctx.columns = append(ctx.columns, value.Name)
				ctx.columnValues = append(ctx.columnValues, value.ToValue())
			}

			// where
			value, has = softdelete.SoftDelTypeNoFVMap[ctx.softDeleteType]
			if has {
				ctx.wb.fieldValue(value.Name, value.ToValue())
			}
		}
		break
	case sqltype.Update:
		if ctx.softDeleteType != softdelete.None {
			// where
			value, has := softdelete.SoftDelTypeNoFVMap[ctx.softDeleteType]
			if has {
				ctx.wb.fieldValue(value.Name, value.ToValue())
			}
		}
		break
	case sqltype.Select:
		if ctx.softDeleteType != softdelete.None {
			// select
			value, has := softdelete.SoftDelTypeYesFVMap[ctx.softDeleteType]
			if has {
				ctx.modelSelectFieldNames = append(ctx.modelSelectFieldNames, value.Name)
			}

			// where
			value, has = softdelete.SoftDelTypeNoFVMap[ctx.softDeleteType]
			if has {
				ctx.wb.fieldValue(value.Name, value.ToValue())
			}
		}
		break

	default:
		break
	}
	return
}

func (ctx ormContext) Copy() ormContext {
	return ormContext{
		ormConf: ctx.ormConf,
		log:     ctx.log,
	}
}

func (ctx ormContext) printSql() {
	if ctx.showSql {
		utils.PrintSql(ctx.disableColor, ctx.originalSql, ctx.dialectSql, ctx.originalArgs...)
	}
}

func (ctx *ormContext) hasErr() bool {
	return ctx.err != nil
}
