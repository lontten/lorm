package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

type DbConfig interface {
	DriverName() string
}

type MysqlConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

type PoolConfig struct {
	maxIdleCount      int           // zero means defaultMaxIdleConns; negative means 0
	maxOpen           int           // <= 0 means unlimited
	maxLifetime       time.Duration // maximum amount of time a connection may be reused
	maxIdleTime       time.Duration // maximum amount of time a connection may be idle before being closed
	cleanerCh         chan struct{}
	waitCount         int64 // Total number of connections waited for.
	maxIdleClosed     int64 // Total number of connections closed due to idle count.
	maxIdleTimeClosed int64 // Total number of connections closed due to idle time.
	maxLifetimeClosed int64 // Total number of connections closed due to max connection lifetime limit.
}

func (c *MysqlConfig) DriverName() string {
	return MYSQL
}

type PgConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Other    string
}

func (c *PgConfig) DriverName() string {
	return POSTGRES
}

type DB struct {
	db        *sql.DB
	dbConfig  DbConfig
	ormConfig OrmConfig
}

func (db DB) OrmConfig() OrmConfig {
	return db.ormConfig
}


func (db *DB) SetOrmConfig(c OrmConfig) {
	db.ormConfig = c
}


type OrmConfig struct {
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
	LogicDeleteValue     func() interface{}
	LogicNotDeleteValue  func() interface{}

	//多租户 tenantIdFieldName不为零值，即开启
	TenantIdFieldName      string
	TenantIdValueFun       func() interface{}
	TenantIdIgnoreTableFun func(tableName string, dest interface{}) string
}

type Engine struct {
	db      *DB
	Base    EngineBase
	Extra   EngineExtra
	Table   EngineTable
	Classic EngineClassic
}

func open(c DbConfig, pc *PoolConfig) (dp *DB, err error) {
	if c == nil {
		fmt.Println("dbconfig canot be nil")
		panic(errors.New("dbconfig canot be nil"))
	}

	var db *sql.DB
	switch c.DriverName() {
	case MYSQL:
		c := c.(*MysqlConfig)
		dsn := c.User + ":" + c.Password +
			"@tcp(" + c.Host +
			":" + c.Port +
			")/" + c.DbName
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
	case POSTGRES:
		c := c.(*PgConfig)
		dsn := "user=" + c.User +
			" password=" + c.Password +
			" dbname=" + c.DbName +
			" host=" + c.Host +
			" port= " + c.Port
		if c.Other == "" {
			dsn += " sslmode=disable TimeZone=Asia/Shanghai"
		}
		dsn += c.Other
		db, err = sql.Open("pgx", dsn)
		if err != nil {
			panic(err)
		}
	default:
		return nil, errors.New("无此db 类型")
	}
	if pc != nil {
		db.SetConnMaxLifetime(pc.maxLifetime)
		db.SetConnMaxIdleTime(pc.maxIdleTime)
		db.SetMaxOpenConns(pc.maxOpen)
		db.SetMaxIdleConns(pc.maxIdleCount)
	}
	return &DB{
		db:        db,
		dbConfig:  c,
		ormConfig: OrmConfig{},
	}, nil
}

func MustConnect(c DbConfig, pc *PoolConfig) *DB {
	db, err := Connect(c, pc)
	if err != nil {
		panic(err)
	}
	return db
}

func Connect(c DbConfig, pc *PoolConfig) (*DB, error) {
	pool, err := open(c, pc)
	if err != nil {
		return nil, err
	}
	err = pool.db.Ping()
	if err != nil {
		return nil, err
	}
	return pool, err
}

func (db *DB) GetEngine() Engine {
	return Engine{
		db:      db,
		Base:    EngineBase{db: db, context: OrmContext{}},
		Extra:   EngineExtra{db: db, context: OrmContext{}},
		Classic: EngineClassic{db: db, context: OrmContext{}},
		Table: EngineTable{
			context:      OrmContext{},
			db:           db,
		},
	}
}

func (db DB) Exec(query string, args ...interface{}) (int64, error) {
	switch db.dbConfig.DriverName() {
	case MYSQL:
	case POSTGRES:
		var i = 1
		for {
			t := strings.Replace(query, " ? ", " $"+strconv.Itoa(i)+" ", 1)
			if t == query {
				break
			}
			i++
			query = t
		}
	default:
		return 0, errors.New("无此db drive 类型")
	}
	log.Println(query, args)

	exec, err := db.db.Exec(query, args...)
	if err != nil {
		return 0, err
	}
	return exec.RowsAffected()
}

func (db DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	switch db.dbConfig.DriverName() {
	case POSTGRES:
		var i = 1
		for {
			t := strings.Replace(query, " ? ", " $"+strconv.Itoa(i)+" ", 1)
			if t == query {
				break
			}
			i++
			query = t
		}
	default:
		return nil, errors.New("无此db drive 类型")
	}
	log.Println(query, args)

	return db.db.Query(query, args...)

}

//
//func (e *Engine) Begin() *Tx {
//	return &Tx{
//		Base:    e.Base,
//		Extra:   e.Extra,
//		Classic: e.Classic,
//		Table:   e.Table,
//	}
//}

type OrmContext struct {
	query  *strings.Builder
	args   []interface{}
	startd bool
}

type OrmSelect struct {
	db      DBer
	context OrmContext
}

type OrmFrom struct {
	db      DBer
	context OrmContext
}

type OrmWhere struct {
	db      DBer
	context OrmContext
}

func selectArgsArr2SqlStr(context OrmContext, args []string) {
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
	sb.WriteString(" ( ")
	for i, v := range args {
		if i == 0 {
			sb.WriteString(v)
		} else {
			sb.WriteString(" , " + v)
		}
	}
	sb.WriteString(" ) ")
	sb.WriteString(" VALUES ")
	sb.WriteString("( ")
	for i := range args {
		if i == 0 {
			sb.WriteString(" ? ")
		} else {
			sb.WriteString(", ? ")
		}
	}
	sb.WriteString(" ) ")
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
