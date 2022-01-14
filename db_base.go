package lorm

import "log"

//对基础SQL的简单封装
type EngineBase struct {
	dialect Dialect
	context OrmContext
}

func (engine EngineBase) Select(args ...string) OrmSelect {
	engine.context.selectArgsArr2SqlStr(args)
	return OrmSelect{engine}
}

func (orm OrmSelect) SelectModel(v interface{}) OrmSelect {
	if v == nil {
		return orm
	}
	return OrmSelect{orm.base}
}

func (orm *OrmSelect) From(arg string) *OrmFrom {
	base := orm.base
	base.context.query.WriteString(" FROM " + arg)
	return &OrmFrom{base}
}

func (orm *OrmFrom) Where(v *WhereBuilder) *OrmWhere {
	base := orm.base
	if v == nil {
		return &OrmWhere{base}
	}

	wheres := v.context.wheres
	for i, where := range wheres {
		if i == 0 {
			base.context.query.WriteString(" WHERE " + where)
			continue
		}
		base.context.query.WriteString(" AND " + where)
	}
	base.context.args = v.context.args
	return &OrmWhere{base}
}

//v0.7
func (orm *OrmWhere) GetOne(dest interface{}) (int64, error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}
	orm.base.context.initScanDestOne(dest)
	orm.base.context.checkScanDestField()
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	base := orm.base
	ctx := orm.base.context
	s := ctx.query.String()
	rows, err := base.dialect.query(s, ctx.args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLn(rows)
}

func (orm *OrmWhere) GetList(dest interface{}) (num int64, err error) {
	if err := orm.base.context.err; err != nil {
		return 0, err
	}
	orm.base.context.initScanDestList(dest)
	orm.base.context.checkScanDestField()
	if err := orm.base.context.err; err != nil {
		return 0, err
	}

	base := orm.base
	ctx := orm.base.context

	query := ctx.query.String()
	log.Println(query)
	args := ctx.args
	log.Printf("args :: %s", args)
	rows, err := base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return ctx.Scan(rows)
}
