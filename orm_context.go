package lorm

import (
	"reflect"
	"strings"
)

type OrmContext struct {
	primaryKeyNames []string

	//当前表名
	tableName string
	//当前struct对象
	dest []interface{}
	//去除 ptr
	destValue reflect.Value
	//去除 ptr slice
	destBaseValue reflect.Value

	columns      []string
	columnValues []interface{}

	query   *strings.Builder
	args    []interface{}
	started bool
	err     error
	log     int
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
func (ctx *OrmContext) tableWhereArgs2SqlStr(args []string, c OrmConf) string {
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
	lgSql := strings.ReplaceAll(c.LogicDeleteNoSql, "lg.", "")
	if c.LogicDeleteNoSql != lgSql {
		sb.WriteString(" AND ")
		sb.WriteString(lgSql)
	}
	return sb.String()
}

//args 为 where 的 字段名列表， 生成where sql
//sql 为 逻辑删除 附加where
//todo 应该改为 统一 where sql 统一生成、  逻辑删除、 多租户
func (ctx *OrmContext) tableWherePrimaryKey2SqlStr(ids []string, c OrmConf) string {

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
func (ctx *OrmContext) tableCreateArgs2SqlStr(args []string) string {
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