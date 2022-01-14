package lorm

import (
	"context"
	"database/sql"
)

func (db DB) Begin() DB {
	t, err := db.dbBase.Begin()
	if err != nil {
		panic(err)
	}
	db.db = t
	return db
}

func (db DB) BeginTx(ctx context.Context, opts *sql.TxOptions) DB {
	t, err := db.dbBase.BeginTx(ctx, opts)
	if err != nil {
		panic(err)
	}
	db.db = t
	return db
}
