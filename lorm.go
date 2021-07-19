package lorm

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"time"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

type Dialect interface {
	ToDialectSql(sql string)string


}


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
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed
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
	db      DB
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
		db.SetConnMaxLifetime(pc.MaxLifetime)
		db.SetConnMaxIdleTime(pc.MaxIdleTime)
		db.SetMaxOpenConns(pc.MaxOpen)
		db.SetMaxIdleConns(pc.MaxIdleCount)
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

func (db DB) GetEngine(c *OrmConfig) Engine {
	if c == nil {
		config := OrmConfig{}
		c = &config
	}
	db.ormConfig = *c

	return Engine{
		db:      db,
		Base:    EngineBase{db: db, context: OrmContext{}},
		Extra:   EngineExtra{db: db, context: OrmContext{}},
		Classic: EngineClassic{db: db, context: OrmContext{}},
		Table: EngineTable{
			context: OrmContext{},
			db:      db,
		},
	}
}

type Ha struct {

}
