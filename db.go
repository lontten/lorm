package lorm

import (
	"database/sql"
)

type DB struct {
	db       *sql.DB
	dbConfig DbConfig
}

func (db DB) Db(c *OrmConf) Engine {
	if c != nil {
		ormConfig = *c
	}
	return Engine{
		db:      db,
		ctx:     OrmContext{},
		dialect: db.dbConfig.Dialect(db.db),
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
