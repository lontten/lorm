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
	"github.com/lontten/lorm/softdelete"
	"github.com/lontten/lorm/utils"
)

type MysqlDialect struct {
	ctx       *ormContext
	dbVersion MysqlVersion
}

// ===----------------------------------------------------------------------===//
// 获取上下文
// ===----------------------------------------------------------------------===//

func (d *MysqlDialect) getCtx() *ormContext {
	return d.ctx
}
func (d *MysqlDialect) copyContext() Dialecter {
	return &MysqlDialect{
		ctx: &ormContext{
			ormConf:      d.ctx.ormConf,
			convertCtx:   ConvertCtx{}.Init(),
			query:        &strings.Builder{},
			wb:           W(),
			insertType:   insert_type.Err,
			disableColor: d.ctx.disableColor,
		},
		dbVersion: d.dbVersion,
	}
}
func (d *MysqlDialect) hasErr() bool {
	return d.ctx.err != nil
}

func (d *MysqlDialect) getErr() error {
	return d.ctx.err
}

// ===----------------------------------------------------------------------===//
// sql 方言化
// ===----------------------------------------------------------------------===//
func (d *MysqlDialect) query(query string, args ...any) (string, []any) {
	return query, args
}

func (d *MysqlDialect) queryBatch(query string) string {
	return query
}

func (d *MysqlDialect) prepare(query string) string {
	return query
}

func (d *MysqlDialect) exec(query string, args ...any) (string, []any) {
	return query, args
}

