package lorm

import (
	"context"
	"database/sql"
	"errors"
)

type LnTXer interface {
	Commit() error
	Rollback() error

	getDB() *sql.DB
	getTX() *sql.Tx

	//原生调用方法
	Query(query string, args ...interface{}) *NativeQuery
	Exec(query string, args ...interface{}) (rowsNum int64, err error)

	//lorm扩展方法
	C()
	R()
	U()
	D()
}

func (db lnDB) Begin() LnTXer {
	t, err := db.db.Begin()
	if err != nil {
		panic(err)
	}
	db.tx = t
	return db
}

func (db lnDB) BeginTx(ctx context.Context, opts *sql.TxOptions) LnTXer {
	t, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	db.tx = t
	return db
}

func (db lnDB) Rollback() error {
	if db.tx == nil {
		return errors.New("not in transaction")
	}
	err := db.tx.Rollback()
	if err != nil {
		return err
	}
	db.ctx.log.Println("rollback")
	return nil
}

func (db lnDB) Commit() error {
	if db.tx == nil {
		return errors.New("not in transaction")
	}
	err := db.tx.Commit()
	if err != nil {
		return err
	}
	db.ctx.log.Println("commit")
	return nil
}

func (db lnDB) C() {
}

func (db lnDB) U() {
}

func (db lnDB) R() {
}

func (db lnDB) D() {
}
