package lorm

import "database/sql"

type MysqlDialect struct {
	db *sql.DB
}

func (m MysqlDialect) exec(query string, args ...interface{}) (int64, error) {
	Log.Println(query, args)

	exec, err := m.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (m MysqlDialect) query(query string, args ...interface{}) (*sql.Rows, error) {
	Log.Println("sql",query,args)
	return m.db.Query(query, args...)
}

func (m MysqlDialect) DriverName() string {
	return MYSQL
}

