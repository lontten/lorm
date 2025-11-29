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

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"
)

type coreDB struct {
	db      *sql.DB
	dialect Dialecter
}

func (db *coreDB) init() Engine {
	return &coreDB{
		db:      db.db,
		dialect: db.dialect.copyContext(),
	}
}

func (db *coreDB) ping() error {
	return db.db.Ping()
}

func (db *coreDB) getCtx() *ormContext {
	return db.dialect.getCtx()
}
func (db *coreDB) getDialect() Dialecter {
	return db.dialect
}
func (db *coreDB) query(query string, args ...any) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}
func (db *coreDB) exec(query string, args ...any) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *coreDB) prepare(query string) (Stmter, error) {
	stmt, err := db.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &coreDBStmt{
		db:      stmt,
		dialect: db.dialect,
	}, nil
}

func (db *coreDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error) {
	tx, err := db.db.BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &coreTX{
		tx:      tx,
		dialect: db.dialect,
	}, nil
}

func (db *coreDB) Begin() (Engine, error) {
	return db.BeginTx(context.Background(), nil)
}

func (db *coreDB) Commit() error {
	return errors.New("this not tx")
}

func (db *coreDB) Rollback() error {
	return errors.New("this not tx")
}

func (db *coreDB) ToWhereSQL(w *WhereBuilder, primaryKeyColumnNames ...string) (string, []any, error) {
	return w.toSql(db.getDialect().parse, primaryKeyColumnNames...)
}
