//  Copyright 2025 lontten lontten@163.com
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

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

	// 解析 WhereBuilder
	ToWhereSQL(w *WhereBuilder, primaryKeyColumnNames ...string) (string, []any, error)
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
	copyContext() Dialecter
	hasErr() bool
	getErr() error

	//对执行语句进行方言处理
	//toSqlInsert()

	//prepare(query string) string
	exec(query string, args ...any) (string, []any)
	execBatch(query string, args [][]any) (string, [][]any)

	query(query string, args ...any) (string, []any)
	queryBatch(query string) string

	parse(c Clause) (string, error)

	tableInsertGen()
	tableDelGen()
	tableUpdateGen()
	tableSelectGen()
	getSql(sql ...string)

	escapeIdentifier(s string) string // 转义 数据库保留字
}
