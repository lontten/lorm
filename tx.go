package lorm

import (
	"context"
	"database/sql"
)

type coreTx struct {
	tx      *sql.Tx
	dialect Dialecter
}

func (tx coreTx) getDB() *sql.DB {
	panic("tx no db")
	return nil
}

func (tx coreTx) rollback() error {
	err := tx.tx.Rollback()
	if err != nil {
		return err
	}
	//db.ctx.log.Println("rollback")
	return nil
}

func (tx coreTx) commit() error {
	err := tx.tx.Commit()
	if err != nil {
		return err
	}
	//db.ctx.log.Println("commit")
	return nil
}
func (tx coreTx) beginTx(ctx context.Context, opts *sql.TxOptions) coreTx {
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
func (tx coreTx) exec(query string, args ...interface{}) (sql.Result, error) {
	query, args = tx.dialect.exec(query, args...)
	return tx.tx.Exec(query, args...)
}

//todo 下面未重构--------------
