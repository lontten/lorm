package lorm

import "database/sql"

type MysqlDialect struct {
	db *sql.DB
}

func (m MysqlDialect) DriverName() string {
	return MYSQL
}

func (m MysqlDialect) ToDialectSql(sql string) string {
	return sql
}