func (d *MysqlDialect) execBatch(query string, args [][]any) (string, [][]any) {

	//var num int64 = 0
	//stmt, err := d.lorm.Prepare(query)
	//if err != nil {
	//	return 0, err
	//}
	//for _, arg := range args {
	//	exec, err := stmt.Exec(arg...)
	//	d.log.Println(query, args)
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
// 转义 危险标识符
func (d MysqlDialect) escapeIdentifier(s string) string {
	_, ok := dangNamesMap[s]
	if ok {
		return "`" + s + "`"
	}
	return s
}

// ===----------------------------------------------------------------------===//
// 中间服务
// ===----------------------------------------------------------------------===//

func (d *MysqlDialect) getSql(sql ...string) {
	if len(sql) == 1 {
		d.ctx.originalSql = sql[0]
	} else {
		d.ctx.originalSql = d.ctx.query.String()
	}
	d.ctx.dialectSql = d.ctx.originalSql
}

// insert 生成
func (d *MysqlDialect) tableInsertGen() {
	ctx := d.ctx
	if ctx.hasErr() {
		return
	}

	extra := ctx.extra
	whenUpdateSet := extra.whenUpdateSet

	columns := ctx.columns
	var query = d.ctx.query

	switch ctx.insertType {
	case insert_type.Err:
		query.WriteString("INSERT INTO ")
		break
	case insert_type.Ignore:
		query.WriteString("INSERT IGNORE ")
		break
	case insert_type.Update:
		query.WriteString("INSERT INTO ")
		break
	case insert_type.Replace:
		query.WriteString("REPLACE INTO ")
		break
	}
	query.WriteString(d.escapeIdentifier(ctx.tableName))

	query.WriteString(" (")
	query.WriteString(escapeJoin(d.escapeIdentifier, columns, ", "))
	query.WriteString(") VALUES (")
	ctx.genInsertValuesSqlBycolumnValues()
	query.WriteString(")")

	switch ctx.insertType {
	case insert_type.Update:
		//从 MySQL 8.0.19 开始 可以用 new 取代 VALUES
		//从 MySQL 8.0.20 开始 VALUES 被弃用。
		// INSERT INTO t1 (a,b,c) VALUES (1,2,3),(4,5,6)
		//  ON DUPLICATE KEY UPDATE c=VALUES(c);

		// INSERT INTO t1 (a,b,c) VALUES (1,2,3),(4,5,6) AS new
		//  ON DUPLICATE KEY UPDATE c = new.c;

		if d.dbVersion >= MysqlVersion8_0_19 {
			query.WriteString(" AS new ")
		}

		query.WriteString("ON DUPLICATE KEY UPDATE ")
		// 当未设置更新字段时，默认为所有效有字段（排除索引）
		columnLen := len(whenUpdateSet.columns)
		if columnLen == 0 && len(whenUpdateSet.fieldNames) == 0 {
			list := append(ctx.columns, extra.columns...)

			for _, name := range list {
				find := utils.Find(extra.duplicateKeyNames, name)
				if find < 0 { // 排除 主键 字段
					whenUpdateSet.fieldNames = append(whenUpdateSet.fieldNames, name)
				}
			}
		}

		// 当 软删除 字段 未删除状态 为 0 时，这里fieldNames 会有 软删除字段，
		// DUPLICATE KEY UPDATE 时，软删除字段是否应该更新，问题分析：当更新时：
		// 1.假设 唯一索引 字段 为 name (因为有软删除，逻辑上这样加唯一索引是错误的。) ，更新值为 abc，则数据库中，一定是只有一条 name为 abc的字段，
		// 更新时，会把 软删除设为未删除状态。
		// 1.1 旧数据 未删除，更新数据，符合预期
		// 1.2 旧数据 已删除，更新数据并变成未删除状态，相当于替换，数据insert成功，只是 id 是 原来的数据，之前的旧数据没有了，勉强算是 符合预期
		// 2. 将 name 和 软删除字段（唯一索引设置正确的情况），设为 符合唯一索引，则 数据库中已有数据为 name=abc,del=0; name=abc,del=大于0的数 （已软删除数据众多）
		// 2.1 只有已删除数据，直接插入，符合预期
		// 2.2 同时有未删除，已删除数据，只对未删除数据 进行更新，成功更新，符合预期

		// DUPLICATE KEY UPDATE 时，软删除字段是否应该更新，问题分析：当不更新时：
		// 1.假设 唯一索引 字段 为 name (因为有软删除，逻辑上这样加唯一索引是错误的。) ，更新值为 abc，则数据库中，一定是只有一条 name为 abc的字段，更新时，不会修改软删除状态。
		// 1.1 旧数据 未删除，更新数据，符合预期
		// 1.2 旧数据 已删除，更新数据，已删除被更新，数据还是被删除状态，无法查询到添加的数据，不符合预期！
		// 2. 将 name 和 软删除字段（唯一索引设置正确的情况），设为 符合唯一索引，则 数据库中已有数据为 name=abc,del=0; name=abc,del=大于0的数 （已软删除数据众多）
		// 2.1 只有已删除数据，直接插入，符合预期
		// 2.2 同时有未删除，已删除数据，同时对未删除数据和已删除 进行更新，勉强算是 符合预期

		// 从上面分析可知，DUPLICATE KEY UPDATE 时，软删除字段不进行更新是最差方案，会出现 不符合预期情况。
		// 软删除字段进行更新 时，如果 唯一索引设置正确，是完美执行；如果 唯一索引 错误，也可以达到 基本复合预期的效果。

		for i, name := range whenUpdateSet.fieldNames {
			if i > 0 {
				query.WriteString(", ")
			}
			name = d.escapeIdentifier(name)

			if d.dbVersion >= MysqlVersion8_0_19 {
				query.WriteString(name + " = new." + name)
			} else {
				query.WriteString(name + " = VALUES(" + name + ")")
			}
		}

		for i, column := range whenUpdateSet.columns {
			query.WriteString(d.escapeIdentifier(column) + " = ?")
			if i < columnLen-1 {
				query.WriteString(", ")
			}
			ctx.originalArgs = append(ctx.originalArgs, whenUpdateSet.columnValues[i].Value)
		}
		break
	default:
		break
	}

	query.WriteString(";")
}

// del 生成
func (d *MysqlDialect) tableDelGen() {
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
func (d *MysqlDialect) tableUpdateGen() {
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
func (d *MysqlDialect) tableSelectGen() {
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

func (d *MysqlDialect) parse(c Clause) (string, error) {
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
