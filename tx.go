package lorm

import (
	"context"
	"database/sql"
)

type LnTXer interface {
	Commit() error
	Rollback() error

	getDB() *sql.DB
	getTX() *sql.Tx

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
	return db.rollback()
}
func (db lnDB) Commit() error {
	return db.commit()
}

func (db lnDB) C() {
}

func (db lnDB) U() {
}

func (db lnDB) R() {
}

func (db lnDB) D() {
}
