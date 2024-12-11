package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func QueryBuild(db Engine) *SqlBuilder {
	db.getDialect().initContext()
	return &SqlBuilder{
		db:              db,
		selectQuery:     &strings.Builder{},
		otherSqlBuilder: &strings.Builder{},
	}
}

type SqlBuilder struct {
	db Engine
	// 最终执行sql
	query string
	// 最终参数列表
	args []any

	// set 部分
	setTokens []string

	// select部分
	selectStatus int8

	selectTokens []string
	selectQuery  *strings.Builder
	selectArgs   []any

	// 其他部分
	otherSqlBuilder *strings.Builder
	otherSqlArgs    []any
	whereStatus     int8

	// page
	other any
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

func (b *SqlBuilder) initSelectSql() {
	b.selectQuery.WriteString("SELECT ")
	b.selectQuery.WriteString(strings.Join(b.selectTokens, ","))
	b.query = b.selectQuery.String() + " " + b.otherSqlBuilder.String()
	b.args = append(b.selectArgs, b.otherSqlArgs...)
}

//
//func (b *SqlBuilder) initQuerySql() {
//	b.query = b.otherSqlBuilder.String()
//	b.args = b.otherSqlArgs
//}

// 显示sql
func (b *SqlBuilder) ShowSql() *SqlBuilder {
	b.db.getCtx().showSql = true
	return b
}

// 添加一个 arg，多个断言
func (b *SqlBuilder) AppendArg(arg any, condition ...bool) *SqlBuilder {
	if b.db.getCtx().hasErr() {
		return b
	}
	for _, c := range condition {
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
func (b *SqlBuilder) AppendSql(sql string) *SqlBuilder {
	b.otherSqlBuilder.WriteString(sql)
	return b
}

// 添加 多个参数
func (b *SqlBuilder) AppendArgs(args ...any) *SqlBuilder {
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
func (b *SqlBuilder) Select(arg string, condition ...bool) *SqlBuilder {
	ctx := b.db.getCtx()
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
func (b *SqlBuilder) SelectModel(v any) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if v == nil {
		return b
	}

	ctx.initScanDestOne(v)
	columns, err := ctx.ormConf.getStructField(ctx.destBaseType)
	if err != nil {
		ctx.err = err
		return b
	}

	b.selectStatus = selectSet
	b.selectTokens = append(b.selectTokens, columns...)
	return b
}

// from 表名
// 状态从 selectNoSet 变成 selectSet
func (b *SqlBuilder) From(name string) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	b.selectStatus = selectDone
	b.otherSqlBuilder.WriteString(" FROM " + name)
	return b
}

// join 联表
func (b *SqlBuilder) Join(name string, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) Arg(arg any, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) Args(args ...any) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder) LeftJoin(name string, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) RightJoin(name string, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) OrderBy(name string, condition ...bool) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" ORDER BY " + name)
	return b
}

func (b *SqlBuilder) Native(sql string, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) OrderDescBy(name string, condition ...bool) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.otherSqlBuilder.WriteString(" ORDER BY " + name + " DESC")
	return b
}

func (b *SqlBuilder) Limit(num int64, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) Offset(num int64, condition ...bool) *SqlBuilder {
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

func (b *SqlBuilder) WhereBuilder(w *WhereBuilder) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if w == nil {
		return b
	}
	sqlStr, err := w.toSql(b.db.getDialect().parse)
	if err != nil {
		b.db.getCtx().err = err
		return b
	}
	if sqlStr == "" {
		return b
	}

	if b.selectStatus != selectDone {
		ctx.err = errors.New("未完成 select 设置")
	}
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

	b.AppendArgs(w.args...)
	return b
}

func (b *SqlBuilder) WhereIng() *SqlBuilder {
	b.selectStatus = selectDone
	b.whereStatus = whereSet
	return b
}
func (b *SqlBuilder) Where(whereStr string, condition ...bool) *SqlBuilder {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if b.selectStatus != selectDone {
		ctx.err = errors.New("Where 设置异常：" + whereStr)
		return b
	}

	for _, c := range condition {
		if !c {
			return b
		}
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

func (b *SqlBuilder) BoolWhere(condition bool, whereStr string, args ...any) *SqlBuilder {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if b.selectStatus != selectDone {
		ctx.err = errors.New("Where 设置异常：" + whereStr)
		return b
	}

	if !condition {
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

func (b *SqlBuilder) WhereIn(whereStr string, args ...any) *SqlBuilder {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	if args == nil {
		return b
	}
	length := len(args)
	if length == 0 {
		return b
	}

	if b.selectStatus != selectDone {
		ctx.err = errors.New("Where 设置异常：" + whereStr)
		return b
	}

	b.AppendArgs(args...)

	var inArgStr = " (" + gen(length) + ")"
	whereStr = whereStr + " IN" + inArgStr

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

// WhereSqlIn
// in ? 当参数列表长度为0时，跳过这个where
func (b *SqlBuilder) WhereSqlIn(whereStr string, args ...any) *SqlBuilder {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return b
	}

	if args == nil {
		return b
	}
	length := len(args)
	if length == 0 {
		return b
	}

	if b.selectStatus != selectDone {
		ctx.err = errors.New("Where 设置异常：" + whereStr)
		return b
	}

	b.AppendArgs(args...)

	var inArgStr = " (" + gen(length) + ")"
	whereStr = strings.Replace(whereStr, "?", inArgStr, -1)

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

func (b *SqlBuilder) ScanOne(dest any) (rowsNum int64, err error) {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return 0, ctx.err
	}
	b.selectStatus = selectDone
	b.whereStatus = whereDone

	ctx.initScanDestOne(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}

	b.initSelectSql()
	query := b.query
	args := b.args
	if ctx.showSql {
		fmt.Println(query, args)
	}

	rows, err := db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLnT(rows)
}

func (b *SqlBuilder) ScanList(dest any) (rowsNum int64, err error) {
	db := b.db
	ctx := db.getCtx()
	if ctx.hasErr() {
		return 0, ctx.err
	}
	b.selectStatus = selectDone
	b.whereStatus = whereDone

	ctx.initScanDestList(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}
	b.initSelectSql()

	query := b.query
	args := b.args

	if ctx.showSql {
		fmt.Println(query, args)
	}
	rows, err := db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanT(rows)
}

func (b *SqlBuilder) Exec() (sql.Result, error) {
	db := b.db
	ctx := db.getCtx()
	b.selectStatus = selectDone
	b.whereStatus = whereDone
	if ctx.hasErr() {
		return nil, ctx.err
	}
	b.initSelectSql()
	if ctx.showSql {
		fmt.Println(b.query, b.args)
	}
	return db.exec(b.query, b.args...)
}
