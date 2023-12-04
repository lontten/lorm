package lorm

import (
	"context"
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

// 当scan的能力，struct的每个filed都要有这个能力，才行。
var ImpValuer = reflect.TypeOf((*driver.Valuer)(nil)).Elem()

// 接收 null的能力
var ImpNuller = reflect.TypeOf((*types.NullEr)(nil)).Elem()

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
		ormConf: OrmConf{
			PoDir:           "src/model/po",
			Author:          "lontten",
			IdType:          0,
			PrimaryKeyNames: []string{"id"},
		},
	}
}

func open(c DbConfig, pc *PoolConf) (dp *lnDB, err error) {
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
	return &lnDB{
		core: coreDb{
			db:      db,
			dialect: c.dialect(ctx),
		},
	}, nil
}

func MustConnect(c DbConfig, pc *PoolConf) DBer {
	db, err := Connect(c, pc)
	if err != nil {
		panic(err)
	}
	return db
}

func MustConnectMock(db *sql.DB, c DbConfig) DBer {
	ctx := genOrmCtx(nil)
	l := lnDB{
		core: coreDb{
			db:      db,
			dialect: c.dialect(ctx),
		},
	}
	return l
}

func Connect(c DbConfig, pc *PoolConf) (DBer, error) {
	db, err := open(c, pc)
	if err != nil {
		return nil, err
	}

	err = db.core.getDB().Ping()
	if err != nil {
		return nil, err
	}
	return db, err
}

type lnDB struct {
	core corer
}

func (db lnDB) getErr() error {
	return db.core.getErr()
}

func (db lnDB) hasErr() bool {
	return db.core.hasErr()
}

func (db lnDB) Rollback() error {
	return db.core.doRollback()
}

func (db lnDB) Commit() error {
	return db.core.doCommit()
}

func (db lnDB) BeginTx(ctx context.Context, opts *sql.TxOptions) TXer {
	tx := db.core.doBeginTx(ctx, opts)
	return lnDB{
		core: tx,
	}
}

func (db lnDB) Query(query string, args ...interface{}) *NativeQuery {
	return db.core.query(query, args...)
}
func (db lnDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	return db.core.doExec(query, args...)
}

//todo 下面未重构--------------
