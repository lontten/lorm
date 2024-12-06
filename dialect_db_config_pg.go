package lorm

import (
	"database/sql"
)

type PgConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Other    string
}

func (c PgConf) dialect(ctx *ormContext) Dialecter {
	ctx.dialectNeedLastInsertId = false
	return &PgDialect{ctx: ctx}
}

func (c PgConf) open() (*sql.DB, error) {
	dsn := "user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DbName +
		" host=" + c.Host +
		" port= " + c.Port +
		" "
	if c.Other == "" {
		dsn += "sslmode=disable TimeZone=Asia/Shanghai"
	}
	dsn += c.Other
	return sql.Open("pgx", dsn)
}
