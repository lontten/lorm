package dbinit

import (
	"os"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/lontten/lcore/v2/logutil"
	"github.com/lontten/lorm"
)

var DB lorm.Engine

func init() {
	initMysql()
}
func initMysql() {
	conf := lorm.MysqlConf{
		Host:     os.Getenv("LORM_MYSQL_HOST"),
		Port:     os.Getenv("LORM_MYSQL_PORT"),
		DbName:   os.Getenv("LORM_MYSQL_DB"),
		User:     os.Getenv("LORM_MYSQL_USER"),
		Password: os.Getenv("LORM_MYSQL_PWD"),
		Version:  lorm.MysqlVersion5,
	}
	logutil.Log(conf)
	db, err := lorm.Connect(conf, nil)
	if err != nil {
		panic(err)
	}
	DB = db
}

func initPg() {
	conf := lorm.PgConf{
		Host:     os.Getenv("LORM_PG_HOST"),
		Port:     os.Getenv("LORM_PG_PORT"),
		DbName:   os.Getenv("LORM_PG_DB"),
		User:     os.Getenv("LORM_PG_USER"),
		Password: os.Getenv("LORM_PG_PWD"),
	}
	db, err := lorm.Connect(conf, nil)
	if err != nil {
		panic(err)
	}
	DB = db
}
