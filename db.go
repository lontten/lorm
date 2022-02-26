package lorm

import "database/sql"

type DB struct {
	db *sql.DB
	tx *sql.Tx

	dbConfig DbConfig
	ctx      OrmContext

	dialect Dialect

	//where tokens
	whereTokens []string

	extraWhereSql string

	orderByTokens []string

	limit  int64
	offset int64

	//where values
	args      []interface{}
	batchArgs [][]interface{}
}

func (db DB) OrmConf(c *OrmConf) DB {
	if c == nil {
		return db
	}
	db.ctx.conf = *c
	return db
}

type Result struct {
	num int64
	err error
}
type Resulter interface {
	Result() (int64, error)
	Err() error
}

func (r Result) Err() error {
	return r.err
}
func (r Result) Result() (int64, error) {
	return r.num, r.err
}

func (db DB) doQuery(query string, args ...interface{}) (*sql.Rows, error) {
	query, args = db.dialect.query(query, args...)
	return db.Db().Query(query, args...)
}

func (db DB) doExec(query string, args ...interface{}) (int64, error) {
	exec, err := db.Db().Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (db DB) doPrepare(query string) (Stmt, error) {
	stmt, err := db.Db().Prepare(query)
	return Stmt{stmt: stmt}, err
}

func (db DB) Db() DBer {
	if db.tx != nil {
		return db.tx
	} else {
		return db.db
	}
}
