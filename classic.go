package lorm

import (
	"errors"
	"strings"
)

// todo 下面未重构--------------
func (db lnDB) Builder() *SqlBuilder {
	return &SqlBuilder{
		db:          db.core,
		selectQuery: &strings.Builder{},
		otherQuery:  &strings.Builder{},
	}
}

type SqlBuilder struct {
	db corer

	query string
	args  []interface{}

	// select部分
	selectTokens []string
	selectQuery  *strings.Builder
	selectArgs   []interface{}
	selectStatus int8

	// 其他部分
	otherQuery  *strings.Builder
	otherArgs   []interface{}
	whereStatus int8

	// page
	other interface{}
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
	}
	if b.selectStatus == selectIng {
		b.selectStatus = selectDone
	}

	if b.whereStatus == whereIng {
		b.whereStatus = whereDone
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

func (b *SqlBuilder) AppendArg(arg interface{}, condition ...bool) *SqlBuilder {
	if b.db.getCtx().err != nil {
		return b
	}
	for _, c := range condition {
		if !c {
			return b
		}
	}
	if b.selectStatus == selectIng {
		b.selectArgs = append(b.selectArgs, arg)
		return b
	}
	b.otherArgs = append(b.otherArgs, arg)
	return b
}

func (b *SqlBuilder) AppendSql(sql string) *SqlBuilder {
	b.otherQuery.WriteString(sql)
	return b
}

func (b *SqlBuilder) AppendArgs(args ...interface{}) *SqlBuilder {
	if b.db.getCtx().err != nil {
		return b
	}
	if b.selectStatus == selectIng {
		b.selectArgs = append(b.selectArgs, args...)
		return b
	}
	b.otherArgs = append(b.otherArgs, args...)
	return b
}

func (b *SqlBuilder) Select(arg string, condition ...bool) *SqlBuilder {
	if b.db.getCtx().err != nil {
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

func (b *SqlBuilder) SelectModel(v interface{}) *SqlBuilder {
	if b.db.getCtx().err != nil {
		return b
	}
	if v == nil {
		return b
	}

	b.db.getCtx().initScanDestOne(v)
	columns, err := b.db.getCtx().ormConf.initColumns(b.db.getCtx().scanDestBaseType)
	if err != nil {
		b.db.getCtx().err = err
		return b
	}

	b.updStatus()
	b.selectTokens = append(b.selectTokens, columns...)
	return b
}

func (b *SqlBuilder) From(name string) *SqlBuilder {
	if b.db.getCtx().err != nil {
		return b
	}
	b.updStatus()
	b.otherQuery.WriteString(" FROM " + name)
	return b
}

func (b *SqlBuilder) Join(name string, condition ...bool) *SqlBuilder {
	if b.db.getCtx().err != nil {
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

func (b *SqlBuilder) Arg(arg interface{}, condition ...bool) *SqlBuilder {
	if b.db.getCtx().err != nil {
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

func (b *SqlBuilder) Args(args ...interface{}) *SqlBuilder {
	if b.db.getCtx().err != nil {
		return b
	}
	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder) LeftJoin(name string, condition ...bool) *SqlBuilder {
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
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
	if b.db.getCtx().err != nil {
		return b
	}
	if w == nil {
		return b
	}
	sql, err := w.toSql(b.db.dialect.parse)
	if err != nil {
		b.db.getCtx().err = err
		return b
	}

	if sql == "" {
		return b
	}

	b.updStatus()

	switch b.whereStatus {
	case whereNone:
		b.whereStatus = whereIng
		b.otherQuery.WriteString(" WHERE ")
		b.otherQuery.WriteString(sql)
	case whereIng:
		b.otherQuery.WriteString(" AND ")
		b.otherQuery.WriteString(sql)
	case whereDone:
		b.db.getCtx().err = errors.New("where has been done")
	}

	b.AppendArgs(w.args...)
	return b
}

func (b *SqlBuilder) Where(whereStr string, condition ...bool) *SqlBuilder {
	if b.db.getCtx().err != nil {
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
		b.db.getCtx().err = errors.New("where has been done")
	}

	return b
}

func (b *SqlBuilder) ScanOne(dest interface{}) (rowsNum int64, err error) {
	if err = b.db.getCtx().err; err != nil {
		return 0, err
	}
	b.updStatus()
	b.db.getCtx().initScanDestOne(dest)
	if b.db.getCtx().scanIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	b.db.getCtx().checkScanDestField()
	if err = b.db.getCtx().err; err != nil {
		return 0, err
	}
	b.initSelectSql()
	rows, err := b.db.doQuery(b.query, b.args...)

	if err != nil {
		return 0, err
	}
	return b.db.getCtx().ScanLn(rows)
}

func (b *SqlBuilder) ScanList(dest interface{}) (rowsNum int64, err error) {
	if err = b.db.getCtx().err; err != nil {
		return 0, err
	}
	b.updStatus()
	b.db.getCtx().initScanDestList(dest)
	b.db.getCtx().checkScanDestField()

	if err = b.db.getCtx().err; err != nil {
		return 0, err
	}
	b.initSelectSql()
	rows, err := b.db.doQuery(b.query, b.args...)

	if err != nil {
		return 0, err
	}
	return b.db.getCtx().Scan(rows)
}

func (b *SqlBuilder) Exec() (rowsNum int64, err error) {
	b.updStatus()
	if err = b.db.getCtx().err; err != nil {
		return 0, err
	}
	b.initSelectSql()
	query, args := b.db.dialect.exec(b.query, b.args...)
	return b.db.doExec(query, args...)
}
