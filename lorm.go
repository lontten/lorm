package lorm

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"reflect"
	"time"
)

var ImpValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()
var ImpScanner = reflect.TypeOf((*sql.Scanner)(nil)).Elem()

type PoolConf struct {
	MaxIdleCount int           // zero means defaultMaxIdleConns; negative means 0
	MaxOpen      int           // <= 0 means unlimited
	MaxLifetime  time.Duration // maximum amount of time a connection may be reused
	MaxIdleTime  time.Duration // maximum amount of time a connection may be idle before being closed

	Logger *log.Logger
}

func genOrmCtx(pc *PoolConf) *ormContext {
	var logger *log.Logger
	if pc == nil || pc.Logger == nil {
		logger = log.New(os.Stdout, "", log.LstdFlags)
		log.SetFlags(log.LstdFlags | log.Llongfile)
	} else {
		logger = pc.Logger
	}
	return &ormContext{
		log: Logger{log: logger},
		ormConf: &OrmConf{
			PoDir:           "src/model/po",
			Author:          "lontten",
			IdType:          0,
			PrimaryKeyNames: []string{"id"},
		},
	}
}

func open(c DbConfig, pc *PoolConf) (Engine, error) {
	if c == nil {
		fmt.Println("dbconfig cannot be nil")
		return nil, errors.New("dbconfig cannot be nil")
	}

	db, err := c.open()
	if err != nil {
		return nil, err
	}

	if pc != nil {
		db.SetConnMaxLifetime(pc.MaxLifetime)
		db.SetConnMaxIdleTime(pc.MaxIdleTime)
		db.SetMaxOpenConns(pc.MaxOpen)
		db.SetMaxIdleConns(pc.MaxIdleCount)
	}
	ctx := genOrmCtx(pc)
	return &DB{dialect: c.dialect(ctx, &coreDB{db: db})}, nil
}

func MustConnect(c DbConfig, pc *PoolConf) Engine {
	db, err := Connect(c, pc)
	if err != nil {
		panic(err)
	}
	return db
}

func MustConnectMock(db *sql.DB, c DbConfig) Engine {
	ctx := genOrmCtx(nil)
	return &DB{dialect: c.dialect(ctx, &coreDB{db: db})}
}

func Connect(c DbConfig, pc *PoolConf) (Engine, error) {
	db, err := open(c, pc)
	if err != nil {
		return nil, err
	}

	err = db.getDialect().getDB().ping()
	if err != nil {
		return nil, err
	}
	return db, err
}
