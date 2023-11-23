package lorm

import (
	"context"
	"database/sql"
)

type DbConfig interface {
	//根据db配置，打开db链接
	open() (*sql.DB, error)
	//根据 ctx生成dialecter，每种数据库类型各一种实现
	dialect(ctx ormContext) Dialecter
}

type DBer interface {
	//用db开启tx事务
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
	// 获取 corer 下面的dialecter的coreDb,coreTx 里面的 ctx
	getCtx() *ormContext
	getDB() *sql.DB
	//getTX() *sql.Tx

	//原生调用方法
	query(query string, args ...interface{}) *NativeQuery
	exec(query string, args ...interface{}) (sql.Result, error)
	prepare(query string) (*sql.Stmt, error)

	//lorm扩展方法
	c()
	r()
	u()
	d()

	beginTx(ctx context.Context, opts *sql.TxOptions) coreTx
	commit() error
	rollback() error
}

type Dialecter interface {
	// 获取coreDb,coreTx 里面的 ctx
	getCtx() *ormContext

	//对执行语句进行方言处理
	exec(query string, args ...interface{}) (string, []interface{})
	execBatch(query string, args [][]interface{}) (string, [][]interface{})

	query(query string, args ...interface{}) (string, []interface{})
	queryBatch(query string) string

	insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})
	insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})

	parse(c Clause) (string, error)

	//prepare(query string) (string, error)
}

//todo 下面未重构--------------

type DBer3 interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}
