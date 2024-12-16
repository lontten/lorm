package lorm

import (
	"github.com/lontten/lorm/field"
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"time"
)

// has
//func (ldb coreDb) doHas() (bool, error) {
//if err := ldb.getCtx().err; err != nil {
//	return false, err
//}
//var bb bytes.Buffer
//
//ctx := ldb.getCtx()
//tableName := ctx.tableName
//
//bb.WriteString("SELECT 1 FROM ")
//bb.WriteString(tableName)
//bb.Write(ldb.genWhereSqlByToken())
//bb.WriteString("LIMIT 1")
//
//rows, err := ldb.doQuery(bb.String(), ldb.args...)
//
//if err != nil {
//	return false, err
//}
//defer rows.Close()
//if rows.Next() {
//	return true, nil
//}
//return false, nil
//}

//-------------------------------initContext------------------------

// 根据 byModel 生成的where token
func (d *MysqlDialect) initByPrimaryKey() {
	//ctx := ldb.ctx
	//if err := ctx.err; err != nil {
	//	return
	//}
	//pkNum := len(ctx.primaryKeyValues)
	//ldb.whereTokens = append(ldb.whereTokens, utils.GenwhereTokenOfBatch(ctx.primaryKeyNames, pkNum))
	//
	//for _, value := range ctx.primaryKeyValues {
	//	ldb.args = append(ldb.args, value...)
	//}
}

// 根据 byModel 生成的where token
func (d *MysqlDialect) initByModel(v any) {
	//if err := ldb.ctx.err; err != nil {
	//	return
	//}
	//if v == nil {
	//	ldb.ctx.err = errors.New("model is nil")
	//	return
	//}
	//
	//columns, values, err := getStructCV(v, ldb.ctx.ormConf)
	//if err != nil {
	//	ldb.ctx.err = err
	//	return
	//}
	//ldb.whereTokens = append(ldb.whereTokens, utils.GenwhereToken(columns)...)
	//ldb.args = append(ldb.args, values...)
}

// 根据 byWhere 生成的where token
func (db *lnDB) initByWhere(w *WhereBuilder) {
	//if err := ldb.core.getCtx().err; err != nil {
	//	return
	//}
	//if w == nil {
	//	ldb.core.getCtx().err = errors.New("ByWhere is nil")
	//	return
	//}
	//
	//args := w.args
	//toSql, err := w.toSql(ldb.core.getDialect().parse)
	//if err != nil {
	//	ldb.core.getCtx().err = err
	//	return
	//}
	//ldb.whereTokens = append(ldb.whereTokens, toSql)
	//ldb.args = append(ldb.args, args...)
}

// init 逻辑删除、租户
func (d *MysqlDialect) initExtra() {
	//if err := ldb.ctx.err; err != nil {
	//	return
	//}
	//
	//if ldb.ctx.ormConf.LogicDeleteYesSql != "" {
	//	ldb.whereTokens = append(ldb.whereTokens, ldb.ctx.ormConf.LogicDeleteYesSql)
	//}
	//
	//if ldb.ctx.ormConf.TenantIdFieldName != "" {
	//	ldb.whereTokens = append(ldb.whereTokens, ldb.ctx.ormConf.TenantIdFieldName)
	//	ldb.args = append(ldb.args, ldb.ctx.ormConf.TenantIdValueFun())
	//}
	//
	//var sb strings.QueryBuild
	//sb.WriteString(ldb.extraWhereSql)
	//
	//if len(ldb.orderByTokens) > 0 {
	//	sb.WriteString(" ORDER BY ")
	//	sb.WriteString(strings.Join(ldb.orderByTokens, ","))
	//}
	//if ldb.limit > 0 {
	//	sb.WriteString(" LIMIT ? ")
	//	ldb.args = append(ldb.args, ldb.limit)
	//}
	//if ldb.offset > 0 {
	//	sb.WriteString(" OFFSET ? ")
	//	ldb.args = append(ldb.args, ldb.offset)
	//}
	//ldb.extraWhereSql = sb.String()

}

