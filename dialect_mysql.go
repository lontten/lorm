package lorm

import (
	"database/sql"
	"strings"
)

type MysqlDialect struct {
	db *sql.DB
}

func (m MysqlDialect) DriverName() string {
	return MYSQL
}

func (m MysqlDialect) query(query string, args ...interface{}) (*sql.Rows, error) {
	Log.Println(query, args)
	return m.db.Query(query, args...)
}

func (m MysqlDialect) insertOrUpdateByPrimaryKey(table string, fields []string, columns []string, args ...interface{}) (int64, error) {
	var query = "INSERT INTO " + table + "(" + strings.Join(columns, ",") +
		") VALUES (" + strings.Repeat("?", len(args)) +
		") ON duplicate key  UPDATE  "

	Log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m MysqlDialect) insertOrUpdateByFields(table string, fields []string, columns []string, args ...interface{}) (int64, error) {
	//todo
	Log.Println(table, args)

	exec, err := m.db.Exec(table, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m MysqlDialect) queryBatch(query string) (*sql.Stmt, error) {
	return m.db.Prepare(query)
}

func (m MysqlDialect) exec(query string, args ...interface{}) (int64, error) {
	Log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m MysqlDialect) execBatch(query string, args [][]interface{}) (int64, error) {
	Log.Println(query, args)

	var num int64 = 0
	stmt, err := m.db.Prepare(query)
	if err != nil {
		return 0, err
	}
	for _, arg := range args {
		exec, err := stmt.Exec(arg...)
		Log.Println(query, args)
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
