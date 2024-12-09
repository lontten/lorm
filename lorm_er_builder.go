package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func SelectBuilder(db Engine) *SqlBuilder {
	return &SqlBuilder{
		db:          db,
		selectQuery: &strings.Builder{},
		otherQuery:  &strings.Builder{},
	}
}

type SqlBuilder struct {
	db Engine

	query string
	args  []any

	// select部分
	selectTokens []string
	selectQuery  *strings.Builder
	selectArgs   []any
	selectStatus int8

	// 其他部分
	otherQuery  *strings.Builder
	otherArgs   []any
	whereStatus int8

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
	b.query = b.selectQuery.String() + " " + b.otherQuery.String()
	b.args = append(b.selectArgs, b.otherArgs...)
}

//
//func (b *SqlBuilder) initQuerySql() {
//	b.query = b.otherQuery.String()
//	b.args = b.otherArgs
//}

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
		b.otherArgs = append(b.otherArgs, arg)
	}
	return b
}

// 添加sql语句
func (b *SqlBuilder) AppendSql(sql string) *SqlBuilder {
	b.otherQuery.WriteString(sql)
	return b
}

// 添加 多个参数
func (b *SqlBuilder) AppendArgs(args ...any) *SqlBuilder {
	ctx := b.db.getCtx()
	if ctx.hasErr() {
		return b
	}
	if b.selectStatus == selectNoSet {
		b.selectArgs = append(b.selectArgs, args...)
	} else {
		b.otherArgs = append(b.otherArgs, args...)
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
	b.otherQuery.WriteString(" FROM " + name)
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
	b.otherQuery.WriteString(" JOIN " + name)
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
	b.otherQuery.WriteString("\n")
	b.otherQuery.WriteString("LEFT JOIN " + name)
	b.otherQuery.WriteString("\n")

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
	b.otherQuery.WriteString(" RIGHT JOIN " + name)
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
	b.otherQuery.WriteString(" ORDER BY " + name)
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
	b.otherQuery.WriteString(" ")
	b.otherQuery.WriteString(sql)
	b.otherQuery.WriteString(" ")
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
	b.otherQuery.WriteString(" ORDER BY " + name + " DESC")
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
	b.otherQuery.WriteString(" LIMIT " + strconv.FormatInt(num, 10))
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
	b.otherQuery.WriteString(" OFFSET " + strconv.FormatInt(num, 10))
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
		b.otherQuery.WriteString(" WHERE ")
		b.otherQuery.WriteString(sqlStr)
	case whereSet:
		b.otherQuery.WriteString(" AND ")
		b.otherQuery.WriteString(sqlStr)
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
		b.otherQuery.WriteString(" WHERE ")
		b.otherQuery.WriteString(whereStr)
	case whereSet:
		b.otherQuery.WriteString(" AND ")
		b.otherQuery.WriteString(whereStr)
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

	b.args = append(b.args, args...)

	var inArgStr = " (" + gen(length) + ")"
	whereStr = strings.Replace(whereStr, "?", inArgStr, -1)

	switch b.whereStatus {
	case whereNoSet:
		b.whereStatus = whereSet
		b.otherQuery.WriteString(" WHERE ")

		b.otherQuery.WriteString(whereStr)

	case whereSet:
		b.otherQuery.WriteString(" AND ")

		b.otherQuery.WriteString(whereStr)

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
	fmt.Println(query, args)

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
	return db.exec(b.query, b.args...)
}