// -------------------------utils------------------------
// 获取comp 的 cv
// 排除 nil 字段
func getCompCV(v any, c *OrmConf) ([]string, []field.Value, error) {
	value := reflect.ValueOf(v)
	_, value, err := basePtrDeepValue(value)
	if err != nil {
		return nil, nil, err
	}

	return getCompValueCV(value)
}

// 排除 nil 字段
func getCompValueCV(v reflect.Value) ([]string, []field.Value, error) {
	if !isCompType(v.Type()) {
		return nil, nil, errors.New("getvcv not comp")
	}
	err := checkCompFieldVS(v)
	if err != nil {
		return nil, nil, err
	}

	cv, err := getStructCV(v)
	if err != nil {
		return nil, nil, err
	}
	if len(cv.columns) < 1 {
		return nil, nil, errors.New("where model valid field need ")
	}
	return cv.columns, cv.columnValues, nil
}

//------------------------gen-sql---------------------------

// 根据 columnValues 生成的 VALUES sql
// INSERT INTO table_name (列1, 列2,...) VALUES (值1, 值2,....)
func (ctx *ormContext) genInsertValuesSqlBycolumnValues() {
	columns := ctx.columns
	values := ctx.columnValues
	var query = ctx.query

	for i, v := range values {
		if i > 0 {
			query.WriteString(" , ")
		}
		switch v.Type {
		case field.None:
			break
		case field.Null:
			query.WriteString("NULL")
			break
		case field.Now:
			query.WriteString("NOW()")
			break
		case field.UnixSecond:
			query.WriteString(strconv.Itoa(time.Now().Second()))
			break
		case field.UnixMilli:
			query.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))
			break
		case field.UnixNano:
			query.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
			break
		case field.Val:
			query.WriteString(" ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Increment:
			query.WriteString(columns[i] + " + ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Expression:
			query.WriteString(v.Value.(string))
			break
		case field.ID:
			if len(ctx.primaryKeyNames) > 0 {
				ctx.err = errors.New("软删除标记为主键id，需要单主键")
				return
			}
			query.WriteString(ctx.primaryKeyNames[0])
			break
		}
	}
}

// 根据 columnValues 生成的set sql
// SET ...
// column1 = value1, column2 = value2, ...
func (ctx *ormContext) genSetSqlBycolumnValues() {
	columns := ctx.columns
	values := ctx.columnValues
	var query = ctx.query

	for i, v := range values {
		if i > 0 {
			query.WriteString(" , ")
		}
		switch v.Type {
		case field.None:
			break
		case field.Null:
			query.WriteString(columns[i])
			query.WriteString(" = NULL")
			break
		case field.Now:
			query.WriteString(columns[i])
			query.WriteString(" = NOW()")
			break
		case field.UnixSecond:
			query.WriteString(columns[i])
			query.WriteString(" = ")
			query.WriteString(strconv.Itoa(time.Now().Second()))
			break
		case field.UnixMilli:
			query.WriteString(columns[i])
			query.WriteString(" = ")
			query.WriteString(strconv.FormatInt(time.Now().UnixMilli(), 10))
			break
		case field.UnixNano:
			query.WriteString(columns[i])
			query.WriteString(" = ")
			query.WriteString(strconv.FormatInt(time.Now().UnixNano(), 10))
			break
		case field.Val:
			query.WriteString(columns[i])
			query.WriteString(" = ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Increment:
			query.WriteString(columns[i])
			query.WriteString(" = ")
			query.WriteString(columns[i] + " + ? ")
			ctx.args = append(ctx.args, v.Value)
			break
		case field.Expression:
			query.WriteString(columns[i])
			query.WriteString(" = ")
			query.WriteString(v.Value.(string))
			break
		case field.ID:
			if len(ctx.primaryKeyNames) > 0 {
				ctx.err = errors.New("软删除标记为主键id，需要单主键")
				return
			}
			query.WriteString(columns[i])
			query.WriteString(" = ")
			query.WriteString(ctx.primaryKeyNames[0])
			break
		}
	}
}
