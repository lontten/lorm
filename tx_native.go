package lorm

import "github.com/pkg/errors"

type NativeTxQuery struct {
	base  Tx
	query string
	args  []interface{}
}

func (tx Tx) Query(query string, args ...interface{}) *NativeTxQuery {
	return &NativeTxQuery{base: tx, query: query, args: args}
}

func (q NativeTxQuery) GetOne(dest interface{}) (rowsNum int64, err error) {
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

	rows, err := q.base.dialect.query(q.query, q.args...)
	if err != nil {
		return 0, err
	}
	return q.base.ctx.ScanLn(rows)
}

func (q NativeTxQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	if err = q.base.ctx.err; err != nil {
		return 0, err
	}
	q.base.ctx.initScanDestList(dest)
	q.base.ctx.checkScanDestField()

	if err = q.base.ctx.err; err != nil {
		return 0, err
	}

	rows, err := q.base.dialect.query(q.query, q.args...)
	if err != nil {
		return 0, err
	}
	return q.base.ctx.Scan(rows)
}

func (tx Tx) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return tx.dialect.exec(query, args...)
}
