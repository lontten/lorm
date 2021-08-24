package lorm

import (
	"reflect"
)

//对 query exec 的简单封装
type EngineNative struct {
	ormConf OrmConf
	dialect Dialect
	context OrmContext
}

type ClassicQuery struct {
	base  EngineNative
	query string
	args  []interface{}
}

type ClassicExec struct {
	base EngineNative
}

func (engine EngineNative) Query(query string, args ...interface{}) *ClassicQuery {
	return &ClassicQuery{base: engine, query: query, args: args}
}

func (q ClassicQuery) GetOne(dest interface{}) (rowsNum int64, err error) {
	_, err = checkScanTypeLn(reflect.TypeOf(dest))
	if err != nil {
		return 0, err
	}

	query := q.query
	args := q.args
	rows, err := q.base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.context.core.ScanLn(rows, dest)
}

func (q ClassicQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	query := q.query
	args := q.args
	rows, err := q.base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.context.core.Scan(rows, dest)
}

func (engine EngineNative) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return engine.dialect.exec(query, args...)
}
