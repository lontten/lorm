package lorm

import (
	"database/sql"
	"reflect"
)

type DBer interface {
}

type Dialect interface {
	DriverName() string

	exec(query string, args ...interface{}) (int64, error)
	query(query string, args ...interface{}) (*sql.Rows, error)
}

type SqlUtil interface {
	selectArgsArr2SqlStr(args []string)
	tableWhereArgs2SqlStr(args []string, c OrmConf) string
	tableCreateArgs2SqlStr(args []string) string
	tableUpdateArgs2SqlStr(args []string) string
	tableWherePrimaryKey2SqlStr(ids []string, c OrmConf) string
}

type OrmCore interface {
	ScanLn(rows *sql.Rows, v interface{}) (int64, error)

	Scan(rows *sql.Rows, v interface{}) (int64, error)

	//获取主键列表
	primaryKeys(tableName string, v reflect.Value) []string
	//获取表名
	tableName(v reflect.Value) (string, error)

	initColumns(v reflect.Value) (c []string, err error)
	initColumnsValue(v reflect.Value) ([]string, []interface{}, error)
	getStructMappingColumns(t reflect.Type) (map[string]int, error)
	getStructMappingColumnsValueNotNil(v reflect.Value) (columns []string, values []interface{}, err error)
	getStructMappingColumnsValueList(v []reflect.Value) (m ModelParam, err error)

	getColFieldIndexLinkMap(columns []string, typ reflect.Type) (ColFieldIndexLinkMap, error)
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
