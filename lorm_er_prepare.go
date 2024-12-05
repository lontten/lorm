package lorm

import (
	"database/sql"
)

type Stmt struct {
	stmt *sql.Stmt
	ctx  ormContext
}

func Prepare(db Engine, query string) (EngineStmt, error) {
	return db.prepare(query)
}

type NativePrepare struct {
	db   EngineStmt
	args []any
}

func (p *NativePrepare) ScanOne(dest any) (int64, error) {
	stmt := p.db
	dialect := stmt.getDialect()
	ctx := dialect.getCtx()
	ctx.initScanDestOne(dest)
	if ctx.err != nil {
		return 0, ctx.err
	}
	rows, err := p.db.query(p.args...)
	if err != nil {
		return 0, err
	}
	return ctx.ScanLnT(rows)
	return 0, nil
}

func (p *NativePrepare) ScanList(dest any) (int64, error) {
	//ctx := p.db.getCtx()
	//ctx.initScanDestList(dest)
	//if ctx.err != nil {
	//	return 0, ctx.err
	//}
	//rows, err := p.db.query(p.args...)
	//if err != nil {
	//	return 0, err
	//}
	//return ctx.ScanLnT(rows)
	return 0, nil
}
