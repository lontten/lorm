package lorm

import (
	"database/sql"
)

type coreTXStmt struct {
	tx      *sql.Stmt
	dialect Dialecter
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

func (s coreTXStmt) Exec(args ...any) (int64, error) {
	exec, err := s.exec(args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (s *coreTXStmt) QueryScan(args ...any) *NativePrepare {
	return &NativePrepare{
		db:   s,
		args: args,
	}
}
