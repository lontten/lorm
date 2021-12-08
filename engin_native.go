package lorm

import "github.com/pkg/errors"

type ClassicQuery struct {
	base  Engine
	query string
	args  []interface{}
}

type ClassicExec struct {
	base Engine
}

func (engine Engine) Query(query string, args ...interface{}) *ClassicQuery {
	return &ClassicQuery{base: engine, query: query, args: args}
}

func (q ClassicQuery) GetOne(dest interface{}) (rowsNum int64, err error) {
	if err = q.base.ctx.err; err != nil {
		return 0, err
	}
	q.base.ctx.initScanDestOne(dest)
	if q.base.ctx.scanIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	q.base.ctx.checkScanDestField()
	if err = q.base.ctx.err; err != nil {
		return 0, err
	}

	query := q.query
	args := q.args
	rows, err := q.base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.ctx.ScanLn(rows)
}

func (q ClassicQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	if err = q.base.ctx.err; err != nil {
		return 0, err
	}
	q.base.ctx.initScanDestList(dest)
	q.base.ctx.checkScanDestField()

	if err = q.base.ctx.err; err != nil {
		return 0, err
	}

	query := q.query
	args := q.args
	rows, err := q.base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.ctx.Scan(rows)
}

func (engine Engine) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return engine.dialect.exec(query, args...)
}
