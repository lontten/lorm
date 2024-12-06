package lorm

import (
	"database/sql"
	"errors"
	"strings"
)

// todo 下面未重构--------------
func (db *lnDB) Builder() *SqlBuilder {
	return &SqlBuilder{
		core:        db.core,
		selectQuery: &strings.Builder{},
		otherQuery:  &strings.Builder{},
	}
}

type SqlBuilder struct {
	core corer

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
	selectNone = iota
	selectIng
	selectDone
)

const (
	whereNone = iota
	whereIng
	whereDone
)

func (b *SqlBuilder) updStatus() {
	if b.selectStatus == selectNone {
		b.selectStatus = selectIng
		return
	}
	if b.selectStatus == selectIng {
		b.selectStatus = selectDone
		return
	}

	if b.whereStatus == whereIng {
		b.whereStatus = whereDone
		return
	}
}

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

func (b *SqlBuilder) AppendArg(arg any, condition ...bool) *SqlBuilder {
	if b.core.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	if b.selectStatus == selectIng {
		b.selectArgs = append(b.selectArgs, arg)
	} else {
		b.otherArgs = append(b.otherArgs, arg)
	}
	return b
}

func (b *SqlBuilder) AppendSql(sql string) *SqlBuilder {
	b.otherQuery.WriteString(sql)
	return b
}

func (b *SqlBuilder) AppendArgs(args ...any) *SqlBuilder {
	if b.core.hasErr() {
		return b
	}
	if b.selectStatus == selectIng {
		b.selectArgs = append(b.selectArgs, args...)
	} else {
		b.otherArgs = append(b.otherArgs, args...)
	}
	return b
}

func (b *SqlBuilder) Select(arg string, condition ...bool) *SqlBuilder {
	if b.core.hasErr() {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}

	b.updStatus()
	b.selectTokens = append(b.selectTokens, arg)

	return b
}

func (b *SqlBuilder) SelectModel(v any) *SqlBuilder {
	if b.core.hasErr() {
		return b
	}
	if v == nil {
		return b
	}

	b.core.getCtx().initScanDestOne(v)
	columns, err := b.core.getCtx().ormConf.getStructField(b.core.getCtx().destBaseType)
	if err != nil {
		b.core.getCtx().err = err
		return b
	}

	b.updStatus()
	b.selectTokens = append(b.selectTokens, columns...)
	return b
}

func (b *SqlBuilder) From(name string) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	b.updStatus()
	b.otherQuery.WriteString(" FROM " + name)
	return b
}

func (b *SqlBuilder) Join(name string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" JOIN " + name)
	return b
}

func (b *SqlBuilder) Arg(arg any, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
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
	if b.core.getCtx().err != nil {
		return b
	}
	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder) LeftJoin(name string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString("\n")
	b.otherQuery.WriteString("LEFT JOIN " + name)
	b.otherQuery.WriteString("\n")

	return b
}

func (b *SqlBuilder) RightJoin(name string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" RIGHT JOIN " + name)
	return b
}

func (b *SqlBuilder) OrderBy(name string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" ORDER BY " + name)
	return b
}

func (b *SqlBuilder) Native(sql string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" ")
	b.otherQuery.WriteString(sql)
	b.otherQuery.WriteString(" ")
	return b
}

func (b *SqlBuilder) OrderDescBy(name string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" ORDER BY " + name + " DESC")
	return b
}

func (b *SqlBuilder) Limit(num int64, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" LIMIT ? ")
	b.AppendArg(num)
	return b
}

func (b *SqlBuilder) Offset(num int64, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.updStatus()
	b.otherQuery.WriteString(" OFFSET ? ")
	b.AppendArg(num)
	return b
}

func (b *SqlBuilder) WhereBuilder(w *WhereBuilder) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	if w == nil {
		return b
	}
	//sql, err := w.toSql(b.core.dialect.parse)
	//if err != nil {
	//	b.core.getCtx().err = err
	//	return b
	//}
	//if sql == "" {
	//	return b
	//}

	b.updStatus()

	switch b.whereStatus {
	case whereNone:
		b.whereStatus = whereIng
		b.otherQuery.WriteString(" WHERE ")
		//b.otherQuery.WriteString(sql)
	case whereIng:
		b.otherQuery.WriteString(" AND ")
		//b.otherQuery.WriteString(sql)
	case whereDone:
		b.core.getCtx().err = errors.New("where has been done")
	}

	b.AppendArgs(w.args...)
	return b
}

func (b *SqlBuilder) Where(whereStr string, condition ...bool) *SqlBuilder {
	if b.core.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}

	b.updStatus()

	switch b.whereStatus {
	case whereNone:
		b.whereStatus = whereIng
		b.otherQuery.WriteString(" WHERE ")
		b.otherQuery.WriteString(whereStr)
	case whereIng:
		b.otherQuery.WriteString(" AND ")
		b.otherQuery.WriteString(whereStr)
	case whereDone:
		b.core.getCtx().err = errors.New("where has been done")
	}

	return b
}

func (b *SqlBuilder) ScanOne(dest any) (rowsNum int64, err error) {
	if err = b.core.getCtx().err; err != nil {
		return 0, err
	}
	b.updStatus()
	b.core.getCtx().initScanDestOne(dest)
	if b.core.getCtx().destIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	b.core.getCtx().checkScanDestField()
	if err = b.core.getCtx().err; err != nil {
		return 0, err
	}
	b.initSelectSql()
	rows, err := b.core.doQuery(b.query, b.args...)

	if err != nil {
		return 0, err
	}
	return b.core.getCtx().ScanLnT(rows)
}

func (b *SqlBuilder) ScanList(dest any) (rowsNum int64, err error) {
	if err = b.core.getCtx().err; err != nil {
		return 0, err
	}
	b.updStatus()
	b.core.getCtx().initScanDestList(dest)
	b.core.getCtx().checkScanDestField()

	if err = b.core.getCtx().err; err != nil {
		return 0, err
	}
	b.initSelectSql()
	rows, err := b.core.doQuery(b.query, b.args...)

	if err != nil {
		return 0, err
	}
	return b.core.getCtx().ScanT(rows)
}

func (b *SqlBuilder) Exec() (sql.Result, error) {
	b.updStatus()
	if err := b.core.getCtx().err; err != nil {
		return nil, err
	}
	b.initSelectSql()
	return b.core.doExec(b.query, b.args...)
}
