package lorm

import (
	"context"
	"database/sql"
	"errors"
)

// DB -----------------DB---------------------
type DB struct {
	dialect Dialecter
}

func (db *DB) getDialect() Dialecter {
	return db.dialect
}
func (db *DB) prepare(query string) (EngineStmt, error) {
	dialect := db.getDialect()
	err := dialect.prepare(query)
	if err != nil {
		return nil, err
	}
	return &DBStmt{dialect: dialect}, nil
}

func (db *DB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error) {
	dialect := db.getDialect()
	err := dialect.beginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &TX{dialect: dialect}, nil
}

func (db *DB) Begin() (Engine, error) {
	return db.BeginTx(context.Background(), nil)
}

func (db *DB) Commit() error {
	return errors.New("this not tx")
}

func (db *DB) Rollback() error {
	return errors.New("this not tx")
}

// -----------------DB-end---------------------

// coreDB -----------------coreDB---------------------

type coreDB struct {
	db *sql.DB
}

func (db *coreDB) ping() error {
	return db.db.Ping()
}

func (db *coreDB) prepare(query string) (Stmter, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &coreDBStmt{db: stmt}, nil
}

func (db *coreDB) query(query string, args ...any) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

func (db *coreDB) exec(query string, args ...any) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *coreDB) beginTx(ctx context.Context, opts *sql.TxOptions) (DBer, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &coreTX{tx: tx}, nil
}

func (db *coreDB) commit() error {
	return errors.New("this is db")
}

func (db *coreDB) rollback() error {
	return errors.New("this is db")
}

// -----------------coreDB-end---------------------

//todo 下面未重构--------------

//----------LnDB-------------

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

func (db *lnDB) tDelete(t baseToken) {
	db.setNameDest(t.dest)
}

func (d *MysqlDialect) tPrimaryKey(t baseToken) {
	d.initPrimaryKeyName()
	d.ctx.initPrimaryKeyValues(t.pk)
	d.initByPrimaryKey()
	d.initExtra()
}

func (d *MysqlDialect) tWhereModel(t baseToken) {
	d.initByModel(t.dest)
	d.initExtra()
}

//func (ldb *MysqlDialect) tWhereBuilder(t baseToken) {
//	ldb.initByWhere(t.wb)
//	ldb.initExtra()
//}
