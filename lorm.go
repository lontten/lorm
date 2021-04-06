package lorm

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"strings"
	"time"
)

type DbConfig interface {
	DriverName() string
}

type DbPoolConfig interface {
	PoolDriverName() string
}

type MysqlConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
}

type MysqlPoolConfig struct {
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
	return "mysql"
}

func (c *MysqlPoolConfig) PoolDriverName() string {
	return "mysql"
}

type PgConfig struct {
	Host     string
	Port     string
	DbName   string
	User     string
	Password string
	Other    string
}

type PgPoolConfig struct {

	// MaxConnLifetime is the duration since creation after which a connection will be automatically closed.
	MaxConnLifetime time.Duration

	// MaxConnIdleTime is the duration after which an idle connection will be automatically closed by the health check.
	MaxConnIdleTime time.Duration

	// MaxConns is the maximum size of the pool.
	MaxConns int32

	// MinConns is the minimum size of the pool. The health check will increase the number of connections to this
	// amount if it had dropped below.
	MinConns int32

	// HealthCheckPeriod is the duration between checks of the health of idle connections.
	HealthCheckPeriod time.Duration

	// If set to true, pool doesn't do any I/O operation on initialization.
	// And connects to the server only when the pool starts to be used.
	// The default is false.
	LazyConnect bool
}

func (c *PgConfig) DriverName() string {
	return "postgresql"
}

func (c *PgPoolConfig) PoolDriverName() string {
	return "postgresql"
}

type DbPool struct {
	context   *OrmContext
	db        interface{}
	dbConfig  DbConfig
	ormConfig OrmConfig
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

func (db *DbPool) Exec(query string, args ...interface{}) (int64, error) {
	log.Println(query, args)
	switch db.dbConfig.DriverName() {
	case "mysql":
		exec, err := db.db.(*sql.DB).Exec(query, args...)
		if err != nil {
			return 0, err
		}
		return exec.RowsAffected()
	case "postgresql":
		exec, err := db.db.(*pgxpool.Pool).Exec(context.Background(), query, args...)
		if err != nil {
			return 0, err
		}
		return exec.RowsAffected(), nil
	default:
		return 0, errors.New("无此db 类型")
	}

}

func open(c DbConfig, pc DbPoolConfig) (*DbPool, error) {
	switch c.DriverName() {
	case "mysql":
		c := c.(*MysqlConfig)
		pc := pc.(*MysqlPoolConfig)

		dsn := c.User + ":" + c.Password +
			"@tcp(" + c.Host +
			":" + c.Port +
			")/" + c.DbName
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			panic(err)
		}
		db.SetConnMaxLifetime(pc.maxLifetime)
		db.SetConnMaxIdleTime(pc.maxIdleTime)
		db.SetMaxOpenConns(pc.maxOpen)
		db.SetMaxIdleConns(pc.maxIdleCount)
		return &DbPool{
			db:        db,
			dbConfig:  c,
			ormConfig: OrmConfig{},
		}, nil
	case "postgresql":
		c := c.(*PgConfig)
		pc := pc.(*PgPoolConfig)

		dsn := "user=" + c.User +
			" password=" + c.Password +
			" dbname=" + c.DbName +
			" host=" + c.Host +
			" port= " + c.Port
		if c.Other == "" {
			dsn += " sslmode=disable TimeZone=Asia/Shanghai"
		}
		dsn += c.Other

		config, err := pgx.ParseConfig(dsn)
		if err != nil {
			return nil, err
		}

		pgpc := pgxpool.Config{
			ConnConfig:        config,
			MaxConnLifetime:   pc.MaxConnLifetime,
			MaxConnIdleTime:   pc.MaxConnIdleTime,
			MaxConns:          pc.MaxConns,
			MinConns:          pc.MinConns,
			HealthCheckPeriod: pc.HealthCheckPeriod,
			LazyConnect:       pc.LazyConnect,
		}
		pool, err := pgxpool.ConnectConfig(context.Background(), &pgpc)
		if err != nil {
			return nil, err
		}
		return &DbPool{
			db:        pool,
			dbConfig:  c,
			ormConfig: OrmConfig{},
		}, nil

	default:
		return nil, errors.New("无此db 类型")

	}
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
	db      *DbPool
	context *OrmContext
}

type OrmFrom struct {
	db      *DbPool
	context *OrmContext
}

type OrmWhere struct {
	db      *DbPool
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
