package lorm

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type OrmContext struct {
	//主键名-列表
	primaryKeyNames []string
	//主键值-列表
	primaryKeyValues [][]interface{}

	//当前表名
	tableName string

	//当前struct对象
	dest interface{}
	//去除 ptr
	destValue reflect.Value
	//用作 参数合法行校验
	destBaseValue reflect.Value
	isSlice       bool
	// dest 的value 列表，用作参数
	destValueArr []reflect.Value

	//字段列表
	columns []string
	//值列表-多个
	columnValues [][]interface{}

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

//args 为 where 的 字段名列表， 生成where sql
//sql 为 逻辑删除 附加where
//todo 应该改为 统一 where sql 统一生成、  逻辑删除、 多租户
func (ctx *OrmContext) tableWhereArgs2SqlStr(args []string) string {
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

	//lgSql := strings.ReplaceAll(c.LogicDeleteNoSql, "lg.", "")
	//if c.LogicDeleteNoSql != lgSql {
	//	sb.WriteString(" AND ")
	//	sb.WriteString(lgSql)
	//}
	return sb.String()
}

// create 生成
func (ctx *OrmContext) tableCreateArgs2SqlStr() string {
	args := ctx.columns
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

func (ctx *OrmContext) checkValidPrimaryKey(v []interface{}) {
	ids := ctx.primaryKeyNames
	//主键名列表长度为1，单主键
	singlePk := len(ids) == 1

	value := reflect.ValueOf(v[0])
	fmt.Println(value.String())
	fmt.Println(value.Kind())
	is, base, err := basePtrValue(value)
	if err != nil { //数值无效，直接返回false，不再进行合法性检查
		ctx.err = err
		return
	}

	is, base, err = checkDestSingle(base)
	if err != nil {
		ctx.err = err
		return
	}
	if is {
		if !singlePk {
			ctx.err = errors.New("PrimaryKey arg is err")
			return
		}
		ctx.args = append(ctx.args, v...)
		return
	}

	err = checkCompField(base)
	if err != nil {
		ctx.err = err
		return
	}

	for _, e := range v {
		value = reflect.ValueOf(e)
		_, base, err = basePtrValue(value)
		if err != nil {
			ctx.err = err
			return
		}
		for _, id := range ids {
			field := base.FieldByName(id)
			ctx.args = append(ctx.args, field.Interface())
		}
	}

}

//v0.7
//生成select sql
func (ctx *OrmContext) genSelectByPrimaryKey() []byte {
	tableName := ctx.tableName
	columns := ctx.columns
	selSql := ormConfig.genSelectSqlCommon(tableName, columns)
	where := ctx.genWhereByPrimaryKey()
	return append(selSql, where...)
}

//v0.6
//生成del sql
func (ctx *OrmContext) genDelByPrimaryKey() []byte {
	keys := ctx.primaryKeyNames
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName, ctx.destBaseValue)
	return ormConfig.genDelSqlCommon(tableName, keys, hasTen)

}

//v0.6
//生成del sql
func (ctx *OrmContext) genDel(keys []string) []byte {
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName, ctx.destBaseValue)
	return ormConfig.genDelSqlCommon(tableName, keys, hasTen)

}

//v0.6
//生成where sql
func (ctx *OrmContext) genWhereByPrimaryKey() []byte {
	keys := ctx.primaryKeyNames
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName, ctx.destBaseValue)
	return ormConfig.GenWhere(keys, hasTen)
}

//v0.6
//生成where sql
func (ctx *OrmContext) genWhere(keys []string) []byte {
	tableName := ctx.tableName
	//开启多租户，并且该表不跳过
	hasTen := ormConfig.TenantIdFieldName != "" && !ormConfig.TenantIgnoreTableFun(tableName, ctx.destBaseValue)
	return ormConfig.GenWhere(keys, hasTen)
}
