package lorm

import (
	"database/sql"
	"errors"
)

type lnDB struct {
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

func (db lnDB) getDB() *sql.DB {
	return db.db
}

func (db lnDB) getTX() *sql.Tx {
	return db.tx
}

func (db lnDB) rollback() error {
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

func (db lnDB) commit() error {
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
