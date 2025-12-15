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
	"errors"
	"strconv"
	"strings"

	"github.com/lontten/lorm/insert-type"
	"github.com/lontten/lorm/return-type"
	"github.com/lontten/lorm/softdelete"
	"github.com/lontten/lorm/utils"
)

type PgDialect struct {
	ctx *ormContext
}

// ===----------------------------------------------------------------------===//
// 获取上下文
// ===----------------------------------------------------------------------===//
func (d *PgDialect) getCtx() *ormContext {
	return d.ctx
}
func (d *PgDialect) copyContext() Dialecter {
	return &PgDialect{
		ctx: &ormContext{
			ormConf:      d.ctx.ormConf,
			convertCtx:   ConvertCtx{}.Init(),
			query:        &strings.Builder{},
			wb:           W(),
			insertType:   insert_type.Err,
			disableColor: d.ctx.disableColor,
		},
	}
}
func (d *PgDialect) hasErr() bool {
	return d.ctx.err != nil
}
func (d *PgDialect) getErr() error {
	return d.ctx.err
}

// ===----------------------------------------------------------------------===//
// sql 方言化
// ===----------------------------------------------------------------------===//

func (d *PgDialect) prepare(query string) string {
	query = toPgSql(query)
	return query
}
func (d *PgDialect) exec(query string, args ...any) (string, []any) {
	query = toPgSql(query)
	return query, args
}

func (d *PgDialect) query(query string, args ...any) (string, []any) {
	query = toPgSql(query)
	return query, args
}

func (d *PgDialect) queryBatch(query string) string {
	query = toPgSql(query)

	//return m.lorm.Prepare(query)
	return query
}

// ===----------------------------------------------------------------------===//
// 工具
// ===----------------------------------------------------------------------===//
// 转义 危险标识符
func (d PgDialect) escapeIdentifier(s string) string {
	_, ok := dangNamesMap[s]
	if ok {
		return "\"" + s + "\""
	}
	return s
}
func (d *PgDialect) getSql(sql ...string) {
	if len(sql) == 1 {
		d.ctx.originalSql = sql[0]
	} else {
		d.ctx.originalSql = d.ctx.query.String()
	}
	d.ctx.dialectSql = toPgSql(d.ctx.originalSql)
}

// insert 生成
func (d *PgDialect) tableInsertGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}
	if ctx.insertType == insert_type.Replace {
		ctx.err = errors.New("pg不支持的插入类型 insert-type.Replace")
		return
	}
	extra := ctx.extra
	whenUpdateSet := extra.whenUpdateSet

	columns := ctx.columns
	var query = d.ctx.query

	query.WriteString("INSERT INTO")
	query.WriteString(" ")

	query.WriteString(d.escapeIdentifier(ctx.tableName))

	query.WriteString(" (")
	query.WriteString(escapeJoin(d.escapeIdentifier, columns, ", "))
	query.WriteString(") VALUES (")
	ctx.genInsertValuesSqlBycolumnValues()
	query.WriteString(")")

	if ctx.insertType == insert_type.Ignore || ctx.insertType == insert_type.Update {
		query.WriteString(" ON CONFLICT (")
		if len(extra.duplicateKeyNames) > 0 {
			query.WriteString(escapeJoin(d.escapeIdentifier, extra.duplicateKeyNames, ","))
		} else {
			query.WriteString(escapeJoin(d.escapeIdentifier, ctx.primaryKeyColumnNames, ","))
		}
		query.WriteString(") DO ")
	}

	switch ctx.insertType {
	case insert_type.Ignore:
		query.WriteString("NOTHING")
		break
	case insert_type.Update:
		query.WriteString("UPDATE SET")
		query.WriteString(" ")

		// 当未设置更新字段时，默认为所有字段
		if len(whenUpdateSet.columns) == 0 && len(whenUpdateSet.fieldNames) == 0 {
			list := append(ctx.columns, extra.columns...)

			for _, name := range list {
				find := utils.Find(extra.duplicateKeyNames, name)
				if find < 0 { // 排除 主键 字段
					whenUpdateSet.fieldNames = append(whenUpdateSet.fieldNames, name)
				}
			}
		}

		for i, name := range whenUpdateSet.fieldNames {
			name = d.escapeIdentifier(name)
			query.WriteString(name + " = EXCLUDED." + name)
			if i < len(whenUpdateSet.fieldNames)-1 {
				query.WriteString(", ")
			}
		}

		for i, column := range whenUpdateSet.columns {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString(d.escapeIdentifier(column) + " = ?")
			ctx.originalArgs = append(ctx.originalArgs, whenUpdateSet.columnValues[i].Value)
		}
		break
	default:
		break
	}

	// 当scan为指针类型时，返回字段。
	if ctx.returnAutoPrimaryKey != pkNoReturn {
		switch expr := ctx.returnType; expr {
		case return_type.None:
			break
		case return_type.Auto:
			query.WriteString(" RETURNING " + escapeJoin(d.escapeIdentifier, ctx.allAutoColumnNames, ","))
		case return_type.ZeroField:
			query.WriteString(" RETURNING " + escapeJoin(d.escapeIdentifier, ctx.modelZeroColumnNames, ","))
		case return_type.AllField:
			query.WriteString(" RETURNING " + escapeJoin(d.escapeIdentifier, ctx.modelAllColumnNames, ","))
		}
	}
	query.WriteString(";")
}

