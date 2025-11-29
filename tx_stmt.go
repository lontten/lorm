package lorm

import (
	"database/sql"
)

type coreTXStmt struct {
	tx      *sql.Stmt
	dialect Dialecter
}

func (db *coreTXStmt) init() Stmter {
	return &coreTXStmt{
		tx:      db.tx,
		dialect: db.dialect.copyContext(),
	}
}
func (db *coreTXStmt) getCtx() *ormContext {
	return db.dialect.getCtx()
}
func (db *coreTXStmt) getDialect() Dialecter {
	return db.dialect
}

func (db *coreTXStmt) query(args ...any) (*sql.Rows, error) {
	return db.tx.Query(args...)
}
func (db *coreTXStmt) exec(args ...any) (sql.Result, error) {
	return db.tx.Exec(args...)
}
