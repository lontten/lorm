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
	SetDb(db DBer)

	exec(query string, args ...interface{}) (int64, error)
	execBatch(query string, args [][]interface{}) (int64, error)
	query(query string, args ...interface{}) (*sql.Rows, error)
	queryBatch(query string) (*sql.Stmt, error)

	insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (int64, error)
	insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (int64, error)
}

type SqlUtil interface {
	selectArgsArr2SqlStr(args []string)
	tableWhereArgs2SqlStr(args []string, c OrmConf) string
	tableCreateArgs2SqlStr(args []string) string
	tableUpdateArgs2SqlStr(args []string) string
	tableWherePrimaryKey2SqlStr(ids []string, c OrmConf) string
}
