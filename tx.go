package lorm

import (
	"context"
	"database/sql"
)

type Tx struct {
	tx       *sql.Tx
	dbConfig DbConfig
	ctx      OrmContext

	dialect Dialect

	//where tokens
	whereTokens []string

	extraWhereSql []byte

	orderByTokens []string

	limit  int64
	offset int64

	//where values
	args      []interface{}
	batchArgs [][]interface{}
}

func (tx Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (db DB) Begin() Tx {
	t, err := db.db.Begin()
	if err != nil {
		panic(err)
	}
	return Tx{
		tx:       t,
		dbConfig: db.dbConfig,
		ctx:      db.ctx,
		dialect:  db.dialect,
	}
}

func (db DB) BeginTx(ctx context.Context, opts *sql.TxOptions) Tx {
	t, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	return Tx{
		tx:       t,
		dbConfig: db.dbConfig,
		ctx:      db.ctx,
		dialect:  db.dialect,
	}
}
