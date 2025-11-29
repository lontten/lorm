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

	"github.com/lontten/lcore/v2/types"
	"github.com/lontten/lorm/utils"
)

func QueryBuild[T any](db Engine) *SqlBuilder[T] {
	return &SqlBuilder[T]{
		db:              db.init(),
		selectQuery:     &strings.Builder{},
		otherSqlBuilder: &strings.Builder{},
	}
}

type SqlBuilder[T any] struct {
	db Engine
	// 最终执行sql
	query string
	// 最终参数列表
	args []any

	// select部分;分页时，需要将select 部分换成 count，这时如果select部分有 arg需要也一起去掉。
	// 所有 select/where 两部分的arg时分开的，需要 selectStatus 来区分，arg放在哪里。
	selectStatus int8

	selectTokens []string
	orderTokens  []string
	selectQuery  *strings.Builder
	selectArgs   []any

	// 分页
	countField   string // 分页时，用于查询总数的字段，默认为 *
	fakeTotalNum int64  // 分页-假数据总数，分页时，跳过查询，直接使用这个总数;默认-1，表示未设置
	noGetList    bool   // 分页-不返回数据，分页时，只查询数量，不返回数据列表;默认false，表示需要获取数据

	// 其他部分
	otherSqlBuilder *strings.Builder
	otherSqlArgs    []any
	// 用来判断是否需要添加 where
	whereStatus int8

	pageConfig *PageConfig
}

const (
	selectNoSet = iota
	selectSet

	selectDone
)

const (
	whereNoSet = iota
	whereSet
	whereDone
)

func (b *SqlBuilder[T]) initSelectSql() {
	dialect := b.db.getDialect()
	if len(b.selectTokens) > 0 {
		b.selectQuery.WriteString("SELECT ")
	}
	b.selectQuery.WriteString(escapeJoin(dialect.escapeIdentifier, b.selectTokens, " ,"))
	b.query = b.selectQuery.String() + " " + b.otherSqlBuilder.String()
	if len(b.orderTokens) > 0 {
		b.query = b.query + " ORDER BY " + escapeJoin(dialect.escapeIdentifier, b.orderTokens, " ,")
	}
	b.args = append(b.selectArgs, b.otherSqlArgs...)
}

// 显示sql
func (b *SqlBuilder[T]) ShowSql(conditions ...bool) *SqlBuilder[T] {
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	b.db.getCtx().showSql = true
	return b
}

// 不执行
func (b *SqlBuilder[T]) NoRun(conditions ...bool) *SqlBuilder[T] {
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	b.db.getCtx().noRun = true
	return b
}

// Convert
// 查询结果转换函数
func (b *SqlBuilder[T]) Convert(c Convert, conditions ...bool) *SqlBuilder[T] {
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	b.db.getCtx().convertCtx.Add(c)
	return b
}

// 添加一个 arg，多个断言
func (b *SqlBuilder[T]) AppendArg(arg any, conditions ...bool) *SqlBuilder[T] {
	if b.db.getCtx().hasErr() {
		return b
	}
	for _, c := range conditions {
		if !c {
			return b
		}
	}
	if b.selectStatus == selectNoSet {
		b.selectArgs = append(b.selectArgs, arg)
	} else {
		b.otherSqlArgs = append(b.otherSqlArgs, arg)
	}
	return b
}

// 添加sql语句
func (b *SqlBuilder[T]) AppendSql(sql string) *SqlBuilder[T] {
	b.otherSqlBuilder.WriteString(sql)
	return b
}

// 添加 多个参数
func (b *SqlBuilder[T]) AppendArgs(args ...any) *SqlBuilder[T] {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if b.selectStatus == selectDone {
		b.otherSqlArgs = append(b.otherSqlArgs, args...)
	} else {
		b.selectArgs = append(b.selectArgs, args...)
	}
	return b
}

// 添加一个 select 字段，多个断言
func (b *SqlBuilder[T]) Select(arg string, condition ...bool) *SqlBuilder[T] {
	ctx := b.db.getCtx()
	if b.selectStatus == selectDone {
		ctx.err = errors.New("Select 代码位置异常")
		return b
	}
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}

	b.selectStatus = selectSet
	b.selectTokens = append(b.selectTokens, arg)

	return b
}

// 添加 多个 select 字段，从 model中
func (b *SqlBuilder[T]) SelectModel(v any, condition ...bool) *SqlBuilder[T] {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	ctx := b.db.getCtx()
	if b.selectStatus == selectDone {
		ctx.err = errors.New("SelectModel 代码位置异常")
		return b
	}
	if ctx.hasErr() {
		return b
	}
	if v == nil {
		return b
	}

	ctx.initScanDestOne(v)
	columns := getStructCAllList(ctx.destBaseType)

	b.selectStatus = selectSet
	b.selectTokens = append(b.selectTokens, columns...)
	return b
}

