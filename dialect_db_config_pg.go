package lorm

import (
	"database/sql"
	"log"
	"os"
)

type PgConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Other    string
}

func (c *PgConf) dialect(ctx *ormContext, pc *PoolConf) Dialecter {
	var logger *log.Logger
	if pc == nil || pc.Logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
		log.SetFlags(log.LstdFlags | log.Llongfile)
	} else {
		logger = pc.Logger
	}
	return &PgDialect{ctx: ctx, log: Logger{log: logger}}
}

func (c *PgConf) Open() (*sql.DB, error) {
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
