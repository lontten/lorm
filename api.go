package lorm

import (
	"context"
	"database/sql"
)

type DbConfig interface {
	//根据db配置，打开db链接
	open() (*sql.DB, error)
	//根据 ctx生成dialecter，每种数据库类型各一种实现
	dialect(ctx *ormContext) Dialecter
}

type DBer interface {
	//用db开启tx事务
	BeginTx(ctx context.Context, opts *sql.TxOptions) TXer

	//原生调用方法
	Query(query string, args ...interface{}) *NativeQuery
	Exec(query string, args ...interface{}) (sql.Result, error)

	//lorm扩展方法
	Delete(v interface{}) OrmTableDelete
}

type TXer interface {
	Commit() error
	Rollback() error

	//原生调用方法
	Query(query string, args ...interface{}) *NativeQuery
	Exec(query string, args ...interface{}) (sql.Result, error)

	//lorm扩展方法
	Delete(v interface{}) OrmTableDelete
}

/*
*
直属lnDb

Dialecter 的实现有两种

	coreDb
	coreTx

内部属性为

	db      *sql.DB
	dialect Dialecter
*/
type corer interface {
	// 获取 corer 下面的dialecter的coreDb,coreTx 里面的 ctx
	getDB() *sql.DB
	getCtx() *ormContext
	hasErr() bool
	getErr() error

	getDialect() Dialecter

	appendBaseToken(token baseToken)

	//原生调用方法
	query(query string, args ...interface{}) *NativeQuery

	//lorm扩展方法
	c()
	r()
	u()
	d()

	//--------------具体执行-------------------------
	//具体执行 创建事务
	doBeginTx(ctx context.Context, opts *sql.TxOptions) coreTx
	//具体执行 创建提交事务
	doCommit() error
	//具体执行 事务回滚
	doRollback() error
	//具体执行 query，返回 *sql.Rows
	doQuery(query string, args ...interface{}) (*sql.Rows, error)
	//具体执行 exec，返回 sql.Result
	doExec(query string, args ...interface{}) (sql.Result, error)
	//具体执行 预处理 返回 *sql.Stmt
	doPrepare(query string) (*sql.Stmt, error)
}

/*
*

	Dialecter 的实现有两种

MysqlDialect
PgDialect

内部属性为 	ctx *ormContext
有	ormConf  OrmConf

	dbConfig DbConfig
	baseTokens []baseToken
*/
type Dialecter interface {
	// 获取coreDb,coreTx 里面的 ctx
	getCtx() *ormContext
	hasErr() bool
	getErr() error

	appendBaseToken(token baseToken)

	//对执行语句进行方言处理
	//toSqlInsert()

	exec(query string, args ...interface{}) (string, []interface{})
	execBatch(query string, args [][]interface{}) (string, [][]interface{})

	query(query string, args ...interface{}) (string, []interface{})
	queryBatch(query string) string

	insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})
	insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (string, []interface{})

	parse(c Clause) (string, error)
}

//todo 下面未重构--------------
