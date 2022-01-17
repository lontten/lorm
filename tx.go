package lorm

import (
	"context"
	"database/sql"
	"errors"
)

func (db DB) Begin() DB {
	t, err := db.db.Begin()
	if err != nil {
		panic(err)
	}

	d := DB{
		db:       db.db,
		tx:       t,
		isTx:     true,
		dbConfig: db.dbConfig,
		ctx:      db.ctx.Copy(),
		dialect:  db.dialect.Copy(t),
	}
	return d
}

func (db DB) BeginTx(ctx context.Context, opts *sql.TxOptions) DB {
	t, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}

	return DB{
		db:       db.db,
		tx:       t,
		isTx:     true,
		dbConfig: db.dbConfig,
		ctx:      db.ctx.Copy(),
		dialect:  db.dialect.Copy(t),
	}
}

func (db DB) Rollback() error {
	if !db.isTx {
		return errors.New("not in transaction")
	}
	err := db.tx.Rollback()
	if err != nil {
		return err
	}
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
	return nil
}
