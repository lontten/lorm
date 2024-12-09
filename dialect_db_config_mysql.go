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
	Other    string
}

func (c MysqlConf) dialect(ctx *ormContext) Dialecter {
	ctx.dialectNeedLastInsertId = true
	return &MysqlDialect{ctx: ctx}
}

func (c MysqlConf) open() (*sql.DB, error) {
	dsn := c.User + ":" + c.Password +
		"@tcp(" + c.Host +
		":" + c.Port +
		")/" + c.DbName + "?"

	if c.Other == "" {
		dsn += "charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		dsn += c.Other
	}
	return sql.Open("mysql", dsn)
}
