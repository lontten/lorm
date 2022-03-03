package lorm

import (
	"errors"
	"strings"
)

func (db DB) Builder() *SqlBuilder {
	return &SqlBuilder{
		db:          db,
		selectQuery: &strings.Builder{},
		otherQuery:  &strings.Builder{},
	}
}

type SqlBuilder struct {
	db    DB
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

func (b *SqlBuilder) initQuerySql() {
	b.query = b.otherQuery.String()
	b.args = b.otherArgs
}

func (b *SqlBuilder) AppendArg(arg interface{}, condition ...bool) *SqlBuilder {
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
		return b
	}
	if v == nil {
		return b
	}

	b.db.ctx.initScanDestOne(v)
	columns, err := b.db.ctx.conf.initColumns(b.db.ctx.scanDestBaseType)
	if err != nil {
		b.db.ctx.err = err
		return b
	}

	b.updStatus()
	b.selectTokens = append(b.selectTokens, columns...)
	return b
}

func (b *SqlBuilder) From(name string) *SqlBuilder {
	if b.db.ctx.err != nil {
		return b
	}
	b.updStatus()
	b.otherQuery.WriteString(" FROM " + name)
	return b
}

func (b *SqlBuilder) Join(name string, condition ...bool) *SqlBuilder {
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
		return b
	}
	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder) LeftJoin(name string, condition ...bool) *SqlBuilder {
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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
	if b.db.ctx.err != nil {
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

func (b *SqlBuilder) WhereBuilder(v *WhereBuilder) *SqlBuilder {
	if b.db.ctx.err != nil {
		return b
	}
	if v == nil {
		return b
	}
	wheres := v.context.wheres
	if len(wheres) == 0 {
		return b
	}

	b.updStatus()
	whereStr := strings.Join(wheres, " AND ")

	switch b.whereStatus {
	case whereNone:
		b.whereStatus = whereIng
		b.otherQuery.WriteString(" WHERE ")
		b.otherQuery.WriteString(whereStr)
	case whereIng:
		b.otherQuery.WriteString(" AND ")
		b.otherQuery.WriteString(whereStr)
	case whereDone:
		b.db.ctx.err = errors.New("where has been done")
	}

	b.AppendArgs(v.context.args...)
	return b
}

func (b *SqlBuilder) Where(whereStr string, condition ...bool) *SqlBuilder {
	if b.db.ctx.err != nil {
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
		b.db.ctx.err = errors.New("where has been done")
	}

	return b
}

func (b *SqlBuilder) ScanOne(dest interface{}) (rowsNum int64, err error) {
	if err = b.db.ctx.err; err != nil {
		return 0, err
	}
	b.updStatus()
	b.db.ctx.initScanDestOne(dest)
	if b.db.ctx.scanIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	b.db.ctx.checkScanDestField()
	if err = b.db.ctx.err; err != nil {
		return 0, err
	}
	b.initQuerySql()
	rows, err := b.db.dialect.query(b.query, b.args...)
	if err != nil {
		return 0, err
	}
	return b.db.ctx.ScanLn(rows)
}

func (b *SqlBuilder) ScanList(dest interface{}) (rowsNum int64, err error) {
	if err = b.db.ctx.err; err != nil {
		return 0, err
	}
	b.updStatus()
	b.db.ctx.initScanDestList(dest)
	b.db.ctx.checkScanDestField()

	if err = b.db.ctx.err; err != nil {
		return 0, err
	}
	b.initQuerySql()
	rows, err := b.db.dialect.query(b.query, b.args...)
	if err != nil {
		return 0, err
	}
	return b.db.ctx.Scan(rows)
}

func (b *SqlBuilder) Exec() (rowsNum int64, err error) {
	b.updStatus()
	if err = b.db.ctx.err; err != nil {
		return 0, err
	}
	b.initQuerySql()
	return b.db.dialect.exec(b.query, b.args...)
}
