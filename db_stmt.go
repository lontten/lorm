package lorm

import (
	"database/sql"
)

type coreDBStmt struct {
	db      *sql.Stmt
	dialect Dialecter
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

func (s coreDBStmt) Exec(args ...any) (int64, error) {
	exec, err := s.exec(args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (s *coreDBStmt) QueryScan(args ...any) *NativePrepare {
	return &NativePrepare{
		db:   s,
		args: args,
	}
}