// del 生成
func (d *PgDialect) tableDelGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}
	var query = d.ctx.query
	tableName := ctx.tableName

	whereStr, args, err := ctx.wb.toSql(d.parse, ctx.primaryKeyColumnNames...)
	if err != nil {
		ctx.err = err
		return
	}

	if !ctx.allowFullTableOp {
		if whereStr == "" {
			ctx.err = errors.New("禁止全表操作")
			return
		}
	}

	//  没有软删除 或者 跳过软删除 ，执行物理删除
	if ctx.softDeleteType == softdelete.None || ctx.skipSoftDelete {
		query.WriteString("DELETE FROM ")
		query.WriteString(d.escapeIdentifier(tableName))
	} else {
		query.WriteString("UPDATE ")
		query.WriteString(d.escapeIdentifier(tableName))

		query.WriteString(" SET ")
		ctx.genSetSqlBycolumnValues(d.escapeIdentifier)
	}
	if len(whereStr) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(whereStr)
	}
	ctx.originalArgs = append(ctx.originalArgs, args...)

	query.WriteString(";")
}

// update 生成
func (d *PgDialect) tableUpdateGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}
	var query = d.ctx.query
	tableName := ctx.tableName
	whereStr, args, err := ctx.wb.toSql(d.parse, ctx.primaryKeyColumnNames...)
	if err != nil {
		ctx.err = err
		return
	}

	if !ctx.allowFullTableOp {
		if whereStr == "" {
			ctx.err = errors.New("禁止全表操作")
			return
		}
	}

	query.WriteString("UPDATE ")
	query.WriteString(d.escapeIdentifier(tableName))
	query.WriteString(" SET ")
	ctx.genSetSqlBycolumnValues(d.escapeIdentifier)

	if len(whereStr) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(whereStr)
	}

	ctx.originalArgs = append(ctx.originalArgs, args...)
	query.WriteString(";")
}

// select 生成
func (d *PgDialect) tableSelectGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}
	var query = d.ctx.query
	tableName := ctx.tableName
	whereStr, args, err := ctx.wb.toSql(d.parse)
	if err != nil {
		ctx.err = err
		return
	}
	if !ctx.allowFullTableOp {
		if whereStr == "" {
			ctx.err = errors.New("禁止全表操作")
			return
		}
	}

	query.WriteString("SELECT ")
	query.WriteString(escapeJoin(d.escapeIdentifier, ctx.modelSelectFieldNames, " ,"))
	query.WriteString(" FROM ")
	query.WriteString(tableName)

	if len(whereStr) > 0 {
		query.WriteString(" WHERE ")
		query.WriteString(whereStr)
	}

	ctx.originalArgs = append(ctx.originalArgs, args...)
	query.WriteString(ctx.lastSql)
	if ctx.limit != nil {
		query.WriteString(" LIMIT ")
		query.WriteString(strconv.FormatInt(*ctx.limit, 10))
	}
	if ctx.offset != nil {
		query.WriteString(" OFFSET ")
		query.WriteString(strconv.FormatInt(*ctx.offset, 10))
	}

	query.WriteString(";")
}

func (d *PgDialect) execBatch(query string, args [][]any) (string, [][]any) {
	query = toPgSql(query)
	//var num int64 = 0
	//stmt, err := m.lorm.Prepare(query)
	//defer stmt.Close()
	//if err != nil {
	//	return 0, err
	//}
	//for _, arg := range args {
	//	exec, err := stmt.Exec(arg...)
	//
	//	m.log.Println(query, arg)
	//	if err != nil {
	//		return num, err
	//	}
	//	rowsAffected, err := exec.RowsAffected()
	//	if err != nil {
	//		return num, err
	//	}
	//	num += rowsAffected
	//}
	return query, args
}

// ===----------------------------------------------------------------------===//
// 工具
// ===----------------------------------------------------------------------===//

func toPgSql(sql string) string {
	var i = 1
	for {
		t := strings.Replace(sql, "?", "$"+strconv.Itoa(i), 1)
		if t == sql {
			break
		}
		i++
		sql = t
	}
	return sql
}

// ===----------------------------------------------------------------------===//
// 中间服务
// ===----------------------------------------------------------------------===//
func (d *PgDialect) toSqlInsert() (string, []any) {
	tableName := d.ctx.tableName
	return tableName, nil
}

func (d *PgDialect) parse(c Clause) (string, error) {
	sb := strings.Builder{}
	switch c.Type {
	case Eq:
		sb.WriteString(c.query + " = ?")
	case Neq:
		sb.WriteString(c.query + " <> ?")
	case Less:
		sb.WriteString(c.query + " < ?")
	case LessEq:
		sb.WriteString(c.query + " <= ?")
	case Greater:
		sb.WriteString(c.query + " > ?")
	case GreaterEq:
		sb.WriteString(c.query + " >= ?")
	case Like:
		sb.WriteString(c.query + " LIKE ?")
	case NotLike:
		sb.WriteString(c.query + " NOT LIKE ?")
	case In:
		length := len(c.args)
		if length == 0 {
			sb.WriteString("1=0")
		} else {
			sb.WriteString(c.query + " IN (")
			sb.WriteString(gen(length))
			sb.WriteString(")")
		}
	case NotIn:
		sb.WriteString(c.query + " NOT IN (")
		sb.WriteString(gen(len(c.args)))
		sb.WriteString(")")
	case Between:
		sb.WriteString(c.query + " BETWEEN ? AND ?")
	case NotBetween:
		sb.WriteString(c.query + " NOT BETWEEN ? AND ?")
	case IsNull:
		sb.WriteString(c.query + " IS NULL")
	case IsNotNull:
		sb.WriteString(c.query + " IS NOT NULL")
	case IsFalse:
		sb.WriteString(c.query + " IS FALSE")
	default:
		return "", errors.New("unknown where token type")
	}

	return sb.String(), nil
}
