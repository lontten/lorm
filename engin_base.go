package lorm

import "log"


type EngineBase struct {
	context OrmContext
	db      DBer
}


func (engine EngineBase) Select(args ...string) OrmSelect {
	context := engine.context
	selectArgsArr2SqlStr(context, args)
	return OrmSelect{db: engine.db, context: context}
}

func (orm OrmSelect) SelectModel(v interface{}) OrmSelect {
	if v == nil {
		return orm
	}
	context := orm.context

	return OrmSelect{db: orm.db, context: context}
}

func (orm *OrmSelect) From(arg string) *OrmFrom {
	context := orm.context
	context.query.WriteString(" FROM " + arg)
	return &OrmFrom{db: orm.db, context: context}
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
	return &OrmWhere{db: orm.db, context: orm.context}
}

func (orm *OrmWhere) GetOne(dest interface{}) (int64, error) {
	s := orm.context.query.String()
	rows, err := orm.db.Query(s, orm.context.args...)
	if err != nil {
		return 0, err
	}
	return StructScanLn(rows, dest)
}

func (orm *OrmWhere) GetList(dest interface{}) (num int64, err error) {
	query := orm.context.query.String()
	log.Println(query)
	args := orm.context.args
	log.Printf("args :: %s", args)
	rows, err := orm.db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	return StructScan(rows, dest)
}
