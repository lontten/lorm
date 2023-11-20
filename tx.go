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
	commit() error
	rollback() error
}

func (db DB) Begin() LnTXer {
	t, err := db.db.Begin()
	if err != nil {
		panic(err)
	}
	db.tx = t
	return db
}

func (db DB) BeginTx(ctx context.Context, opts *sql.TxOptions) LnTXer {
	t, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	db.tx = t
	return db
}

func (db DB) Rollback() error {
	return db.rollback()
}
func (db DB) Commit() error {
	return db.commit()
}
