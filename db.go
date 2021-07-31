package lorm

import (
	"database/sql"
)

type DB struct {
	db       *sql.DB
	dbConfig DbConfig
}

func (db DB) Db(c *OrmConf) Engine {
	conf := OrmConf{}
	if c != nil {
		conf = *c
	}
	return Engine{
		db:       db,
		lormConf: conf,
		Base: EngineBase{
			core:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(db.db),
		},
		Extra: EngineExtra{
			core:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(db.db),
		},
		Classic: EngineNative{
			core:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(db.db),
		},
		Table: EngineTable{
			core:    conf,
			ctx:     OrmContext{},
			dialect: db.dbConfig.Dialect(db.db),
		},
	}
}

type OrmSelect struct {
	base EngineBase
}

type OrmFrom struct {
	base EngineBase
}

type OrmWhere struct {
	base EngineBase
}
