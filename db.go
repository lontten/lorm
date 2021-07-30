package lorm

import (
	"database/sql"
	"strings"
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
			dialect: db.dbConfig.Dialect(conf),
		},
		Extra: EngineExtra{
			core:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
		Classic: EngineNative{
			core:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
		Table: EngineTable{
			core:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
	}
}

type OrmContext struct {
	query  *strings.Builder
	args   []interface{}
	startd bool
	err    error
	log    int
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
