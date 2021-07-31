package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Tx struct {
	db        *sql.Tx
	dbConfig  DbConfig
	ormConfig OrmConf
}

func (tx Tx) exec(query string, args ...interface{}) (int64, error) {
	switch tx.dbConfig.DriverName() {
	case MYSQL:
	case POSTGRES:
		var i = 1
		for {
			t := strings.Replace(query, " ? ", " $"+strconv.Itoa(i)+" ", 1)
			if t == query {
				break
			}
			i++
			query = t
		}
	default:
		return 0, errors.New("无此db drive 类型")
	}
	log.Println(query, args)

	exec, err := tx.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (tx Tx) query(query string, args ...interface{}) (*sql.Rows, error) {
	switch tx.dbConfig.DriverName() {
	case POSTGRES:
		var i = 1
		for {
			t := strings.Replace(query, " ? ", " $"+strconv.Itoa(i)+" ", 1)
			if t == query {
				break
			}
			i++
			query = t
		}
	default:
		return nil, errors.New("无此db drive 类型")
	}
	log.Println(query, args)

	return tx.db.Query(query, args...)

}

func (tx Tx) OrmConfig() OrmConf {
	return tx.ormConfig
}

func (tx TxEngine) Commit() error {
	return tx.tx.Commit()
}

func (tx TxEngine) Rollback() error {
	return tx.tx.Rollback()
}

type TxEngine struct {
	tx      *sql.Tx
	Base    EngineBase
	Extra   EngineExtra
	Table   EngineTable
	Classic EngineNative
}

func (e Engine) Begin() TxEngine {
	t, err := e.db.db.Begin()
	if err != nil {
		panic(err)
	}
	tx := Tx{
		db:       t,
		dbConfig: e.db.dbConfig,
	}

	fmt.Println(tx)
	return TxEngine{
		tx: t,
	}
}
