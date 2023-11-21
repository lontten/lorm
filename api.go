package lorm

import (
	"context"
	"database/sql"
)

type DbConfig interface {
	Open() (*sql.DB, error)
	dialect(ctx *ormContext, pc *PoolConf) Dialecter
}

type DBer interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) TXer

	//原生调用方法
	Query(query string, args ...interface{}) *NativeQuery
	Exec(query string, args ...interface{}) (rowsNum int64, err error)

	//lorm扩展方法
	C()
	R()
	U()
	D()
}

type TXer interface {
	Commit() error
	Rollback() error

	//原生调用方法
	Query(query string, args ...interface{}) *NativeQuery
	Exec(query string, args ...interface{}) (rowsNum int64, err error)

	//lorm扩展方法
	C()
	R()
	U()
	D()
}

type corer interface {
	getDB() *sql.DB
	//getTX() *sql.Tx

	//原生调用方法
	query(query string, args ...interface{}) (*sql.Rows, error)
	exec(query string, args ...interface{}) (rowsNum int64, err error)

	//lorm扩展方法
	c()
	r()
	u()
	d()

	beginTx(ctx context.Context, opts *sql.TxOptions) corer
	commit() error
	rollback() error
}

type DBer3 interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}

type Dialecter interface {
	exec(query string, args ...interface{}) (string, []interface{})
	execBatch(query string, args [][]interface{}) (string, [][]interface{})

	query(query string, args ...interface{}) (string, []interface{})
	queryBatch(query string) string

	insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})
	insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})

	parse(c Clause) (string, error)

	//prepare(query string) (string, error)
}
