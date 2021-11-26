package lorm

import (
	"bytes"
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

type OrmContext struct {
	//主键名-列表
	primaryKeyNames []string

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

//args 为 where 的 字段名列表， 生成where sql
//sql 为 逻辑删除 附加where
//todo 应该改为 统一 where sql 统一生成、  逻辑删除、 多租户
func (ctx *OrmContext) tableWherePrimaryKey2SqlStr(ids []string) string {
	c := ormConfig
	var sb strings.Builder
	for i, where := range ids {
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
	lgSql := strings.ReplaceAll(c.LogicDeleteNoSql, "lg.", "")
	if c.LogicDeleteNoSql != lgSql {
		sb.WriteString(" AND ")
		sb.WriteString(lgSql)
	}

	if c.TenantIdFieldName != "" {
		sb.WriteString(" AND ")
		sb.WriteString(c.TenantIdFieldName)
		sb.WriteString(" = ? ")

		ctx.args = append(ctx.args, c.TenantIdValueFun())
	}
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

func (c *OrmContext) checkValidPrimaryKey(v []interface{}) {
	ids := c.primaryKeyNames
	//主键名列表长度为1，单主键
	singlePk := len(ids) == 1

	value := reflect.ValueOf(v[0])
	is, base := basePtrValue(value)
	if is && value.IsNil() { //数值无效，直接返回false，不再进行合法性检查
		c.err = errors.New("PrimaryKey  : is nil")
		return
	}

	is, base, err := checkDestSingle(base)
	if err != nil {
		c.err = err
		return
	}
	if is {
		if !singlePk {
			c.err = errors.New("PrimaryKey arg is err")
			return
		}
		c.args = append(c.args, v...)
		return
	}

	err = checkStructValidFieldNuller(base)
	if err != nil {
		c.err = err
		return
	}

	for _, e := range v {
		value = reflect.ValueOf(e)
		_, base = basePtrValue(value)
		for _, id := range ids {
			field := base.FieldByName(id)
			c.args = append(c.args, field.Interface())
		}
	}

}

func (ctx *OrmContext) genDelSql() []byte {
	var bb bytes.Buffer
	keys := ctx.primaryKeyNames
	tableName := ctx.tableName

	whereSql := genWhere(keys)

	logicDeleteSetSql := ormConfig.LogicDeleteSetSql
	logicDeleteYesSql := ormConfig.LogicDeleteYesSql
	lgSql := strings.ReplaceAll(logicDeleteSetSql, "lg.", "")
	logicDeleteYesSql = strings.ReplaceAll(logicDeleteYesSql, "lg.", "")
	if logicDeleteSetSql == lgSql {
		bb.WriteString("DELETE FROM ")
		bb.WriteString(tableName)
		bb.WriteString("WHERE ")
		bb.WriteString(string(whereSql))
	} else {
		bb.WriteString("UPDATE ")
		bb.WriteString(tableName)
		bb.WriteString(" SET ")
		bb.WriteString(lgSql)
		bb.WriteString("WHERE ")
		bb.WriteString(string(whereSql))
		bb.WriteString(" and ")
		bb.WriteString(logicDeleteYesSql)
	}
	return bb.Bytes()

}
