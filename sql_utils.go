package lorm

import "strings"
//select 生成
func selectArgsArr2SqlStr(context OrmContext, args []string) {
	query := context.query
	if context.startd {
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
			context.startd = true
		}
	}
}

//args 为 where 的 字段名列表， 生成where sql
//sql 为 逻辑删除 附加where
//todo 应该改为 统一 where sql 统一生成、  逻辑删除、 多租户
func tableWhereArgs2SqlStr(args []string, sql string) string {
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
	lgSql := strings.ReplaceAll(sql, "lg.", "")
	if sql != lgSql {
		sb.WriteString(" AND ")
		sb.WriteString(lgSql)
	}
	return sb.String()
}

func tableSelectArgs2SqlStr(args []string) string {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	for i, column := range args {
		if i == 0 {
			sb.WriteString(column)
		} else {
			sb.WriteString(" , ")
			sb.WriteString(column)
		}
	}
	return sb.String()
}

// create 生成
func tableCreateArgs2SqlStr(args []string) string {
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
func tableUpdateArgs2SqlStr(args []string) string {
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

