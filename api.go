// Copyright (c) 2024 lontten
// lorm is licensed under Mulan PSL v2.
// You can use this software according to the terms and conditions of the Mulan PSL v2.
// You may obtain a copy of Mulan PSL v2 at:
// http://license.coscl.org.cn/MulanPSL2
// THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
// EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
// MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
// See the Mulan PSL v2 for more details.

package lorm

import "database/sql"

import (
	"context"
)

type DbConfig interface {
	//根据db配置，打开db链接
	open() (*sql.DB, error)
	//根据 ctx生成dialecter，每种数据库类型各一种实现
	dialect(ctx *ormContext) Dialecter
}

type Stmter interface {
	init() Stmter
	getCtx() *ormContext
	getDialect() Dialecter

	query(args ...any) (*sql.Rows, error)
	exec(args ...any) (sql.Result, error)

	Exec(args ...any) (int64, error)
	QueryScan(args ...any) *NativePrepare
}
type Engine interface {
	init() Engine
	ping() error
	getCtx() *ormContext
	getDialect() Dialecter

	query(query string, args ...any) (*sql.Rows, error)
	exec(query string, args ...any) (sql.Result, error)
	prepare(query string) (Stmter, error)

	//用db开启tx事务
	BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error)
	Begin() (Engine, error)

	Commit() error
	Rollback() error
}

type DBer interface {
	//用db开启tx事务
	BeginTx(ctx context.Context, opts *sql.TxOptions) TXer

	//lorm扩展方法
	Delete(v any) OrmTableDelete
}

type TXer interface {
	Commit() error
	Rollback() error

	//lorm扩展方法
	Delete(v any) OrmTableDelete
}

/*
*
直属lnDb

Dialecter 的实现有两种

	coreDb
	coreTx

内部属性为

	ldb      *sql.DB
	dialect Dialecter
*/
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
	initContext() Dialecter
	hasErr() bool
	getErr() error

	appendBaseToken(token baseToken)

	//对执行语句进行方言处理
	//toSqlInsert()

	//prepare(query string) string
	exec(query string, args ...any) (string, []any)
	execBatch(query string, args [][]any) (string, [][]any)

	query(query string, args ...any) (string, []any)
	queryBatch(query string) string

	parse(c Clause) (string, error)

	tableInsertGen()
	getSql() string
}
