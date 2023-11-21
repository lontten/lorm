package lorm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
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

// ----------LnDB-------------

type coreDb struct {
	db      *sql.DB
	dialect Dialecter
}

func (db coreDb) getDB() *sql.DB {
	return db.db
}

func (db coreDb) beginTx(ctx context.Context, opts *sql.TxOptions) corer {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	return coreTx{tx: tx}
}

func (db coreDb) rollback() error {
	return errors.New("this not tx")
}

func (db coreDb) commit() error {
	return errors.New("this not tx")
}

func (db coreDb) c() {
}
func (db coreDb) r() {
}
func (db coreDb) u() {
}
func (db coreDb) d() {
}

func (db coreDb) query(query string, args ...interface{}) *NativeQuery {
	rows, err := db.db.Query()
	return &NativeQuery{core: db, query: query, args: args}
}

func (db coreDb) exec(query string, args ...interface{}) (rowsNum int64, err error) {
	query, args = db.dialect.exec(query, args...)
	//return db.doExec(query, args...)
	return 0, nil
}

//----------LnDB-------------

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
	Result() (int64, error) //sql执行影响了多少行数时和err
	Err() error             //当用户不在意，sql执行影响了多少行数时，可以使用这个直接获取err，不用再想之前一样还要用_接受
}

func (r Result) Err() error {
	return r.err
}
func (r Result) Result() (int64, error) {
	return r.num, r.err
}

func (db lnDB) doQuery(query string, args ...interface{}) (*sql.Rows, error) {
	//query, args = db.dialect.query(query, args...)
	//return db.Db().Query(query, args...)
	return nil, nil
}

func (db lnDB) doExec(query string, args ...interface{}) (int64, error) {
	//exec, err := db.Db().Exec(query, args...)
	//if err != nil {
	//	return 0, err
	//}
	//return exec.RowsAffected()
	return 0, nil
}

func (db lnDB) doPrepare(query string) (Stmt, error) {
	//stmt, err := db.Db().Prepare(query)
	//return Stmt{stmt: stmt}, err
	return Stmt{}, nil
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
