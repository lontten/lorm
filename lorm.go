package lorm

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

type DB struct {
	context   *OrmContext
	db        *sql.DB
	ormConfig OrmConfig
}
type OrmConfig struct {
	//驱动名称
	DriverName string

	DbConfig string

	//po生成文件目录
	PoDir string
	//是否覆盖，默认true
	IsFileOverride bool

	//作者
	Author string
	//是否开启ActiveRecord模式,默认false
	IsActiveRecord bool

	IdType int

	//表名
	TableNamePrefix string
	TableNameFun    func(tableName string, dest interface{}) string

	//字段名
	FieldNamePrefix string

	//主键 默认为id
	IdName    string
	IdNameFun func(tableName string, dest interface{}) string

	//逻辑删除 logicDeleteFieldName不为零值，即开启
	LogicDeleteFieldName string
	LogicDeleteValue     interface{}
	LogicNotDeleteValue  interface{}

	//多租户 tenantIdFieldName不为零值，即开启
	TenantIdFieldName      string
	TenantIdValueFun       func() interface{}
	TenantIdIgnoreTableFun func(tableName string, dest interface{}) string
}

type Engine struct {
	Base    *EngineBase
	Extra   *EngineExtra
	Table   *EngineTable
	Classic *EngineClassic
}

func (db *DB) DriverName() string {
	return db.ormConfig.DriverName
}

func (db *DB) Exec(query string, args ...interface{}) (int64, error) {
	log.Println(query, args)
	exec, err := db.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &DB{db: db,
		context: &OrmContext{
			query:  &strings.Builder{},
			startd: false,
		},
	}, nil

}

func MustConnect(config OrmConfig) *Engine {
	db, err := Connect(config)
	if err != nil {
		panic(err)
	}
	return db
}

func Connect(config OrmConfig) (*Engine, error) {
	db, err := open(config.DriverName, config.DbConfig)
	if err != nil {
		return nil, err
	}

	err = db.db.Ping()
	if err != nil {
		return nil, err
	}

	db.ormConfig = config

	return &Engine{
		Base:    &EngineBase{db},
		Extra:   &EngineExtra{db},
		Classic: &EngineClassic{db},
		Table: &EngineTable{
			db:           db,
			idName:       "",
			tableName:    "",
			dest:         nil,
			columns:      nil,
			columnValues: nil,
		},
	}, nil

}

func (e *Engine) Begin() *Tx {
	return &Tx{
		Base:    e.Base,
		Extra:   e.Extra,
		Classic: e.Classic,
		Table:   e.Table,
	}
}

type OrmContext struct {
	query  *strings.Builder
	args   []interface{}
	startd bool
}

type OrmSelect struct {
	db      *DB
	context *OrmContext
}

type OrmFrom struct {
	db      *DB
	context *OrmContext
}

type OrmWhere struct {
	db      *DB
	context *OrmContext
}

func selectArgsArr2SqlStr(context *OrmContext, args []string) {
	query := context.query
	if context.startd {
		for _, name := range args {
			query.WriteString(", " + name)
		}
	} else {
		query.WriteString("SELECT ")
		for i := range args {
			if i == 0 {
				query.WriteString(args[i])
			} else {
				query.WriteString(", " + args[i])
			}
		}
		if len(args) > 0 {
			context.startd = true
		}
	}
}

func tableWhereArgs2SqlStr(args []string) string {
	var sb strings.Builder
	for i, where := range args {
		if i == 0 {
			sb.WriteString(" WHERE ")
			sb.WriteString(where)
			sb.WriteString(" = ? ")
			continue
		}
		sb.WriteString(" AND ")
		sb.WriteString(where)
		sb.WriteString(" = ? ")
	}
	return sb.String()
}

func tableSelectArgs2SqlStr(args []string) string {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, column := range args {
		if i == 0 {
			sb.WriteString(column)
		} else {
			sb.WriteString(" , ")
			sb.WriteString(column)
		}
	}
	return sb.String()
}

func tableCreateArgs2SqlStr(args []string) string {
	var sb strings.Builder
	sb.WriteString("( ")
	for i, v := range args {
		if i == 0 {
			sb.WriteString(v)
		} else {
			sb.WriteString(", " + v)
		}
	}
	sb.WriteString(" )")
	sb.WriteString(" VALUES ")
	sb.WriteString("( ")
	for i := range args {
		if i == 0 {
			sb.WriteString(" ? ")
		} else {
			sb.WriteString(", ?")
		}
	}
	sb.WriteString(" )")
	return sb.String()
}

func tableUpdateArgs2SqlStr(args []string) string {
	var sb strings.Builder
	l := len(args)
	for i, v := range args {
		if i != l-1 {
			sb.WriteString(v + " = ? ,")
		} else {
			sb.WriteString(v + " = ? ")
		}
	}
	return sb.String()
}
