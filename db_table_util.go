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
	"strconv"
	"time"

	"github.com/lontten/lorm/field"
	"github.com/pkg/errors"
)

// init 逻辑删除、租户
func (d *MysqlDialect) initExtra() {
	//if err := lorm.ctx.err; err != nil {
	//	return
	//}
	//
	//if lorm.ctx.ormConf.LogicDeleteYesSql != "" {
	//	lorm.whereTokens = append(lorm.whereTokens, lorm.ctx.ormConf.LogicDeleteYesSql)
	//}
	//
	//if lorm.ctx.ormConf.TenantIdFieldName != "" {
	//	lorm.whereTokens = append(lorm.whereTokens, lorm.ctx.ormConf.TenantIdFieldName)
	//	lorm.args = append(lorm.args, lorm.ctx.ormConf.TenantIdValueFun())
	//}
	//
	//var sb strings.QueryBuild
	//sb.WriteString(lorm.whereSql)
	//
	//if len(lorm.orderByTokens) > 0 {
	//	sb.WriteString(" ORDER BY ")
	//	sb.WriteString(strings.Join(lorm.orderByTokens, ","))
	//}
	//if lorm.limit > 0 {
	//	sb.WriteString(" LIMIT ? ")
	//	lorm.args = append(lorm.args, lorm.limit)
	//}
	//if lorm.offset > 0 {
	//	sb.WriteString(" OFFSET ? ")
	//	lorm.args = append(lorm.args, lorm.offset)
	//}
	//lorm.whereSql = sb.String()

}

// -------------------------utils------------------------
// 获取comp 的 cv
// 排除 nil 字段
func getCompCV(v any, c *OrmConf) ([]string, []field.Value, error) {
	value := reflect.ValueOf(v)
	_, value, err := basePtrDeepValue(value)
	if err != nil {
		return nil, nil, err
	}

	return getCompValueCV(value)
}

// 排除 nil 字段
func getCompValueCV(v reflect.Value) ([]string, []field.Value, error) {
	if !isCompType(v.Type()) {
		return nil, nil, errors.New("getvcv not comp")
	}
	err := checkCompFieldVS(v)
	if err != nil {
		return nil, nil, err
	}

	cv := getStructCVMap(v)
	if len(cv.columns) < 1 {
		return nil, nil, errors.New("where model valid field need ")
	}
	return cv.columns, cv.columnValues, nil
}

//------------------------gen-sql---------------------------

// 根据 columnValues 生成的 VALUES sql
// INSERT INTO table_name (列1, 列2,...) VALUES (值1, 值2,....)
func (ctx *ormContext) genInsertValuesSqlBycolumnValues() {
	columns := ctx.columns
	values := ctx.columnValues
	var query = ctx.query

	for i, v := range values {
		if i > 0 {
			query.WriteString(", ")
		}
		switch v.Type {
		case field.None:
			break
		case field.Null:
			query.WriteString("NULL")
			break
		case field.Now:
			query.WriteString("NOW()")
			break
		case field.UnixSecond:
			query.WriteString(strconv.Itoa(time.Now().Second()))
			break
		case field.UnixMilli:
			query.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))
			break
		case field.UnixNano:
			query.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
			break
		case field.Val:
			query.WriteString("?")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Increment:
			query.WriteString(columns[i] + "+ ?")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Expression:
			query.WriteString(v.Value.(string))
			break
		case field.ID:
			if len(ctx.primaryKeyColumnNames) != 1 {
				ctx.err = errors.New("软删除标记为主键id，需要单主键")
				return
			}
			query.WriteString(ctx.primaryKeyColumnNames[0])
			break
		}
	}
}

// 根据 columnValues 生成的set sql
// SET ...
// column1 = value1, column2 = value2, ...
func (ctx *ormContext) genSetSqlBycolumnValues(fn escapeFun) {
	columns := ctx.columns
	values := ctx.columnValues
	var query = ctx.query

	for i, v := range values {
		column := columns[i]
		column = fn(column)
		if i > 0 {
			query.WriteString(", ")
		}
		switch v.Type {
		case field.None:
			break
		case field.Null:
			query.WriteString(column)
			query.WriteString(" = NULL")
			break
		case field.Now:
			query.WriteString(column)
			query.WriteString(" = NOW()")
			break
		case field.UnixSecond:
			query.WriteString(column)
			query.WriteString(" = ")
			query.WriteString(strconv.Itoa(time.Now().Second()))
			break
		case field.UnixMilli:
			query.WriteString(column)
			query.WriteString(" = ")
			query.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))
			break
		case field.UnixNano:
			query.WriteString(column)
			query.WriteString(" = ")
			query.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
			break
		case field.Val:
			query.WriteString(column)
			query.WriteString(" = ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Increment:
			query.WriteString(column)
			query.WriteString(" = ")
			query.WriteString(column + " + ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Expression:
			query.WriteString(column)
			query.WriteString(" = ")
			query.WriteString(v.Value.(string))
			break
		case field.ID:
			if len(ctx.primaryKeyColumnNames) != 1 {
				ctx.err = errors.New("软删除标记为主键id，需要单主键")
				return
			}
			query.WriteString(column)
			query.WriteString(" = ")
			query.WriteString(ctx.primaryKeyColumnNames[0])
			break
		default:
			ctx.err = errors.New("genSetSqlBycolumnValues not support type")
		}

	}
}
