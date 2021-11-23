package lorm

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"time"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)
var ormConfig=OrmConf{
	PoDir:                "src/model/po",
	IsFileOverride:       false,
	Author:               "lontten",
	IsActiveRecord:       false,
	IdType:               0,
	TableNamePrefix:      "",
	FieldNamePrefix:      "",
	PrimaryKeyNames:     []string{"id"},
	LogicDeleteYesSql:    "",
	LogicDeleteNoSql:     "",
	LogicDeleteSetSql:    "",
	TenantIdFieldName:    "",
	TenantIdValueFun:     nil,
	TenantIgnoreTableFun: nil,
}

type DbConfig interface {
	DriverName() string
	Open() (*sql.DB, error)
	Dialect(db *sql.DB) Dialect
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

func (c *MysqlConf) Dialect(db *sql.DB) Dialect {
	return MysqlDialect{db}
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

func (c *PgConf) Dialect(db *sql.DB) Dialect {
	return PgDialect{db}
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

type Engine struct {
	db      DB

	Base    EngineBase
	Extra   EngineExtra
	Table   EngineTable
	Classic EngineNative
}

type EngineEr interface {
	Db(c *OrmConf) Engine
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

func MustConnect(c DbConfig, pc *PoolConf) EngineEr {
	db, err := Connect(c, pc)
	if err != nil {
		panic(err)
	}
	return db
}

func Connect(c DbConfig, pc *PoolConf) (EngineEr, error) {
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