// from 表名
// 状态从 selectNoSet 变成 selectSet
func (b *SqlBuilder[T]) From(name string) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	b.otherSqlBuilder.WriteString(" FROM " + name)
	return b
}

// join 联表
func (b *SqlBuilder[T]) Join(name string, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" JOIN " + name)
	return b
}

func (b *SqlBuilder[T]) Arg(arg any, condition ...bool) *SqlBuilder[T] {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.AppendArgs(arg)
	return b
}

func (b *SqlBuilder[T]) Args(args ...any) *SqlBuilder[T] {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder[T]) LeftJoin(name string, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString("\n")
	b.otherSqlBuilder.WriteString("LEFT JOIN " + name)
	b.otherSqlBuilder.WriteString("\n")

	return b
}

func (b *SqlBuilder[T]) RightJoin(name string, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" RIGHT JOIN " + name)
	return b
}

func (b *SqlBuilder[T]) OrderBy(name string, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.orderTokens = append(b.orderTokens, name+" ASC")
	return b
}

func (b *SqlBuilder[T]) OrderDescBy(name string, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.orderTokens = append(b.orderTokens, name+" DESC")
	return b
}
func (b *SqlBuilder[T]) Native(sql string, condition ...bool) *SqlBuilder[T] {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" ")
	b.otherSqlBuilder.WriteString(sql)
	b.otherSqlBuilder.WriteString(" ")
	return b
}

func (b *SqlBuilder[T]) Limit(num int64, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" LIMIT " + strconv.FormatInt(num, 10))
	return b
}

func (b *SqlBuilder[T]) Offset(num int64, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}

	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" OFFSET " + strconv.FormatInt(num, 10))
	return b
}

func (b *SqlBuilder[T]) WhereBuilder(w *WhereBuilder) *SqlBuilder[T] {
	b.selectStatus = selectDone
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if w == nil {
		return b
	}
	sqlStr, args, err := w.toSql(b.db.getDialect().parse)
	if err != nil {
		b.db.getCtx().err = err
		return b
	}
	if sqlStr == "" {
		return b
	}
	sqlStr = "(" + sqlStr + ")"

	switch b.whereStatus {
	case whereNoSet:
		b.whereStatus = whereSet
		b.otherSqlBuilder.WriteString(" WHERE ")
		b.otherSqlBuilder.WriteString(sqlStr)
	case whereSet:
		b.otherSqlBuilder.WriteString(" AND ")
		b.otherSqlBuilder.WriteString(sqlStr)
	case whereDone:
		b.db.getCtx().err = errors.New("where has been done")
	}

	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder[T]) Where(whereStr string, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b._whereArg(whereStr)
	return b
}

func (b *SqlBuilder[T]) _whereArg(whereStr string, args ...any) *SqlBuilder[T] {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	b.AppendArgs(args...)
	switch b.whereStatus {
	case whereNoSet:
		b.whereStatus = whereSet
		b.otherSqlBuilder.WriteString(" WHERE ")
		b.otherSqlBuilder.WriteString(whereStr)
	case whereSet:
		b.otherSqlBuilder.WriteString(" AND ")
		b.otherSqlBuilder.WriteString(whereStr)
	case whereDone:
		ctx.err = errors.New("where has been done")
	}

	return b
}
func (b *SqlBuilder[T]) BoolWhere(condition bool, whereStr string, args ...any) *SqlBuilder[T] {
	b.selectStatus = selectDone
	if !condition {
		return b
	}
	b._whereArg(whereStr, args...)
	return b
}

func (b *SqlBuilder[T]) BoolWhereIn(condition bool, whereStr string, args ...any) *SqlBuilder[T] {
	b.selectStatus = selectDone
	if !condition {
		return b
	}
	b.WhereIn(whereStr, args...)
	return b
}

func (b *SqlBuilder[T]) WhereIn(whereStr string, args ...any) *SqlBuilder[T] {
	b.selectStatus = selectDone
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	length := len(args)
	if length == 0 {
		whereStr = "1=0"
	} else {
		b.AppendArgs(args...)

		var inArgStr = " (" + gen(length) + ")"
		whereStr = whereStr + " IN" + inArgStr
	}

	switch b.whereStatus {
	case whereNoSet:
		b.whereStatus = whereSet
		b.otherSqlBuilder.WriteString(" WHERE ")

		b.otherSqlBuilder.WriteString(whereStr)

	case whereSet:
		b.otherSqlBuilder.WriteString(" AND ")

		b.otherSqlBuilder.WriteString(whereStr)

	case whereDone:
		ctx.err = errors.New("where has been done")
	}

	return b
}

