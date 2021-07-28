package lorm

import (
	"fmt"
	"log"
	"reflect"
)

//对 query exec 的简单封装
type EngineClassic struct {
	db DBer
	lormConf Lorm
	dialect  Dialect

	context OrmContext
}

type ClassicQuery struct {
	base  EngineClassic
	query string
	args  []interface{}
}

type ClassicExec struct {
	base EngineClassic
}

func (engine EngineClassic) Query(query string, args ...interface{}) *ClassicQuery {
	return &ClassicQuery{base: engine, query: query, args: args}
}

func (q ClassicQuery) GetOne(dest interface{}) (rowsNum int64, err error) {
	_, err = checkScanTypeLn(reflect.TypeOf(dest))
	if err != nil {
		return 0, err
	}

	query := q.query
	args := q.args
	log.Println(query, args)
	rows, err := q.base.db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.lormConf.ScanLn(rows, dest)
}

func (q ClassicQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	query := q.query
	args := q.args
	fmt.Println(query+"k")
	fmt.Println(args...)
	query = q.base.dialect.ToDialectSql(query)
	fmt.Println(query)
	rows, err := q.base.db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.lormConf.Scan(rows, dest)
}

func (engine EngineClassic) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return engine.db.exec(query, args...)
}
