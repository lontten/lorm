package lorm

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"time"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

type Dialect interface {
	DriverName() string
	ToDialectSql(sql string) string
}

type MysqlDialect struct {
	lormConf LormConf
}

type PgDialect struct {
	lormConf LormConf
}

type DbConfig interface {
	DriverName() string
	Open() (*sql.DB, error)
}

type PoolConf struct {
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed
}

type MysqlConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

func (c *MysqlConf) DriverName() string {
	return MYSQL
}

func (c *MysqlConf) Open() (*sql.DB, error) {
	dsn := c.User + ":" + c.Password +
		"@tcp(" + c.Host +
		":" + c.Port +
		")/" + c.DbName
	return sql.Open("mysql", dsn)
}

type PgConf struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Other    string
}

func (c *PgConf) DriverName() string {
	return POSTGRES
}

func (c *PgConf) Open() (*sql.DB, error) {
	dsn := "user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DbName +
		" host=" + c.Host +
		" port= " + c.Port
	if c.Other == "" {
		dsn += " sslmode=disable TimeZone=Asia/Shanghai"
	}
	dsn += c.Other
	return sql.Open("pgx", dsn)
}

type LormConf struct {
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
	//TableNameFun >  tag > TableNamePrefix
	TableNamePrefix string
	TableNameFun    func(structName string, dest interface{}) string

	//字段名
	FieldNamePrefix string

	//主键 默认为id
	PrimaryKeyNames   []string
	PrimaryKeyNameFun func(tableName string, base reflect.Value) []string

	//逻辑删除 logicDeleteFieldName不为零值，即开启
	// LogicDeleteYesSql   lg.deleted_at is null
	// LogicDeleteNoSql   lg.deleted_at is not null
	// LogicDeleteSetSql   lg.deleted_at = now()
	LogicDeleteYesSql string
	LogicDeleteNoSql  string
	LogicDeleteSetSql string

	//多租户 tenantIdFieldName不为零值，即开启
	TenantIdFieldName      string
	TenantIdValueFun       func() interface{}
	TenantIdIgnoreTableFun func(structName string, dest interface{}) string
}

type Engine struct {
	db       DB
	lormConf LormConf

	Base    EngineBase
	Extra   EngineExtra
	Table   EngineTable
	Classic EngineClassic
}

func open(c DbConfig, pc *PoolConf) (dp *DB, err error) {
	if c == nil {
		fmt.Println("dbconfig canot be nil")
		return nil, errors.New("dbconfig canot be nil")
	}

	db, err := c.Open()
	if err != nil {
		return nil, err
	}

	if pc != nil {
		db.SetConnMaxLifetime(pc.MaxLifetime)
		db.SetConnMaxIdleTime(pc.MaxIdleTime)
		db.SetMaxOpenConns(pc.MaxOpen)
		db.SetMaxIdleConns(pc.MaxIdleCount)
	}
	return &DB{
		db:       db,
		dbConfig: c,
	}, nil
}

func MustConnect(c DbConfig, pc *PoolConf) *DB {
	db, err := Connect(c, pc)
	if err != nil {
		panic(err)
	}
	return db
}

func Connect(c DbConfig, pc *PoolConf) (*DB, error) {
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

func (db DB) GetEngine(c *LormConf) Engine {
	if c == nil {
		c = &LormConf{}
	}
	return Engine{
		db:       db,
		lormConf: *c,
		Base:     EngineBase{db: db, context: OrmContext{}},
		Extra:    EngineExtra{db: db, context: OrmContext{}},
		Classic:  EngineClassic{db: db, context: OrmContext{}},
		Table: EngineTable{
			context: OrmContext{},
			db:      db,
		},
	}
}

type Ha struct {
}
