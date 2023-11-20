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

	C()
	R()
	U()
	D()
}

func (db lnDB) OrmConf(c *OrmConf) lnDB {
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

func (db lnDB) doQuery(query string, args ...interface{}) (*sql.Rows, error) {
	query, args = db.dialect.query(query, args...)
	return db.Db().Query(query, args...)
}

func (db lnDB) doExec(query string, args ...interface{}) (int64, error) {
	exec, err := db.Db().Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (db lnDB) doPrepare(query string) (Stmt, error) {
	stmt, err := db.Db().Prepare(query)
	return Stmt{stmt: stmt}, err
}

func (db lnDB) Db() DBer {
	if db.tx != nil {
		return db.tx
	} else {
		return db.db
	}
}

func (db *lnDB) Do() Resulter {
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

func (db *lnDB) DoInsert() Resulter {
	return nil
}

func (db *lnDB) DoUpdate() Resulter {
	return nil
}

func (db *lnDB) DoDelete() Resulter {
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

func (db *lnDB) DoHas() Resulter {
	return nil
}

func (db *lnDB) DoSelect() Resulter {
	return nil
}

func (db *lnDB) DoCount() Resulter {
	return nil
}

func (db *lnDB) tDelete(t baseToken) {
	db.setTargetDest2TableName(t.dest)
}

func (db *lnDB) tPrimaryKey(t baseToken) {
	db.initPrimaryKeyName()
	db.ctx.initPrimaryKeyValues(t.pk)
	db.initByPrimaryKey()
	db.initExtra()
}

func (db *lnDB) tWhereModel(t baseToken) {
	db.initByModel(t.dest)
	db.initExtra()
}

func (db *lnDB) tWhereBuilder(t baseToken) {
	db.initByWhere(t.where)
	db.initExtra()
}
