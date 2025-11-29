package lorm

import (
	"database/sql"
)

type coreDBStmt struct {
	db      *sql.Stmt
	dialect Dialecter
}

func (db *coreDBStmt) init() Stmter {
	return &coreDBStmt{
		db:      db.db,
		dialect: db.dialect.copyContext(),
	}
}
func (db *coreDBStmt) getCtx() *ormContext {
	return db.dialect.getCtx()
}
func (db *coreDBStmt) getDialect() Dialecter {
	return db.dialect
}

func (db *coreDBStmt) query(args ...any) (*sql.Rows, error) {
	return db.db.Query(args...)
}
func (db *coreDBStmt) exec(args ...any) (sql.Result, error) {
	return db.db.Exec(args...)
}
