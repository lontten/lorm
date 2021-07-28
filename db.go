package lorm

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
)

type DB struct {
	db        *sql.DB
	dbConfig  DbConfig
}



func (db DB) exec(query string, args ...interface{}) (int64, error) {

	switch db.dbConfig.DriverName() {
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
	Log.Println(query, args)

	exec, err := db.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (db DB) query(query string, args ...interface{}) (*sql.Rows, error) {
	Log.Println("sql",query,args)
	return db.db.Query(query, args...)
}

type OrmContext struct {
	query  *strings.Builder
	args   []interface{}
	startd bool
	err error
	log int
}

type OrmSelect struct {
	db DBer
	lormConf Lorm
	dialect  Dialect
	context OrmContext
}

type OrmFrom struct {
	db DBer
	lormConf Lorm
	dialect  Dialect
	context OrmContext
}

type OrmWhere struct {
	db DBer
	lormConf Lorm
	dialect  Dialect
	context OrmContext
}

