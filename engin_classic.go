package lorm

import (
	"log"
	"reflect"
)

//对 query exec 的简单封装
type EngineClassic struct {
	context OrmContext
	db      DBer
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
	fieldNamePrefix := q.base.db.OrmConfig().FieldNamePrefix
	log.Println(query,args)
	rows, err := q.base.db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return StructScanLn(rows, dest,fieldNamePrefix)
}

func (q ClassicQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	query := q.query
	args := q.args
	fieldNamePrefix := q.base.db.OrmConfig().FieldNamePrefix
	log.Println(query,args)
	rows, err := q.base.db.query(query, args...)
	if err != nil {
		return 0, err
	}
	return StructScan(rows, dest,fieldNamePrefix)
}

func (engine EngineClassic) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return engine.db.exec(query, args...)
}
