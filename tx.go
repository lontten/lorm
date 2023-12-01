package lorm

import (
	"context"
	"database/sql"
)

type coreTx struct {
	tx      *sql.Tx
	dialect Dialecter
}

func (tx coreTx) getCtx() *ormContext {
	return tx.dialect.getCtx()
}
func (tx coreTx) hasErr() bool {
	return tx.dialect.hasErr()
}
func (db coreTx) getDialect() Dialecter {
	return db.dialect
}

func (tx coreTx) appendBaseToken(token baseToken) {
	tx.dialect.appendBaseToken(token)
}
func (tx coreTx) getDB() *sql.DB {
	panic("tx no db")
	return nil
}

func (tx coreTx) doRollback() error {
	err := tx.tx.Rollback()
	if err != nil {
		return err
	}
	return nil
}

func (tx coreTx) doCommit() error {
	err := tx.tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
func (tx coreTx) doBeginTx(ctx context.Context, opts *sql.TxOptions) coreTx {
	panic("tx err again beginTX")
	return tx
}
func (tx coreTx) c() {
}
func (tx coreTx) r() {
}
func (tx coreTx) u() {
}
func (tx coreTx) d() {
}
func (tx coreTx) query(query string, args ...interface{}) *NativeQuery {
	return &NativeQuery{core: tx, query: query, args: args}
}
func (tx coreTx) doExec(query string, args ...interface{}) (sql.Result, error) {
	query, args = tx.dialect.exec(query, args...)
	return tx.tx.Exec(query, args...)
}
func (tx coreTx) doPrepare(query string) (*sql.Stmt, error) {
	return tx.tx.Prepare(query)
}

//todo 下面未重构--------------
