package lorm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

type coreTX struct {
	tx      *sql.Tx
	dialect Dialecter
}

func (db *coreTX) init() Engine {
	return &coreTX{
		tx:      db.tx,
		dialect: db.dialect.initContext(),
	}
}
func (db *coreTX) ping() error {
	return errors.New("this is tx")
}
func (db *coreTX) getCtx() *ormContext {
	return db.dialect.getCtx()
}
func (db *coreTX) getDialect() Dialecter {
	return db.dialect
}
func (db *coreTX) query(query string, args ...any) (*sql.Rows, error) {
	return db.tx.Query(query, args...)
}
func (db *coreTX) exec(query string, args ...any) (sql.Result, error) {
	return db.tx.Exec(query, args...)
}

func (db *coreTX) prepare(query string) (Stmter, error) {
	stmt, err := db.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &coreDBStmt{
		db:      stmt,
		dialect: db.dialect,
	}, nil
}

func (db *coreTX) BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error) {
	return nil, errors.New("this is tx")
}

func (db *coreTX) Begin() (Engine, error) {
	return nil, errors.New("this is tx")
}

func (db *coreTX) Commit() error {
	return db.tx.Commit()
}

func (db *coreTX) Rollback() error {
	return db.tx.Rollback()
}
