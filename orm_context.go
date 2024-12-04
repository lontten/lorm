package lorm

import (
	"github.com/lontten/lorm/field"
	"github.com/lontten/lorm/insert-type"
	"github.com/lontten/lorm/return-type"
	"github.com/lontten/lorm/soft-delete"
	"github.com/lontten/lorm/sql-type"
	"github.com/lontten/lorm/utils"
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
	ormConf *OrmConf
	extra   *ExtraContext

	// model 参数，用于校验字段类型是否合法
	paramModelBaseV reflect.Value

	// dest
	scanDest  any
	scanIsPtr bool

	// dest去除ptr的value
	destV              reflect.Value
	destBaseType       reflect.Type
	destBaseTypeIsComp bool
	// scan 为slice时，里面item是否是ptr
	destIsSlice        bool
	destSliceItemIsPtr bool

	log Logger
	err error

	tableSqlType tableSqlType //单表，sql类型crud

	baseTokens []baseToken

	isLgDel bool //是否启用了逻辑删除
	isTen   bool //是否启用了多租户

	// ------------------主键----------------------

	// id = 1
	//主键名-列表,这里考虑到多主键
	primaryKeyNames []string
	//主键值-列表
	primaryKeyValues [][]field.Value

	// id != 1 ,使用场景 更新名字时，检查名字重复，排除自己
	//主键名-列表,这里考虑到多主键-排除
	filterPrimaryKeyNames []string
	//主键值-列表-排除
	filterPrimaryKeyValues [][]any

	// ------------------conf----------------------

	insertType     insert_type.InsertType
	returnType     return_type.ReturnType
	softDeleteType soft_delete.SoftDelType
	skipSoftDelete bool   // 跳过软删除
	tableName      string //当前表名
	checkParam     bool   // 是否检查参数
	showSql        bool   // 是否打印sql
	// ------------------conf-end----------------------

	// ------------------字段名：字段值----------------------

	columns      []string      // 有效字段列表
	columnValues []field.Value // 有效字段值

	modelZeroFieldNames []string // model 零值字段列表
	modelAllFieldNames  []string // model 所有字段列表
	// ------------------字段名：字段值-end----------------------

	//------------------scan----------------------
	//true query,false exec
	sqlIsQuery bool
	sqlType    sql_type.SqlType

	//要执行的sql语句
	query *strings.Builder
	//参数
	args []any

	started bool
}

func (ctx *ormContext) initExtra(extra ...*ExtraContext) {
	var e *ExtraContext
	if len(extra) > 0 && extra[0] != nil {
		e = extra[0]
	} else {
		e = Extra()
	}
	// err 上抛到 ormContext
	if e.GetErr() != nil {
		ctx.err = e.GetErr()
		return
	}
	ctx.extra = e
	ctx.insertType = e.insertType
	ctx.returnType = e.returnType
	ctx.showSql = e.showSql
	ctx.skipSoftDelete = e.skipSoftDelete
	ctx.tableName = e.tableName
}

// 初始化 表名,主键
func (ctx *ormContext) initConf() {
	if ctx.hasErr() {
		return
	}
	v := ctx.destV
	dest := ctx.scanDest
	t := ctx.destBaseType
	ctx.softDeleteType = utils.GetSoftDelType(t)

	if ctx.tableName == "" {
		tableName := ctx.ormConf.tableName(v, dest)
		ctx.tableName = tableName
	}

	primaryKeys := ctx.ormConf.primaryKeys(v, dest)
	ctx.primaryKeyNames = primaryKeys
}

// 获取struct对应的字段名 和 其值，
// slice为全部，一个为非nil字段。
func (ctx *ormContext) initColumnsValue() {
	if ctx.hasErr() {
		return
	}
	cv, err := ctx.ormConf.getStructFV(ctx.destV)
	if err != nil {
		ctx.err = err
		return
	}
	ctx.columns = cv.columns
	ctx.columnValues = cv.columnValues

	ctx.modelZeroFieldNames = cv.modelZeroFieldNames

	ctx.modelAllFieldNames = cv.modelAllFieldNames

	ctx.initColumnsValueSet()
	ctx.initColumnsValueExtra()
	ctx.initColumnsValueSoftDel()
	return
}
func (ctx *ormContext) initColumnsValueSet() {
	if ctx.hasErr() {
		return
	}
	e := ctx.extra
	set := e.set
	if set.hasModel {
		oc := &ormContext{
			ormConf:        ctx.ormConf,
			skipSoftDelete: true,
		}
		oc.initModelDest(set.model) //初始化参数
		oc.initColumnsValue()       //初始化cv

		set.columns = append(set.columns, oc.columns...)
		set.columnValues = append(set.columnValues, oc.columnValues...)
	}

	return
}
func (ctx *ormContext) initColumnsValueExtra() {
	if ctx.hasErr() {
		return
	}
	e := ctx.extra
	if e == nil {
		return
	}
	for i, column := range e.columns {
		cv := e.columnValues[i]
		if cv.Type == field.Null || cv.Type == field.Now {
			ctx.modelZeroFieldNames = append(ctx.modelZeroFieldNames, column)
		}
		find := utils.Find(ctx.columns, column)
		if find == -1 {
			ctx.columns = append(ctx.columns, column)
			ctx.columnValues = append(ctx.columnValues, e.columnValues[i])
		} else {
			ctx.columnValues[i] = e.columnValues[i]
		}
	}
	return
}
func (ctx *ormContext) initColumnsValueSoftDel() {
	if ctx.hasErr() {
		return
	}
	if ctx.skipSoftDelete {
		return
	}

	switch ctx.sqlType {
	case sql_type.Insert:
		value, has := soft_delete.SoftDelTypeNoFVMap[ctx.softDeleteType]
		if has && value.Type != field.None {
			ctx.columns = append(ctx.columns, value.Name)
			ctx.columnValues = append(ctx.columnValues, value.ToValue())
		}
		break
	case sql_type.Delete:
		value, has := soft_delete.SoftDelTypeYesFVMap[ctx.softDeleteType]
		if has && value.Type != field.None {
			ctx.columns = append(ctx.columns, value.Name)
			ctx.columnValues = append(ctx.columnValues, value.ToValue())
		}
		break
	default:
		break
	}
	return
}

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

func (ctx *ormContext) initPrimaryKeyValues(v []any) {
	if ctx.hasErr() {
		return
	}

	idLen := len(v)
	if idLen == 0 {
		ctx.err = errors.New("ByPrimaryKey arg len num 0")
		return
	}
	pkLen := len(ctx.primaryKeyNames)

	idValuess := make([][]field.Value, 0)

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

			idValues := make([]field.Value, 1)
			idValues[0] = field.Value{
				Type:  field.Val,
				Value: value.Interface(),
			}
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

			idValues := make([]field.Value, 0)
			idValues = append(idValues, values...)
			idValuess = append(idValuess, idValues)
		}
	}

	ctx.primaryKeyValues = idValuess
}

func (ctx *ormContext) initSelfPrimaryKeyValues() {
	if ctx.hasErr() {
		return
	}

	keyNum := len(ctx.primaryKeyNames)
	idValues := make([]field.Value, 0)
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

func (ctx *ormContext) hasErr() bool {
	return ctx.err != nil
}
