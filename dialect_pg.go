package lorm

import (
	"database/sql"
	"strconv"
	"strings"
)

type PgDialect struct {
	db *sql.DB
}

func (m PgDialect) DriverName() string {
	return POSTGRES
}

func (m PgDialect) query(query string, args ...interface{}) (*sql.Rows, error) {
	query = toPgSql(query)
	Log.Println("sql", query, args)
	return m.db.Query(query, args...)
}

func (m PgDialect) exec(query string, args ...interface{}) (int64, error) {
	query = toPgSql(query)
	Log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func toPgSql(sql string) string {
	var i = 1
	for {
		t := strings.Replace(sql, " ? ", " $"+strconv.Itoa(i)+" ", 1)
		if t == sql {
			break
		}
		i++
		sql = t
	}
	return sql
}
