package lorm

import (
	"errors"
	"strings"
)

type SqlBuilder struct {
	base DB

	queryTokens []string
	query       *strings.Builder
	args        []interface{}
	start       bool
}

func (b *SqlBuilder) AppendArg(arg interface{}) *SqlBuilder {
	b.args = append(b.args, arg)
	return b
}

func (b *SqlBuilder) initSelectSql() *SqlBuilder {
	if b.start {
		b.query.WriteString(strings.Join(b.queryTokens, ","))
	}
	b.start = false
	return b
}

func (b *SqlBuilder) AppendArgs(args ...interface{}) *SqlBuilder {
	b.args = append(b.args, args...)
	return b
}

func (db DB) Builder() *SqlBuilder {
	return &SqlBuilder{
		base:  db,
		query: &strings.Builder{},
	}
}

func (b *SqlBuilder) Select(arg string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}

	if !b.start {
		b.query.WriteString("SELECT ")
		b.start = true
	}
	b.queryTokens = append(b.queryTokens, arg)
	return b
}

func (b *SqlBuilder) SelectModel(v interface{}) *SqlBuilder {
	if v == nil {
		return b
	}

	if !b.start {
		b.query.WriteString("SELECT ")
		b.start = true
	}

	b.base.ctx.initScanDestOne(v)
	columns, err := b.base.ctx.conf.initColumns(b.base.ctx.scanDestBaseType)
	if err != nil {
		b.base.ctx.err = err
		return b
	}
	b.queryTokens = append(b.queryTokens, columns...)
	return b
}

func (b *SqlBuilder) From(name string) *SqlBuilder {
	b.initSelectSql()
	b.query.WriteString(" FROM " + name)
	return b
}

func (b *SqlBuilder) Join(name string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" JOIN " + name)
	return b
}

func (b *SqlBuilder) Arg(arg interface{}, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.AppendArgs(arg)
	return b
}

func (b *SqlBuilder) Args(args ...interface{}) *SqlBuilder {
	b.AppendArgs(args...)
	return b
}

func (b *SqlBuilder) LeftJoin(name string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString("\n")
	b.query.WriteString("LEFT JOIN " + name)
	b.query.WriteString("\n")

	return b
}

func (b *SqlBuilder) RightJoin(name string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" RIGHT JOIN " + name)
	return b
}

func (b *SqlBuilder) OrderBy(name string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" ORDER BY " + name)
	return b
}

func (b *SqlBuilder) Native(sql string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" ")
	b.query.WriteString(sql)
	b.query.WriteString(" ")
	return b
}

func (b *SqlBuilder) OrderDescBy(name string, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" ORDER BY " + name + " DESC")
	return b
}

func (b *SqlBuilder) Limit(num int64, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" LIMIT ? ")
	b.AppendArg(num)
	return b
}

func (b *SqlBuilder) Offset(num int64, condition ...bool) *SqlBuilder {
	for _, c := range condition {
		if !c {
			return b
		}
	}
	b.initSelectSql()
	b.query.WriteString(" OFFSET ? ")
	b.AppendArg(num)
	return b
}

func (b *SqlBuilder) Where(v *WhereBuilder) *SqlBuilder {
	if v == nil {
		return b
	}

	wheres := v.context.wheres
	for i, where := range wheres {
		if i == 0 {
			b.query.WriteString(" WHERE " + where)
			continue
		}
		b.query.WriteString(" AND " + where)
	}
	b.initSelectSql()
	b.AppendArgs(v.context.args)
	return b
}

func (b *SqlBuilder) ScanOne(dest interface{}) (rowsNum int64, err error) {
	b.initSelectSql()
	b.base.ctx.initScanDestOne(dest)
	if b.base.ctx.scanIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	b.base.ctx.checkScanDestField()
	if err = b.base.ctx.err; err != nil {
		return 0, err
	}

	rows, err := b.base.dialect.query(b.query.String(), b.args...)
	if err != nil {
		return 0, err
	}
	return b.base.ctx.ScanLn(rows)
}

func (b *SqlBuilder) ScanList(dest interface{}) (rowsNum int64, err error) {
	if err = b.base.ctx.err; err != nil {
		return 0, err
	}
	b.initSelectSql()
	b.base.ctx.initScanDestList(dest)
	b.base.ctx.checkScanDestField()

	if err = b.base.ctx.err; err != nil {
		return 0, err
	}

	rows, err := b.base.dialect.query(b.query.String(), b.args...)
	if err != nil {
		return 0, err
	}
	return b.base.ctx.Scan(rows)
}

func (b *SqlBuilder) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	b.initSelectSql()
	return b.base.dialect.exec(query, args...)
}
