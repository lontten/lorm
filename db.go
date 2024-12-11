package lorm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type coreDB struct {
	db      *sql.DB
	dialect Dialecter
}

func (db *coreDB) init() Engine {
	return &coreDB{
		db:      db.db,
		dialect: db.dialect.initContext(),
	}
}

func (db *coreDB) ping() error {
	return db.db.Ping()
}

func (db *coreDB) getCtx() *ormContext {
	return db.dialect.getCtx()
}
func (db *coreDB) getDialect() Dialecter {
	return db.dialect
}
func (db *coreDB) query(query string, args ...any) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}
func (db *coreDB) exec(query string, args ...any) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *coreDB) prepare(query string) (Stmter, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &coreDBStmt{
		db:      stmt,
		dialect: db.dialect,
	}, nil
}

func (db *coreDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &coreTX{
		tx:      tx,
		dialect: db.dialect,
	}, nil
}

func (db *coreDB) Begin() (Engine, error) {
	return db.BeginTx(context.Background(), nil)
}

func (db *coreDB) Commit() error {
	return errors.New("this not tx")
}

func (db *coreDB) Rollback() error {
	return errors.New("this not tx")
}

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