// BoolWhereSqlIn
// in ? 当参数列表长度为0时，为 1=0   false条件
func (b *SqlBuilder[T]) BoolWhereSqlIn(condition bool, whereStr string, args ...any) *SqlBuilder[T] {
	b.selectStatus = selectDone
	if !condition {
		return b
	}
	b.WhereSqlIn(whereStr, args...)
	return b
}

// WhereSqlIn
// in ? 当参数列表长度为0时，为 1=0   false条件
func (b *SqlBuilder[T]) WhereSqlIn(whereStr string, args ...any) *SqlBuilder[T] {
	b.selectStatus = selectDone
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	length := len(args)
	if length == 0 {
		whereStr = "1=0"
	} else {
		b.AppendArgs(args...)

		var inArgStr = " (" + gen(length) + ")"
		whereStr = strings.Replace(whereStr, "?", inArgStr, -1)
	}

	switch b.whereStatus {
	case whereNoSet:
		b.whereStatus = whereSet
		b.otherSqlBuilder.WriteString(" WHERE ")

		b.otherSqlBuilder.WriteString(whereStr)

	case whereSet:
		b.otherSqlBuilder.WriteString(" AND ")

		b.otherSqlBuilder.WriteString(whereStr)

	case whereDone:
		ctx.err = errors.New("where has been done")
	}

	return b
}

func (b *SqlBuilder[T]) Between(whereStr string, begin, end any, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	for _, c := range condition {
		if !c {
			return b
		}
	}
	has1 := !utils.IsNil(begin)
	has2 := !utils.IsNil(end)

	if has1 {
		if has2 {
			b._whereArg(whereStr+" BETWEEN ? AND ?", begin, end)
			return b
		}
		b._whereArg(whereStr+" >= ?", begin)
		return b
	}
	if has2 {
		b._whereArg(whereStr+" <= ?", end)
		return b
	}
	return b
}

func (b *SqlBuilder[T]) Like(key *string, fields ...string) *SqlBuilder[T] {
	b.selectStatus = selectDone
	b._like(key, 1, fields...)
	return b
}
func (b *SqlBuilder[T]) LikeLeft(key *string, fields ...string) *SqlBuilder[T] {
	b.selectStatus = selectDone
	b._like(key, 2, fields...)
	return b
}
func (b *SqlBuilder[T]) LikeRight(key *string, fields ...string) *SqlBuilder[T] {
	b.selectStatus = selectDone
	b._like(key, 3, fields...)
	return b
}

// likeType
// 1 表示 like '%key%';
// 2 表示 like '%key';
// 3 表示 like 'key%';
func (b *SqlBuilder[T]) _like(key *string, likeType int, fields ...string) *SqlBuilder[T] {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if b.selectStatus != selectDone {
		ctx.err = errors.New("Where 设置异常：like ")
		return b
	}

	if types.NilToZero(key) == "" {
		return b
	}
	if len(fields) == 0 {
		return b
	}
	var args []any
	var k = ""
	if likeType == 1 {
		k = "%" + *key + "%"
	} else if likeType == 2 {
		k = "%" + *key
	} else if likeType == 3 {
		k = *key + "%"
	}

	var tokens []string
	for _, field := range fields {
		tokens = append(tokens, field+" LIKE ? ")
		args = append(args, k)
	}
	var whereStr = "(" + strings.Join(tokens, " OR ") + ")"
	b._whereArg(whereStr, args...)
	return b
}

// BetweenDateTimeOfDate
// 用 Date类型，去查询 DateTime 字段
func (b *SqlBuilder[T]) BetweenDateTimeOfDate(whereStr string, dateBegin, dateEnd *types.LocalDate, condition ...bool) *SqlBuilder[T] {
	b.selectStatus = selectDone
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	for _, c := range condition {
		if !c {
			return b
		}
	}

	var dateTimeBegin *types.LocalDateTime = nil
	if dateBegin != nil {
		dateTimeBegin = dateBegin.ToDateTimeP()
	}

	var dateTimeEnd *types.LocalDateTime = nil
	if dateEnd != nil {
		dateTimeEnd = dateEnd.Add(types.Duration().Day(1)).ToDateTimeP()
	}

	if dateTimeBegin != nil {
		b._whereArg(whereStr+" >= ?", dateTimeBegin)
	}
	if dateTimeEnd != nil {
		b._whereArg(whereStr+" < ?", dateTimeEnd)
	}

	return b
}
