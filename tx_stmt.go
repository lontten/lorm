package lorm

import (
	"database/sql"
)

// TXStmt -----------------TXStmt---------------------
type TXStmt struct {
	dialect Dialecter
}

func (s *TXStmt) getDialect() Dialecter {
	return s.dialect
}

func (s *TXStmt) Exec(args ...any) (int64, error) {
	return s.dialect.getStmt().Exec(args...)
}

func (s *TXStmt) QueryScan(args ...any) *NativePrepare {
	return &NativePrepare{
		db:   s,
		args: args,
	}
}

// -----------------DBStmt-end---------------------

// coreTXStmt -----------------coreTXStmt---------------------

type coreTXStmt struct {
	tx *sql.Stmt
}

func (s *coreTXStmt) query(args ...any) (*sql.Rows, error) {
	return s.tx.Query(args...)
}
func (s *coreTXStmt) exec(args ...any) (sql.Result, error) {
	return s.tx.Exec(args...)
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

// -----------------coreTXStmt-end---------------------
