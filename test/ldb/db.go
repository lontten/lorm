package ldb

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/lontten/lorm"
)

//import _ "github.com/jackc/pgx/v5/stdlib"

var DB lorm.Engine

func init() {
	conf := lorm.MysqlConf{
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "test",
		User:     "root",
		Password: "123456",
	}
	db, err := lorm.Connect(conf, nil)
	if err != nil {
		panic(err)
	}
	DB = db
}

func init2() {
	conf := lorm.PgConf{
		Host:     "127.0.0.1",
		Port:     "5432",
		DbName:   "test",
		User:     "postgres",
		Password: "123456",
	}
	db, err := lorm.Connect(conf, nil)
	if err != nil {
		panic(err)
	}
	DB = db
}
