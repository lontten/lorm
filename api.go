package lorm

import (
	"database/sql"
)

type DBer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

type Dialect interface {
	Copy(db DBer) Dialect

	exec(query string, args ...interface{}) (int64, error)
	execBatch(query string, args [][]interface{}) (int64, error)
	query(query string, args ...interface{}) (*sql.Rows, error)
	queryBatch(query string) (*sql.Stmt, error)

	insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (int64, error)
	insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (int64, error)
}
