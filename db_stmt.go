package lorm

import (
	"database/sql"
)

// DBStmt -----------------DBStmt---------------------
type DBStmt struct {
	dialect Dialecter
}

func (s *DBStmt) getDialect() Dialecter {
	return s.dialect
}

func (s *DBStmt) Exec(args ...any) (int64, error) {
	return s.dialect.getStmt().Exec(args...)
}

func (s *DBStmt) QueryScan(args ...any) *NativePrepare {
	return &NativePrepare{
		db:   s,
		args: args,
	}
}

// -----------------DBStmt-end---------------------

// coreDBStmt -----------------coreDBStmt---------------------

type coreDBStmt struct {
	db *sql.Stmt
}

func (s *coreDBStmt) query(args ...any) (*sql.Rows, error) {
	return s.db.Query(args...)
}
func (s *coreDBStmt) exec(args ...any) (sql.Result, error) {
	return s.db.Exec(args...)
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

// -----------------coreDBStmt-end---------------------
