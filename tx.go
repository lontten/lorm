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

	t, err := db.db.Begin()
	if err != nil {
		panic(err)
	}
	db.tx = t
	db.isTx = true
	db.dialect.SetDb(t)
	return db
}

func (db DB) BeginTx(ctx context.Context, opts *sql.TxOptions) DB {
	t, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	db.tx = t
	db.isTx = true
	db.dialect.SetDb(t)
	return db
}

func (db DB) Rollback() error {
	if !db.isTx {
		return errors.New("not in transaction")
	}
	err := db.tx.Rollback()
	if err != nil {
		return err
	}
	db.isTx = false
	db.dialect.SetDb(db.db)
	return nil
}

func (db DB) Commit() error {
	if !db.isTx {
		return errors.New("not in transaction")
	}
	err := db.tx.Commit()
	if err != nil {
		return err
	}
	db.isTx = false
	db.dialect.SetDb(db.db)
	return nil
}
