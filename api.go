package lorm

import (
	"database/sql"
)

type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

func Select(q Queryer, dest interface{}, query string, args ...interface{}) error {
	rows, err := q.Query(query, args...)
	if err != nil {
		return err
	}

	defer rows.Close()
	_, err = StructScan(rows, dest)
	return err
}

func Exec(e Execer, query string, args ...interface{}) (sql.Result, error) {
	result, err := e.Exec(query, args...)
	if err != nil {
		return nil, err
	}
	return result, err
}



type DBer interface {
	exec(query string, args ...interface{}) (int64, error)
	query(query string, args ...interface{}) (*sql.Rows, error)
	OrmConfig() OrmConfig
}
