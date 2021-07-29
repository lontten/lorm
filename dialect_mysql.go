package lorm


type MysqlDialect struct {
	lormConf OrmConf
}

func (m MysqlDialect) DriverName() string {
	return MYSQL
}

func (m MysqlDialect) ToDialectSql(sql string) string {
	return sql
}
