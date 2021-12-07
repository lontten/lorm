package lorm

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type OrmContext struct {
	isLgDel bool
	isTen   bool

	//主键名-列表
	primaryKeyNames []string
	//主键值-列表
	primaryKeyValues [][]interface{}

	//当前表名
	tableName string

	//字段列表
	columns []string
	//值列表-多个
	columnValues []interface{}

	//-------------------target---------------------
	//当前struct对象
	scanDest interface{}
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
	err     error

	//log的层级
	log int
}

//select 生成
func (ctx *OrmContext) selectArgsArr2SqlStr(args []string) {
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
func (ctx *OrmContext) tableCreateGen() string {
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

func (ctx *OrmContext) createSqlGenera(args []string) string {
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
func (ctx *OrmContext) tableUpdateArgs2SqlStr(args []string) string {
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

func (ctx *OrmContext) initPrimaryKeyValues(v []interface{}) {
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

			if !isSingleType(value.Type()) {
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

			columns, values, err := getCompValueCV(value)
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

func (ctx *OrmContext) initSelfPrimaryKeyValues() {
	if ctx.err != nil {
		return
	}

	keyNum := len(ctx.primaryKeyNames)
	idValues := make([]interface{}, 0)
	columns, values, err := getCompCV(ctx.scanDest)
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

//生成select sql
func (ctx *OrmContext) genSelectByPrimaryKey() []byte {
	tableName := ctx.tableName
	columns := ctx.columns
	selSql := ormConfig.genSelectSqlCommon(tableName, columns)
	where := ctx.genWhereByPrimaryKey()
	return append(selSql, where...)
}

//生成del sql
func (ctx *OrmContext) genDelByPrimaryKey() []byte {
	return ormConfig.genDelSqlCommon(ctx.tableName, ctx.primaryKeyNames)
}

//生成del sql
func (ctx *OrmContext) genDelByKeys(keys []string) []byte {
	return ormConfig.genDelSqlCommon(ctx.tableName, keys)
}

//生成del sql
func (ctx *OrmContext) genDelByWhere(where []byte) []byte {
	return ormConfig.genDelSqlByWhere(ctx.tableName, where)
}

//生成where sql
func (ctx *OrmContext) genWhereByPrimaryKey() []byte {
	keys := ctx.primaryKeyNames
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName)
	return ormConfig.GenWhere(keys, hasTen)
}

//生成where sql
func (ctx *OrmContext) genWhere(keys []string) []byte {
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName)
	return ormConfig.GenWhere(keys, hasTen)
}

//为where语句附加上，租户，逻辑删除等。。。
func (ctx *OrmContext) whereExtra(where []byte) []byte {
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName)
	return ormConfig.whereExtra(where, hasTen)
}
