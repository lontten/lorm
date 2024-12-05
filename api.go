package lorm

import (
	"context"
	"database/sql"
)

type DbConfig interface {
	//根据db配置，打开db链接
	open() (*sql.DB, error)
	//根据 ctx生成dialecter，每种数据库类型各一种实现
	dialect(ctx *ormContext, db DBer) Dialecter
}

type Stmter interface {
	query(args ...any) (*sql.Rows, error)
	exec(args ...any) (sql.Result, error)

	Exec(args ...any) (int64, error)
	QueryScan(args ...any) *NativePrepare
}

type DBer interface {
	// util
	ping() error
	// stmt
	prepare(query string) (Stmter, error)
	// tx
	beginTx(ctx context.Context, opts *sql.TxOptions) (DBer, error)
	commit() error
	rollback() error
	// db
	query(query string, args ...any) (*sql.Rows, error)
	exec(query string, args ...any) (sql.Result, error)
}

type Engine interface {
	// 上下文
	getDialect() Dialecter
	// stmt
	prepare(query string) (EngineStmt, error)
	// tx
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error)
	Begin() (Engine, error)
	Commit() error
	Rollback() error
}

type EngineStmt interface {
	// 上下文
	getDialect() Dialecter

	Exec(args ...any) (int64, error)
	QueryScan(args ...any) *NativePrepare
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
	// 上下文
	getDB() DBer
	getStmt() Stmter
	getCtx() *ormContext
	initContext() *ormContext
	hasErr() bool
	getErr() error
	// tx
	beginTx(ctx context.Context, opts *sql.TxOptions) error
	commit() error
	rollback() error
	// stmt
	prepare(query string) error
	// db
	exec(query string, args ...any) (string, []any)
	execBatch(query string, args [][]any) (string, [][]any)
	query(query string, args ...any) (string, []any)
	queryBatch(query string) string
	// utils
	parse(c Clause) (string, error)
	appendBaseToken(token baseToken)
	tableInsertGen()
	getSql() string
	//对执行语句进行方言处理
	//toSqlInsert()
}

type corer interface {
	//===----------------------------------------------------------------------===//
	// 获取上下文
	//===----------------------------------------------------------------------===//
	// 获取 corer 下面的dialecter的coreDb,coreTx 里面的 ctx
	getDB() *sql.DB
	getCtx() *ormContext
	hasErr() bool
	getErr() error
	getDialect() Dialecter

	//===----------------------------------------------------------------------===//
	// 具体执行
	//===----------------------------------------------------------------------===//
	//具体执行 创建事务
	//具体执行 创建提交事务
	doCommit() error
	//具体执行 事务回滚
	doRollback() error
	//具体执行 query，返回 *sql.Rows
	doQuery(query string, args ...any) (*sql.Rows, error)
	//具体执行 exec，返回 sql.Result
	doExec(query string, args ...any) (sql.Result, error)
	//具体执行 预处理 返回 *sql.Stmt
	doPrepare(query string) (Stmt, error)

	//===----------------------------------------------------------------------===//
	// 工具
	//===----------------------------------------------------------------------===//
	appendBaseToken(token baseToken)
}
type lnDB struct {
	core corer
}
