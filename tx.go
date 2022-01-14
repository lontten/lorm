package lorm

import (
	"context"
	"database/sql"
	"errors"
)

func (db DB) Begin() DB {
	if db.isTx {
		return db
	}
	t, err := db.db.(*sql.DB).Begin()
	if err != nil {
		panic(err)
	}
	db.db = t
	db.isTx = true
	db.dialect.SetDb(t)
	return db
}

func (db DB) BeginTx(ctx context.Context, opts *sql.TxOptions) DB {
	t, err := db.db.(*sql.DB).BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	db.db = t
	db.isTx = true
	db.dialect.SetDb(t)
	return db
}

func (db DB) Rollback() error {
	if !db.isTx {
		return errors.New("not in transaction")
	}
	return db.db.(*sql.Tx).Rollback()
}

func (db DB) Commit() error {
	if !db.isTx {
		return errors.New("not in transaction")
	}
	return db.db.(*sql.Tx).Commit()
}
