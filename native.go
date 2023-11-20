package lorm

import "github.com/pkg/errors"

type NativeQuery struct {
	base  lnDB
	query string
	args  []interface{}
}

func (db lnDB) Query(query string, args ...interface{}) *NativeQuery {
	return &NativeQuery{base: db, query: query, args: args}
}

func (q NativeQuery) ScanOne(dest interface{}) (rowsNum int64, err error) {
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

	rows, err := q.base.doQuery(q.query, q.args...)

	if err != nil {
		return 0, err
	}
	return q.base.ctx.ScanLn(rows)
}

func (q NativeQuery) ScanList(dest interface{}) (rowsNum int64, err error) {
	if err = q.base.ctx.err; err != nil {
		return 0, err
	}
	q.base.ctx.initScanDestList(dest)
	q.base.ctx.checkScanDestField()

	if err = q.base.ctx.err; err != nil {
		return 0, err
	}

	rows, err := q.base.doQuery(q.query, q.args...)
	if err != nil {
		return 0, err
	}
	return q.base.ctx.Scan(rows)
}

func (db lnDB) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	query, args = db.dialect.exec(query, args...)
	return db.doExec(query, args...)
}
