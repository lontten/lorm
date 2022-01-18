package lorm

import "database/sql"

type DB struct {
	db   *sql.DB
	tx   *sql.Tx
	isTx bool

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

func (db DB) OrmConf(c *OrmConf) DB {
	if c == nil {
		return db
	}
	db.ctx.conf = *c
	return db
}
