package lorm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/lontten/lorm/types"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"time"
)

var ImpValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
var ImpNuller = reflect.TypeOf((*types.NullEr)(nil)).Elem()

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

type DbConfig interface {
	DriverName() string
	Open() (*sql.DB, error)
	Dialect(db *sql.DB, log *log.Logger) Dialect
}

type PoolConf struct {
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed

	Logger *log.Logger
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

func (c *MysqlConf) Dialect(db *sql.DB, logger *log.Logger) Dialect {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}
	return MysqlDialect{db: db, log: Logger{log: logger}}
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

func (c *PgConf) Dialect(db *sql.DB, logger *log.Logger) Dialect {
	if logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
		log.SetFlags(log.LstdFlags | log.Llongfile)
	}
	return PgDialect{db: db, log: Logger{log: logger}}
}

func (c *PgConf) DriverName() string {
	return POSTGRES
}

func (c *PgConf) Open() (*sql.DB, error) {
	dsn := "user=" + c.User +
		" password=" + c.Password +
		" dbname=" + c.DbName +
		" host=" + c.Host +
		" port= " + c.Port +
		" "
	if c.Other == "" {
		dsn += "sslmode=disable TimeZone=Asia/Shanghai"
	}
	dsn += c.Other
	return sql.Open("pgx", dsn)
}

var _ormCtx = OrmContext{
	ormConf: OrmConf{
		PoDir:           "src/model/po",
		Author:          "lontten",
		IdType:          0,
		PrimaryKeyNames: []string{"id"},
	},
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
		ctx:      _ormCtx,
		dialect:  c.Dialect(db, pc.Logger),
	}, nil
}

func MustConnect(c DbConfig, pc *PoolConf) *DB {
	db, err := Connect(c, pc)
	if err != nil {
		panic(err)
	}
	return db
}

func MustConnectMock(db *sql.DB, c DbConfig) *DB {
	return &DB{
		db:       db,
		dbConfig: c,
		ctx:      _ormCtx,
		dialect:  c.Dialect(db, nil),
	}
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
