package lorm

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func TestWhereBuilder_toSql(t *testing.T) {
	conf := MysqlConf{
		Host:     "127.0.0.1",
		Port:     "3306",
		DbName:   "test",
		User:     "root",
		Password: "123456",
	}
	db, err := Connect(conf, nil)
	if err != nil {
		panic(err)
	}

	w1 := W().Eq("a", 1)
	w2 := W().Eq("x", 3)

	w11 := W().Or(w1).And(w1).Or(w2).And(w1)
	w22 := W().Or(w1)
	w22 = w22.Or(w2)

	builder := W().
		And(w11).
		And(w22)
	sql, err := builder.toSql(db.getDialect().parse)
	fmt.Println(err)
	fmt.Println(sql)

}
