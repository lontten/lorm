package lorm

import "github.com/pkg/errors"

// todo 下面未重构--------------
type NativeQuery struct {
	core  corer
	query string
	args  []interface{}
}

func (q NativeQuery) ScanOne(dest interface{}) (rowsNum int64, err error) {
	if err = q.core.getCtx().err; err != nil {
		return 0, err
	}
	q.core.getCtx().initScanDestOne(dest)
	if q.core.getCtx().scanIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	q.core.getCtx().checkScanDestField()
	if err = q.core.getCtx().err; err != nil {
		return 0, err
	}

	rows, err := q.core.doQuery(q.query, q.args...)

	if err != nil {
		return 0, err
	}
	return q.core.getCtx().ScanLn(rows)
}

func (q NativeQuery) ScanList(dest interface{}) (rowsNum int64, err error) {
	if err = q.core.getCtx().err; err != nil {
		return 0, err
	}
	q.core.getCtx().initScanDestList(dest)
	q.core.getCtx().checkScanDestField()

	if err = q.core.getCtx().err; err != nil {
		return 0, err
	}

	rows, err := q.base.doQuery(q.query, q.args...)
	if err != nil {
		return 0, err
	}
	return q.core.getCtx().Scan(rows)
}

//func (db dbCore) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
//	query, args = db.dialect.exec(query, args...)
//	return db.doExec(query, args...)
//}
//func (db dbCore) Query(query string, args ...interface{}) *NativeQuery {
//	return &NativeQuery{base: db, query: query, args: args}
//}
