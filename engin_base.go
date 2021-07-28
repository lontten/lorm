package lorm

import "log"

//对基础SQL的简单封装
type EngineBase struct {
	db       DBer
	lormConf Lorm
	dialect  Dialect

	context OrmContext
}

func (engine EngineBase) Select(args ...string) OrmSelect {
	context := engine.context
	selectArgsArr2SqlStr(context, args)
	return OrmSelect{
		db:       engine.db,
		lormConf: engine.lormConf,
		dialect:  engine.dialect,
		context:  context,
	}
}

func (orm OrmSelect) SelectModel(v interface{}) OrmSelect {
	if v == nil {
		return orm
	}

	context := orm.context

	return OrmSelect{
		db:       orm.db,
		lormConf: orm.lormConf,
		dialect:  orm.dialect,
		context:  context,
	}
}

func (orm *OrmSelect) From(arg string) *OrmFrom {
	context := orm.context
	context.query.WriteString(" FROM " + arg)
	return &OrmFrom{
		db:       orm.db,
		lormConf: orm.lormConf,
		dialect:  orm.dialect,
		context:  context,
	}
}

func (orm *OrmFrom) Where(v *WhereBuilder) *OrmWhere {
	if v == nil {
		return &OrmWhere{db: orm.db, context: orm.context}
	}
	wheres := v.context.wheres
	for i, where := range wheres {
		if i == 0 {
			orm.context.query.WriteString(" WHERE " + where)
			continue
		}
		orm.context.query.WriteString(" AND " + where)
	}
	orm.context.args = v.context.args
	return &OrmWhere{
		db:       orm.db,
		lormConf: orm.lormConf,
		dialect:  orm.dialect,
		context:  orm.context,
	}
}

func (orm *OrmWhere) GetOne(dest interface{}) (int64, error) {
	s := orm.context.query.String()
	rows, err := orm.db.query(s, orm.context.args...)
	if err != nil {
		return 0, err
	}
	return orm.lormConf.ScanLn(rows, dest)
}

func (orm *OrmWhere) GetList(dest interface{}) (num int64, err error) {
	query := orm.context.query.String()
	log.Println(query)
	args := orm.context.args
	log.Printf("args :: %s", args)
	rows, err := orm.db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return orm.lormConf.Scan(rows, dest)
}
