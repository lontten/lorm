package lorm

import (
	"database/sql"
)

type Stmt struct {
	stmt *sql.Stmt

	ctx OrmContext
}

func (db DB) Prepare(query string) (Stmt, error) {
	return db.dialect.prepare(query)
}

func (s *Stmt) Exec(args ...interface{}) (int64, error) {
	exec, err := s.stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (s *Stmt) Query(args ...interface{}) Prepare {

	return Prepare{}
}

type Prepare struct {
}

func (p Prepare) ScanOne(v interface{}) (int64, error) {
	return 0, nil
}

func (p Prepare) ScanList(v interface{}) (int64, error) {
	return 0, nil
}

func (p Prepare) ScanFirst(v interface{}) (int64, error) {
	return 0, nil
}
