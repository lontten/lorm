package lorm

import (
	"database/sql"
)

type Stmt struct {
	stmt *sql.Stmt
}

func (db DB) Prepare(query string) (s Stmt, err error) {
	if db.isTx {
		stmt, err := db.tx.Prepare(query)
		if err != nil {
			return s, err
		}
		s.stmt = stmt
		return s, nil
	}
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return
	}
	s.stmt = stmt
	return
}

func (s *Stmt) Query(args ...interface{}) (int64, error) {
	return 0, nil
}

func (s *Stmt) Exec(args ...interface{}) Prepare {

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
