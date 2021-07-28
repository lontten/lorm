package lorm

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	MYSQL    = "mysql"
	POSTGRES = "postgres"
)

type Lorm interface {
	ScanLn(rows *sql.Rows, v interface{}) (int64, error)
	Scan(rows *sql.Rows, v interface{}) (int64, error)
}

func (c OrmConf) ScanLn(rows *sql.Rows, v interface{}) (num int64, err error) {
	defer rows.Close()
	value := reflect.ValueOf(v)
	code, base := basePtrStructBaseValue(value)
	if code == -1 {
		return 0, errors.New("dest need a  ptr")
	}
	if code == -2 {
		return 0, errors.New("need a ptr struct or base type")
	}

	num = 1
	t := base.Type()

	columns, err := rows.Columns()
	if err != nil {
		return
	}
	cfm, err := getColFieldIndexLinkMap(columns, t, c.FieldNamePrefix)
	if err != nil {
		return
	}
	if rows.Next() {
		box, _, v := createColBox(t, cfm)
		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return
		}
		base.Set(v)
	}

	if rows.Next() {
		return 0, errors.New("result to many for one")
	}
	return
}

func (c OrmConf) Scan(rows *sql.Rows, v interface{}) (int64, error) {
	defer rows.Close()
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Ptr {
		return 0, errors.New("need a ptr type")
	}
	arr := value.Elem()
	if arr.Kind() != reflect.Slice {
		return 0, errors.New("need a slice type")
	}

	slice := arr.Type()

	base := slice.Elem()
	isPtr := base.Kind() == reflect.Ptr
	code, base := baseStructBaseType(base)
	if code == -2 {
		return 0, errors.New("need a struct or base type in  slice")
	}

	columns, err := rows.Columns()
	if err != nil {
		return 0, err
	}
	cfm, err := getColFieldIndexLinkMap(columns, base, c.FieldNamePrefix)
	fmt.Println(len(cfm))
	fmt.Println("------")
	if err != nil {
		return 0, err
	}
	var num int64 = 0
	for rows.Next() {
		box, vp, v := createColBox(base, cfm)

		err = rows.Scan(box...)
		if err != nil {
			fmt.Println(err)
			return 0, err
		}
		if isPtr {
			arr.Set(reflect.Append(arr, vp))
		} else {
			arr.Set(reflect.Append(arr, v))
		}
		num++
	}
	return num, nil
}

type Dialect interface {
	DriverName() string
	ToDialectSql(sql string) string

}

type MysqlDialect struct {
	lormConf OrmConf
}

func (m MysqlDialect) DriverName() string {
	return MYSQL
}

func (m MysqlDialect) ToDialectSql(sql string) string {
	return sql
}

type PgDialect struct {
	lormConf OrmConf
}

func (m PgDialect) DriverName() string {
	return POSTGRES
}

func (m PgDialect) ToDialectSql(sql string) string {
	var i = 1
	for {
		t := strings.Replace(sql, " ? ", " $"+strconv.Itoa(i)+" ", 1)
		if t == sql {
			break
		}
		i++
		sql = t
	}
	return sql
}

type DbConfig interface {
	DriverName() string
	Open() (*sql.DB, error)
	Dialect(c OrmConf) Dialect
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

func (c *MysqlConf) Dialect(cf OrmConf) Dialect {
	return MysqlDialect{cf}
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

func (c *PgConf) Dialect(cf OrmConf) Dialect {
	return PgDialect{cf}
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

type OrmConf struct {
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
	lormConf OrmConf

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

func (db DB) GetEngine(c *OrmConf) Engine {
	conf := OrmConf{}
	if c != nil {
		conf = *c
	}
	return Engine{
		db:       db,
		lormConf: conf,
		Base: EngineBase{
			db:      db,
			lorm:    conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
		Extra: EngineExtra{
			db:      db,
			lormConf: conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
		Classic: EngineClassic{
			db:      db,
			lormConf: conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
		Table: EngineTable{
			db:      db,
			lormConf: conf,
			context: OrmContext{},
			dialect: db.dbConfig.Dialect(conf),
		},
	}
}

type Ha struct {
}
