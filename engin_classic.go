package lorm

import "log"

type EngineClassic struct {
	context   OrmContext
	db DbPool
}

type ClassicQuery struct {
	base  EngineClassic
	query string
	args  []interface{}
}

type ClassicExec struct {
	base *EngineClassic
}

func (engine EngineClassic) Query(query string, args ...interface{}) *ClassicQuery {
	return &ClassicQuery{base: engine, query: query, args: args}
}

func (q ClassicQuery) GetOne(dest interface{}) (rowsNum int64, err error) {
	query := q.query
	args := q.args
	log.Println(query,args)
	rows, err := q.base.db.db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	return StructScanLn(rows, dest)
}

func (q ClassicQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	query := q.query
	args := q.args
	log.Println(query,args)
	rows, err := q.base.db.db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	return StructScan(rows, dest)
}

func (engine EngineClassic) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return engine.db.Exec(query, args)
}
