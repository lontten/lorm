package lorm

import (
	"context"
	"database/sql"
)

type tableSqlType int

const (
	dInsert tableSqlType = iota
	dUpdate
	dDelete
	dSelect
	dHas
	dCount
)

type LnDBer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) LnTXer
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

func (db *DB) Do() Resulter {
	switch db.typ {
	case dInsert:
		return db.DoInsert()
	case dUpdate:
		return db.DoUpdate()
	case dDelete:
		return db.DoDelete()
	case dSelect:
		return db.DoSelect()
	case dHas:
		return db.DoHas()
	case dCount:
		return db.DoCount()
	}
	return nil
}

func (db *DB) DoInsert() Resulter {
	return nil
}

func (db *DB) DoUpdate() Resulter {
	return nil
}

func (db *DB) DoDelete() Resulter {
	for _, token := range db.baseTokens {
		switch token.typ {
		case tDelete:
			db.tDelete(token)
		case tPrimaryKey:
			db.tPrimaryKey(token)
		case tWhereModel:
			db.tPrimaryKey(token)
		case tWhereBuilder:
			db.tWhereBuilder(token)
		}
	}
	num, err := db.doDel()
	return Result{num: num, err: err}
}

func (db *DB) DoHas() Resulter {
	return nil
}

func (db *DB) DoSelect() Resulter {
	return nil
}

func (db *DB) DoCount() Resulter {
	return nil
}

func (db *DB) tDelete(t baseToken) {
	db.setTargetDest2TableName(t.dest)
}

func (db *DB) tPrimaryKey(t baseToken) {
	db.initPrimaryKeyName()
	db.ctx.initPrimaryKeyValues(t.pk)
	db.initByPrimaryKey()
	db.initExtra()
}

func (db *DB) tWhereModel(t baseToken) {
	db.initByModel(t.dest)
	db.initExtra()
}

func (db *DB) tWhereBuilder(t baseToken) {
	db.initByWhere(t.where)
	db.initExtra()
}
