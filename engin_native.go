package lorm

//对 query exec 的简单封装
type EngineNative struct {
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
	if err = q.base.context.err; err != nil {
		return 0, err
	}
	q.base.context.initScanDestSlice(dest)
	q.base.context.checkScanDestField()
	if err = q.base.context.err; err != nil {
		return 0, err
	}

	query := q.query
	args := q.args
	rows, err := q.base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.context.ScanLn(rows)
}

func (q ClassicQuery) GetList(dest interface{}) (rowsNum int64, err error) {
	if err = q.base.context.err; err != nil {
		return 0, err
	}
	q.base.context.initScanDestSlice(dest)
	q.base.context.checkScanDestField()

	if err = q.base.context.err; err != nil {
		return 0, err
	}

	query := q.query
	args := q.args
	rows, err := q.base.dialect.query(query, args...)
	if err != nil {
		return 0, err
	}
	return q.base.context.Scan(rows)
}

func (engine EngineNative) Exec(query string, args ...interface{}) (rowsNum int64, err error) {
	return engine.dialect.exec(query, args...)
}
