package lorm

import (
	"database/sql"
	"reflect"
)


type DBer interface {
	exec(query string, args ...interface{}) (int64, error)
	query(query string, args ...interface{}) (*sql.Rows, error)
}


type Dialect interface {
	DriverName() string
	ToDialectSql(sql string) string

}



type OrmCore interface {
	ScanLn(rows *sql.Rows, v interface{}) (int64, error)
	Scan(rows *sql.Rows, v interface{}) (int64, error)

	//获取主键列表
	primaryKeys(tableName string, v reflect.Value)[]string
	//获取表名
	tableName(v reflect.Value)(string,error)
}





type Queryer interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Preparer interface {
	Prepare(query string) (*sql.Stmt, error)
}

