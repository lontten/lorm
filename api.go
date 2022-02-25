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
	exec(query string, args ...interface{}) (string, []interface{})
	execBatch(query string, args [][]interface{}) (string, [][]interface{})

	query(query string, args ...interface{}) (string, []interface{})
	queryBatch(query string) string

	insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})
	insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})

	parse(c Clause) (string, error)

	//prepare(query string) (string, error)
}
