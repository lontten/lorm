package lorm

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
)

type DB struct {
	ctx OrmContext
	db        *sql.DB
	dbConfig  DbConfig
	ormConfig OrmConfig
}

func (db DB) OrmConfig() OrmConfig {
	return db.ormConfig
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
	switch db.dbConfig.DriverName() {
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
	db      DBer
	context OrmContext
}

type OrmFrom struct {
	db      DBer
	context OrmContext
}

type OrmWhere struct {
	db      DBer
	context OrmContext
}

