package lorm

import (
	"strconv"
	"strings"
)

type PgDialect struct {
	lormConf OrmConf
}

func (m PgDialect) DriverName() string {
	return POSTGRES
}

func (m PgDialect) ToDialectSql(sql string) string {
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

