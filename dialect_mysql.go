package lsql

import (
	"database/sql"
	"errors"
	"github.com/lontten/lsql/utils"
	"strings"
)

type MysqlDialect struct {
	db  DBer
	log Logger
}

func (m MysqlDialect) Copy(db DBer) Dialect {
	return &MysqlDialect{
		db:  db,
		log: m.log,
	}
}

func (m MysqlDialect) query(query string, args ...interface{}) (*sql.Rows, error) {
	m.log.Println(query, args)
	return m.db.Query(query, args...)
}

func (m MysqlDialect) insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (int64, error) {
	cs := make([]string, 0)
	vs := make([]interface{}, 0)

	for i, column := range columns {
		if utils.Contains(fields, column) {
			continue
		}
		cs = append(cs, column)
		vs = append(vs, args[i])
	}

	var query = "INSERT INTO " + table + "(" + strings.Join(columns, ",") +
		") VALUES (" + strings.Repeat("?", len(args)) +
		") ON duplicate key UPDATE " + strings.Join(cs, "=?, ") + "=?"

	args = append(args, vs...)
	m.log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m MysqlDialect) insertOrUpdateByUnique(table string, fields []string, columns []string, args ...interface{}) (int64, error) {
	return 0, errors.New("MySQL insertOrUpdateByUnique not support, please use insertOrUpdateByPrimaryKey")
}

func (m MysqlDialect) queryBatch(query string) (*sql.Stmt, error) {
	return m.db.Prepare(query)
}

func (m MysqlDialect) exec(query string, args ...interface{}) (int64, error) {
	m.log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m MysqlDialect) execBatch(query string, args [][]interface{}) (int64, error) {
	m.log.Println(query, args)

	var num int64 = 0
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	for _, arg := range args {
		exec, err := stmt.Exec(arg...)
		m.log.Println(query, args)
		if err != nil {
			return num, err
		}
		rowsAffected, err := exec.RowsAffected()
		if err != nil {
			return num, err
		}
		num += rowsAffected
	}
	return num, nil
}
