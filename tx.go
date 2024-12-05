package lorm

import (
	"context"
	"database/sql"
	"github.com/pkg/errors"
)

// TX -----------------TX---------------------
type TX struct {
	dialect Dialecter
}

func (db *TX) getDialect() Dialecter {
	return db.dialect
}
func (db *TX) prepare(query string) (EngineStmt, error) {
	dialect := db.getDialect()
	err := dialect.prepare(query)
	if err != nil {
		return nil, err
	}
	return &TXStmt{dialect: dialect}, nil
}

func (db *TX) BeginTx(ctx context.Context, opts *sql.TxOptions) (Engine, error) {
	return nil, errors.New("this is tx")
}

func (db *TX) Begin() (Engine, error) {
	return nil, errors.New("this is tx")
}

func (db *TX) Commit() error {
	return db.getDialect().commit()
}

func (db *TX) Rollback() error {
	return db.getDialect().rollback()
}

// -----------------DB-end---------------------

// coreTX -----------------coreTX---------------------

type coreTX struct {
	tx *sql.Tx
}

func (db *coreTX) ping() error {
	return errors.New("this is tx")
}

func (db *coreTX) prepare(query string) (Stmter, error) {
	stmt, err := db.tx.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &coreTXStmt{tx: stmt}, nil
}

func (db *coreTX) query(query string, args ...any) (*sql.Rows, error) {
	return db.tx.Query(query, args...)
}
func (db *coreTX) exec(query string, args ...any) (sql.Result, error) {
	return db.tx.Exec(query, args...)
}

func (db *coreTX) beginTx(ctx context.Context, opts *sql.TxOptions) (DBer, error) {
	return nil, errors.New("this is tx")
}

func (db *coreTX) commit() error {
	return db.tx.Commit()
}

func (db *coreTX) rollback() error {
	return db.tx.Rollback()
}

// -----------------coreTX-end---------------------
