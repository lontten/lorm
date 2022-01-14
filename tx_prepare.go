package lorm

import (
	"database/sql"
)

type StmtTx struct {
	stmt *sql.Stmt
}

func (tx Tx) Prepare(query string) (s Stmt, err error) {
	stmt, err := tx.tx.Prepare(query)
	if err != nil {
		return
	}
	s.stmt = stmt
	return
}
