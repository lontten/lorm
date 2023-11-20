package lorm

import (
	"database/sql"
	"errors"
)

type DB struct {
	db *sql.DB
	tx *sql.Tx

	dbConfig DbConfig
	ctx      OrmContext

	dialect Dialect

	//where tokens
	whereTokens []string

	extraWhereSql string

	orderByTokens []string

	limit  int64
	offset int64

	//where values
	args      []interface{}
	batchArgs [][]interface{}

	// insert
	typ tableSqlType

	baseTokens []baseToken
}

func (db DB) getDB() *sql.DB {
	return db.db
}

func (db DB) getTX() *sql.Tx {
	return db.tx
}

func (db DB) rollback() error {
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

func (db DB) commit() error {
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
