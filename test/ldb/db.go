package ldb

import "github.com/lontten/lorm"

//import _ "github.com/go-sql-driver/mysql"

import _ "github.com/jackc/pgx/v5/stdlib"

var DB lorm.Engine

func init2() {
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

func init() {
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
