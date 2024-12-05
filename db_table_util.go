package lorm

import (
	"bytes"
	"github.com/lontten/lorm/field"
	"github.com/pkg/errors"
	"reflect"
)

// todo 下面未重构--------------
// update
func (db *lnDB) doUpdate() (int64, error) {
	if err := db.core.getCtx().err; err != nil {
		return 0, err
	}
	var bb bytes.Buffer

	ctx := db.core.getCtx()
	tableName := ctx.tableName
	cs := ctx.columns

	bb.WriteString("UPDATE ")
	bb.WriteString(tableName)
	bb.WriteString(" SET ")
	bb.WriteString(ctx.tableUpdateArgs2SqlStr(cs))
	bb.Write(db.genWhereSqlByToken())

	//return ldb.core.doExec(bb.String(), append(ctx.columnValues, ldb.args...)...)
	return 0, nil
}

// del
//func (ldb coreDb) doDel() (int64, error) {
//if err := ldb.getCtx().err; err != nil {
//	return 0, err
//}
//var bb bytes.Buffer
//tableName := ldb.getCtx().tableName
//w := ldb.genWhereSqlByToken()
//
//if ldb.getCtx().ormConf.LogicDeleteSetSql == "" {
//	bb.WriteString("DELETE FROM ")
//	bb.WriteString(tableName)
//	bb.Write(w)
//} else {
//	bb.WriteString("UPDATE ")
//	bb.WriteString(tableName)
//	bb.WriteString(" SET ")
//	bb.WriteString(ldb.getCtx().ormConf.LogicDeleteSetSql)
//	bb.Write(w)
//}
//return ldb.doExec(bb.String(), ldb.args...)
//return 0, nil
//}

// update
//func (ldb coreDb) doSelect() (int64, error) {
//	var bb bytes.Buffer
//
//	ctx := ldb.getCtx()
//	tableName := ctx.tableName
//	columns := ctx.columns
//
//	bb.WriteString("SELECT ")
//	for i, column := range columns {
//		if i == 0 {
//			bb.WriteString(column)
//		} else {
//			bb.WriteString(" , ")
//			bb.WriteString(column)
//		}
//	}
//	bb.WriteString(" FROM ")
//	bb.WriteString(tableName)
//bb.Write(ldb.genWhereSqlByToken())
//
//return ldb.query(bb.String(), ldb.args...)
//return 0, nil
//}

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
	//var sb strings.Builder
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

// 初始化逻辑删除
func (db *lnDB) initLgDel() {
	//if err := ldb.ctx.err; err != nil {
	//	return
	//}
	//if ldb.ctx.ormConf.LogicDeleteYesSql != "" {
	//	ldb.extraWhereSql = ldb.ctx.ormConf.LogicDeleteYesSql
	//}
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

	return getCompValueCV(value, c)
}

// 排除 nil 字段
func getCompValueCV(v reflect.Value, c *OrmConf) ([]string, []field.Value, error) {
	if !isCompType(v.Type()) {
		return nil, nil, errors.New("getvcv not comp")
	}
	err := checkCompFieldVS(v)
	if err != nil {
		return nil, nil, err
	}

	cv, err := c.getStructCV(v)
	if err != nil {
		return nil, nil, err
	}
	if len(cv.columns) < 1 {
		return nil, nil, errors.New("where model valid field need ")
	}
	return cv.columns, cv.columnValues, nil
}

// ------------------------query--------------------------
func (db *lnDB) query(query string, args ...any) (int64, error) {
	//rows, err := ldb.doQuery(query, args...)
	//if err != nil {
	//	return 0, err
	//}
	//if ldb.ctx.scanIsSlice {
	//	return ldb.ctx.Scan(rows)
	//}
	//return ldb.ctx.ScanLn(rows)
	return 0, nil
}

func (db *lnDB) queryBatch(query string, args [][]any) (int64, error) {
	//query = ldb.dialect.queryBatch(query)
	//
	//stmt, err := ldb.doPrepare(query)
	//if err != nil {
	//	return 0, err
	//}
	//
	//rowss := make([]*sql.Rows, 0)
	//for _, arg := range args {
	//	rows, err := stmt.Query(arg...)
	//	if err != nil {
	//		return 0, err
	//	}
	//	rowss = append(rowss, rows)
	//}
	//return ldb.ctx.ScanBatch(rowss)
	return 0, nil
}

//------------------------gen-sql---------------------------

// 根据whereTokens生成的where sql
func (db *lnDB) genWhereSqlByToken() []byte {
	//if len(ldb.whereTokens) == 0 && ldb.extraWhereSql == "" {
	//	return nil
	//}
	//var buf bytes.Buffer
	//buf.WriteString(" WHERE ")
	//for i, token := range ldb.whereTokens {
	//	if i > 0 {
	//		buf.WriteString(" AND ")
	//	}
	//	buf.WriteString(token)
	//}
	//buf.WriteString(ldb.extraWhereSql)
	//return buf.Bytes()
	return nil
}
