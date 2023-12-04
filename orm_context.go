package lorm

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type tableSqlType int

// d前缀是单表的意思，tableSqlType 只用于单表操作
const (
	dInsert tableSqlType = iota
	dInsertOrUpdate
	dUpdate
	dDelete
	dSelect
	dGetOrInsert
	dHas
	dCount
)

type ormContext struct {
	ormConf  OrmConf
	dbConfig DbConfig

	log Logger
	err error

	tableSqlType tableSqlType //单表，sql类型crud

	baseTokens []baseToken

	isLgDel bool //是否启用了逻辑删除
	isTen   bool //是否启用了多租户

	//主键名-列表,这里考虑到多主键
	primaryKeyNames []string
	//主键值-列表
	primaryKeyValues [][]interface{}

	//主键名-列表,这里考虑到多主键-排除
	filterPrimaryKeyNames []string
	//主键值-列表-排除
	filterPrimaryKeyValues [][]interface{}

	//当前表名
	tableName string

	//字段列表-not nil
	columns []string
	//值列表-多个-not nil
	columnValues []interface{}

	//字段列表-nil
	nilColumns []string
	//值列表-多个-nil
	nilColumnValues []interface{}

	//字段列表-all
	allColumns []string
	//值列表-多个-all
	allColumnValues []interface{}

	//-------------------target---------------------
	//当前struct对象
	scanDest  interface{}
	destIsPtr bool
	//去除 ptr
	destValue reflect.Value
	//用作 参数合法行校验
	destBaseValue reflect.Value
	destBaseType  reflect.Type

	//------------------scan----------------------
	//scan base type
	scanDestBaseType reflect.Type
	//scan 是comp，false是single
	scanDestBaseTypeIsComp bool
	//scan 接收返回
	scanIsSlice bool
	//scan 为slice时，里面item是否是ptr
	scanSliceItemIsPtr bool

	//要执行的sql语句
	query *strings.Builder
	//参数
	args []interface{}

	started bool
}

// v03
// 初始化 表名
func (ctx *ormContext) initTableName() {
	if ctx.err != nil {
		return
	}
	if ctx.tableName != "" {
		ctx.err = errors.New("表名已经存在，不可再次初始化")
		return
	}
	tableName, err := ctx.ormConf.tableName(ctx.destBaseType)
	if err != nil {
		ctx.err = err
		return
	}
	ctx.tableName = tableName
}

// 获取struct对应的字段名 和 其值，
// slice为全部，一个为非nil字段。
func (ctx *ormContext) initColumnsValue() {
	if ctx.err != nil {
		return
	}
	cv, err := ctx.ormConf.getCompCV(ctx.destValue)
	if err != nil {
		ctx.err = err
		return
	}
	ctx.columns = cv.columns
	ctx.columnValues = cv.columnValues

	ctx.nilColumns = cv.nilColumns
	ctx.nilColumnValues = cv.nilColumnValues

	ctx.allColumns = cv.allColumns
	ctx.allColumnValues = cv.allColumnValues
	return
}

//todo 下面未重构--------------

func (ctx ormContext) Copy() ormContext {
	return ormContext{
		ormConf: ctx.ormConf,
		log:     ctx.log,
	}
}

// select 生成
func (ctx *ormContext) selectArgsArr2SqlStr(args []string) {
	query := ctx.query
	if ctx.started {
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
			ctx.started = true
		}
	}
}

