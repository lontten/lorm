package lorm

import (
	"database/sql"
	"log"
	"os"
)

type MysqlConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

func (c *MysqlConf) dialect(ctx *ormContext, pc *PoolConf) Dialecter {
	var logger *log.Logger
	if pc == nil || pc.Logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
		log.SetFlags(log.LstdFlags | log.Llongfile)
	} else {
		logger = pc.Logger
	}
	return &PgDialect{ctx: ctx, log: Logger{log: logger}}
}

func (c *MysqlConf) Open() (*sql.DB, error) {
	dsn := c.User + ":" + c.Password +
		"@tcp(" + c.Host +
		":" + c.Port +
		")/" + c.DbName
	return sql.Open("mysql", dsn)
}
