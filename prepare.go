package lorm

import (
	"database/sql"
	"github.com/pkg/errors"
)

type Stmt struct {
	stmt *sql.Stmt

	ctx OrmContext
}

func (db DB) Prepare(query string) (Stmt, error) {
	return db.dialect.prepare(query)
}

func (s *Stmt) Exec(args ...interface{}) (int64, error) {
	exec, err := s.stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (s Stmt) Query(args ...interface{}) Prepare {
	return Prepare{
		db:   s,
		args: args,
	}
}

type Prepare struct {
	db   Stmt
	args []interface{}
}

func (p Prepare) ScanOne(dest interface{}) (int64, error) {
	if err := p.db.ctx.err; err != nil {
		return 0, err
	}
	p.db.ctx.initScanDestOne(dest)
	if p.db.ctx.scanIsSlice {
		return 0, errors.New("not support GetOne for slice")
	}
	p.db.ctx.checkScanDestField()
	if err := p.db.ctx.err; err != nil {
		return 0, err
	}

	rows, err := p.db.stmt.Query(p.args...)
	if err != nil {
		return 0, err
	}
	return p.db.ctx.ScanLn(rows)
}

func (p Prepare) ScanList(dest interface{}) (int64, error) {
	if err := p.db.ctx.err; err != nil {
		return 0, err
	}
	p.db.ctx.initScanDestList(dest)
	p.db.ctx.checkScanDestField()
	if err := p.db.ctx.err; err != nil {
		return 0, err
	}

	rows, err := p.db.stmt.Query(p.args...)
	if err != nil {
		return 0, err
	}
	return p.db.ctx.ScanLn(rows)
}