// create 生成
func (ctx *ormContext) tableInsertGen() string {
	args := ctx.columns
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString(ctx.tableName + " ")

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

// v03
// 单表sql生成，insert
func (p *PgDialect) tGenInsert() string {
	args := p.ctx.columns
	var sb strings.Builder

	sb.WriteString("INSERT INTO ")
	sb.WriteString(p.ctx.tableName + " ")

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

func (ctx *ormContext) createSqlGenera(args []string) string {
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

// upd 生成
func (ctx *ormContext) tableUpdateArgs2SqlStr(args []string) string {
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

func (ctx *ormContext) initPrimaryKeyValues(v []interface{}) {
	if ctx.err != nil {
		return
	}

	idLen := len(v)
	if idLen == 0 {
		ctx.err = errors.New("ByPrimaryKey arg len num 0")
		return
	}
	pkLen := len(ctx.primaryKeyNames)

	idValuess := make([][]interface{}, 0)

	if pkLen == 1 { //单主键
		for _, i := range v {
			value := reflect.ValueOf(i)
			_, value, err := basePtrDeepValue(value)
			if err != nil {
				ctx.err = err
				return
			}

			if !isValuerType(value.Type()) {
				ctx.err = errors.New("ByPrimaryKey typ err,not single")
				return
			}

			idValues := make([]interface{}, 1)
			idValues[0] = value.Interface()
			idValuess = append(idValuess, idValues)
		}

	} else {
		for _, i := range v {
			value := reflect.ValueOf(i)
			_, value, err := basePtrDeepValue(value)
			if err != nil {
				ctx.err = err
				return
			}
			if !isCompType(value.Type()) {
				ctx.err = errors.New("ByPrimaryKey typ err,not comp")
				return
			}

			columns, values, err := getCompValueCV(value, ctx.ormConf)
			if err != nil {
				ctx.err = err
				return
			}
			if len(columns) != pkLen {
				ctx.err = errors.New("复合主键，filed数量 len err")
				return
			}

			idValues := make([]interface{}, 0)
			idValues = append(idValues, values...)
			idValuess = append(idValuess, idValues)
		}
	}

	ctx.primaryKeyValues = idValuess
}

func (ctx *ormContext) initSelfPrimaryKeyValues() {
	if ctx.err != nil {
		return
	}

	keyNum := len(ctx.primaryKeyNames)
	idValues := make([]interface{}, 0)
	columns, values, err := getCompCV(ctx.scanDest, ctx.ormConf)
	if err != nil {
		ctx.err = err
		return
	}
	//只要主键字段
	for _, key := range ctx.primaryKeyNames {
		for i, c := range columns {
			if c == key {
				idValues = append(idValues, values[i])
				continue
			}
		}
	}
	idLen := len(idValues)
	if idLen == 0 {
		ctx.err = errors.New("no pk")
		return
	}
	if keyNum != idLen {
		ctx.err = errors.New("comp pk num err")
		return
	}

	ctx.primaryKeyValues = append(ctx.primaryKeyValues, idValues)
}

// 生成select sql
func (ctx *ormContext) genSelectByPrimaryKey() []byte {
	tableName := ctx.tableName
	columns := ctx.columns
	selSql := ctx.ormConf.genSelectSqlCommon(tableName, columns)
	where := ctx.genWhereByPrimaryKey()
	return append(selSql, where...)
}

// 生成del sql
func (ctx *ormContext) genDelByPrimaryKey() []byte {
	return ctx.ormConf.genDelSqlCommon(ctx.tableName, ctx.primaryKeyNames)
}

// 生成del sql
func (ctx *ormContext) genDelByKeys(keys []string) []byte {
	return ctx.ormConf.genDelSqlCommon(ctx.tableName, keys)
}

// 生成del sql
func (ctx *ormContext) genDelByWhere(where []byte) []byte {
	return ctx.ormConf.genDelSqlByWhere(ctx.tableName, where)
}

// 生成where sql
func (ctx *ormContext) genWhereByPrimaryKey() []byte {
	keys := ctx.primaryKeyNames
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ctx.ormConf.TenantIdFieldName != "" && !ctx.ormConf.TenantIgnoreTableFun(tableName)
	return ctx.ormConf.GenWhere(keys, hasTen)
}

// 生成where sql
func (ctx *ormContext) genWhere(keys []string) []byte {
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ctx.ormConf.TenantIdFieldName != "" && !ctx.ormConf.TenantIgnoreTableFun(tableName)
	return ctx.ormConf.GenWhere(keys, hasTen)
}

// 为where语句附加上，租户，逻辑删除等。。。
func (ctx *ormContext) whereExtra(where []byte) []byte {
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ctx.ormConf.TenantIdFieldName != "" && !ctx.ormConf.TenantIgnoreTableFun(tableName)
	return ctx.ormConf.whereExtra(where, hasTen)
}
