package lorm

import (
	"database/sql"
)

type MysqlConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

func (c MysqlConf) dialect(ctx *ormContext, db DBer) Dialecter {
	return &MysqlDialect{
		ctx: ctx,
		db:  db,
	}
}

func (c MysqlConf) open() (*sql.DB, error) {
	dsn := c.User + ":" + c.Password +
		"@tcp(" + c.Host +
		":" + c.Port +
		")/" + c.DbName
	return sql.Open("mysql", dsn)
}
